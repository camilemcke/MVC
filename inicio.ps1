Write-Host "BEAUTY SOURCE"
Write-Host ""

Write-Host "1. Iniciando contenedores"
docker compose up --build -d

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: No se pudieron iniciar los contenedores."
    exit 1
}

Write-Host ""
Write-Host "2. Esperando a que los servicios inicien"
Start-Sleep -Seconds 8

Write-Host ""
Write-Host "3. Probando servicios"

$tests = @(
    @{
        Name = "Middleware"
        URL = "http://localhost:8082/heartbeat"
    },
    @{
        Name = "Products"
        URL = "http://localhost:8082/product?id=1"
    },
    @{
        Name = "Inventory"
        URL = "http://localhost:8082/inventory?id=1"
    },
    @{
        Name = "Analytics"
        URL = "http://localhost:8082/analytics"
    }
)

$allOk = $true

foreach ($test in $tests) {
    try {
        $response = Invoke-WebRequest -Uri $test.URL -UseBasicParsing

        if ($response.StatusCode -eq 200) {
            Write-Host "[OK] $($test.Name)"
        }
    }
    catch {
        Write-Host "[ERROR] $($test.Name)"
        $allOk = $false
    }
}

Write-Host ""

if ($allOk) {
    Write-Host "Todos los servicios estan funcionando"
    Write-Host "Frontend: http://localhost:8081"
}
else {
    Write-Host "Uno o mas servicios no respondieron correctamente"
    exit 1
}