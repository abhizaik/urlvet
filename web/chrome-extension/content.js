(async () => {
  const url = window.location.href;

  chrome.runtime.sendMessage({ type: "analyze_url", url }, (response) => {
    if (!response?.success) {
      console.error("url.vet API error:", response?.error);
      return;
    }

    const data = response.data;
    const finalScore = data.result?.final_score ?? 100;

    if (finalScore < 80) {
      showWarningPopup(data);
    }
  });
})();

function showWarningPopup(data) {
  if (document.getElementById("urlvet-warning")) return; // prevent duplicates

  const container = document.createElement("div");
  container.id = "urlvet-warning";
  container.style.position = "fixed";
  container.style.bottom = "20px";
  container.style.right = "20px";
  container.style.zIndex = "99999";
  container.style.padding = "15px";
  container.style.backgroundColor = "#ff4d4f";
  container.style.color = "#fff";
  container.style.fontFamily = "Arial, sans-serif";
  container.style.fontSize = "14px";
  container.style.borderRadius = "8px";
  container.style.boxShadow = "0 4px 10px rgba(0,0,0,0.3)";
  container.style.maxWidth = "300px";

  container.innerHTML = `
    <strong>⚠️ Risky Website Detected!</strong><br>
    URL: ${data.url}<br>
    Verdict: ${data.result.verdict}<br>
    Trust Score: ${data.result.final_score}
    <button id="urlvet-close" style="
      margin-top: 8px;
      padding: 3px 6px;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      background: #fff;
      color: #ff4d4f;
      font-weight: bold;
    ">Close</button>
  `;

  document.body.appendChild(container);

  document.getElementById("urlvet-close").onclick = () => {
    container.remove();
  };
}
