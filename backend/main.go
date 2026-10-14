package main


import (
    "bytes"
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
	"context"
	"os/signal"
	"syscall"

    _ "github.com/lib/pq"

    "backend/controllers"
    "backend/models"
)



func registerWithConsul() {
    instanceID := os.Getenv("INSTANCE_NAME")

    if instanceID == "" {
        log.Println("INSTANCE_NAME is empty; skipping Consul registration")
        return
    }

    service := map[string]interface{}{
        "ID":      instanceID,
        "Name":    "products",
        "Address": instanceID,
        "Port":    8080,
        "Check": map[string]interface{}{
            "HTTP":     fmt.Sprintf("http://%s:8080/health", instanceID),
            "Interval": "10s",
            "Timeout":  "3s",
        },
    }

    data, err := json.Marshal(service)
    if err != nil {
        log.Println("Error creating Consul registration:", err)
        return
    }

    client := &http.Client{Timeout: 3 * time.Second}

    for attempt := 1; attempt <= 10; attempt++ {
        request, err := http.NewRequest(
            http.MethodPut,
            "http://consul:8500/v1/agent/service/register",
            bytes.NewReader(data),
        )
        if err != nil {
            log.Println("Error creating request:", err)
            return
        }

        request.Header.Set("Content-Type", "application/json")

        response, err := client.Do(request)

        if err == nil {
            response.Body.Close()

            if response.StatusCode == http.StatusOK {
                log.Println("Registered in Consul:", instanceID)
                return
            }

            log.Printf("Consul returned status %d", response.StatusCode)
        } else {
            log.Println("Consul registration error:", err)
        }

        log.Println("Retrying Consul registration...")
        time.Sleep(2 * time.Second)
    }

    log.Println("Could not register in Consul:", instanceID)
}


func deregisterFromConsul() {
    instanceID := os.Getenv("INSTANCE_NAME")

    if instanceID == "" {
        log.Println("INSTANCE_NAME is empty; cannot deregister")
        return
    }

    consulURL := fmt.Sprintf(
        "http://consul:8500/v1/agent/service/deregister/%s",
        instanceID,
    )

    request, err := http.NewRequest(
        http.MethodPut,
        consulURL,
        nil,
    )

    if err != nil {
        log.Println("Error creating deregistration request:", err)
        return
    }

    client := &http.Client{Timeout: 3 * time.Second}

    response, err := client.Do(request)
    if err != nil {
        log.Println("Error deregistering from Consul:", err)
        return
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        log.Printf(
            "Consul deregistration returned status %d",
            response.StatusCode,
        )
        return
    }

    log.Println("Deregistered from Consul:", instanceID)
}





func main() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to configure database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	registerWithConsul()

	userModel := &models.UserModel{DB: db}
	userController := &controllers.UserController{UserModel: userModel}
	productModel := &models.ProductModel{DB: db}
	productController := &controllers.ProductController{ProductModel: productModel}

	
	http.HandleFunc("/user", userController.GetUserHandler)
	http.HandleFunc("/product", productController.GetProductHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	
	server := &http.Server{
		Addr: ":8080",
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChannel

		log.Println("Shutdown signal received")

		deregisterFromConsul()

		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Println("Error shutting down server:", err)
		}
	}()

	fmt.Println("Server running on port 8080...")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
