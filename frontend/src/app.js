const API_BASE = "http://localhost:8082";

const formEl = document.getElementById("product-form");
const idInput = document.getElementById("product-id");
const resultEl = document.getElementById("result");

formEl.addEventListener("submit", async (e) => {
  e.preventDefault();

  const res = await fetch(`${API_BASE}/product?id=${idInput.value}`);

  if (!res.ok) {
    resultEl.textContent = "Producto no encontrado";
    return;
  }

  const product = await res.json();

  resultEl.innerHTML = `
    <h3>${product.name}</h3>
    <p><strong>SKU:</strong> ${product.sku}</p>
    <p><strong>Marca:</strong> ${product.brand}</p>
    <p><strong>Categoría:</strong> ${product.category}</p>
    <p><strong>ID:</strong> ${product.id}</p>
  `;
});
