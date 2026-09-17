const API_BASE = "http://localhost:8082";

const formEl = document.getElementById("product-form");
const idInput = document.getElementById("product-id");
const resultEl = document.getElementById("result");

formEl.addEventListener("submit", async (e) => {
  e.preventDefault();

  const productId = idInput.value;

  const productRes = await fetch(`${API_BASE}/product?id=${productId}`);

  if (!productRes.ok) {
    resultEl.textContent = "Producto no encontrado";
    return;
  }

  const product = await productRes.json();

  const inventoryRes = await fetch(`${API_BASE}/inventory?id=${productId}`);

  if (!inventoryRes.ok) {
    resultEl.textContent = "No se pudo consultar el inventario";
    return;
  }

  const inventory = await inventoryRes.json();

  resultEl.innerHTML = `
    <h3>${product.name}</h3>
    <p><strong>SKU:</strong> ${product.sku}</p>
    <p><strong>Marca:</strong> ${product.brand}</p>
    <p><strong>Categoría:</strong> ${product.category}</p>
    <p><strong>Stock:</strong> ${inventory.stock} unidades</p>
    <p><strong>ID:</strong> ${product.id}</p>
  `;
});
