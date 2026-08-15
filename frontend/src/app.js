const API_BASE = "http://localhost:8080";

const formEl = document.getElementById("user-form");
const idInput = document.getElementById("user-id");
const resultEl = document.getElementById("result");

formEl.addEventListener("submit", async (e) => {
  e.preventDefault();
  const res = await fetch(`${API_BASE}/user?id=${idInput.value}`);
  const text = await res.text();
  resultEl.textContent = text;
});
