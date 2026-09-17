const API_BASE = "http://localhost:8082";

const formEl = document.getElementById("product-form");
const idInput = document.getElementById("product-id");
const resultEl = document.getElementById("result");

async function loadAnalytics() {
  const res = await fetch(`${API_BASE}/analytics`);

  if (!res.ok) {
    return;
  }

  const analytics = await res.json();

  document.getElementById("total-products").textContent =
    analytics.total_products;

  document.getElementById("total-stock").textContent =
    analytics.total_stock;

  document.getElementById("low-stock").textContent =
    analytics.low_stock_products;

  document.getElementById("out-of-stock").textContent =
  analytics.out_of_stock_products;

  document.getElementById("total-entries").textContent =
    analytics.total_entries;

  document.getElementById("total-exits").textContent =
    analytics.total_exits;
}

loadAnalytics();

async function loadBrandChart() {
  const res = await fetch(`${API_BASE}/analytics/brands`);

  if (!res.ok) {
    return;
  }

  const brands = await res.json();
  const chart = document.getElementById("brand-chart");

  chart.innerHTML = "";

  const maxStock = Math.max(...brands.map((brand) => brand.stock));

  brands.forEach((brand) => {
    const percentage = maxStock > 0
      ? (brand.stock / maxStock) * 100
      : 0;

    const row = document.createElement("div");
    row.className = "chart-row";

    row.innerHTML = `
      <div class="chart-label">${brand.brand}</div>

      <div class="chart-bar-container">
        <div
          class="chart-bar"
          style="width: ${percentage}%"
        ></div>
      </div>

      <div class="chart-value">${brand.stock}</div>
    `;

    chart.appendChild(row);
  });
}

loadBrandChart();

async function loadInventory() {
  const res = await fetch(`${API_BASE}/inventory/all`);

  if (!res.ok) {
    return;
  }

  const inventory = await res.json();
  const tableBody = document.getElementById("inventory-table");

  tableBody.innerHTML = "";

  inventory.forEach((item) => {
    let status = "Disponible";

    if (item.stock === 0) {
      status = "Sin stock";
    } else if (item.stock <= 5) {
      status = "Stock bajo";
    }

    const row = document.createElement("tr");

    row.innerHTML = `
      <td>${item.sku}</td>
      <td>${item.brand}</td>
      <td>${item.name}</td>
      <td>${item.entries}</td>
      <td>${item.exits}</td>
      <td>${item.stock}</td>
      <td>${status}</td>
    `;

    tableBody.appendChild(row);
  });
}

loadInventory();
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

  let stockStatus = "Stock disponible";

if (inventory.stock === 0) {
  stockStatus = "Sin stock";
} else if (inventory.stock <= 5) {
  stockStatus = "Stock bajo";
}

resultEl.innerHTML = `
  <h3>${product.name}</h3>
  <p><strong>SKU:</strong> ${product.sku}</p>
  <p><strong>Marca:</strong> ${product.brand}</p>
  <p><strong>Entradas:</strong> ${inventory.entries}</p>
  <p><strong>Salidas:</strong> ${inventory.exits}</p>
  <p><strong>Stock actual:</strong> ${inventory.stock} unidades</p>
  <p><strong>Estado:</strong> ${stockStatus}</p>
  <p><strong>ID:</strong> ${product.id}</p>
`;
});
