(async function () {
  if (!window.ethereum) {
    document.body.textContent = "MetaMask is required to continue.";
    return;
  }

  const session_id = window.location.pathname.replace(/^\/+/, "");

  try {
    const [address] = await ethereum.request({ method: "eth_requestAccounts" });

    const authRes = await fetch("http://localhost:8080/api/auth", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ session_id, address }),
    });
    if (!authRes.ok) throw new Error("Auth failed");
    const { message } = await authRes.json();

    window.onbeforeunload = () => "Please wait—verifying your signature...";
    const signature = await ethereum.request({
      method: "personal_sign",
      params: [message, address],
    });

    const verifyRes = await fetch("http://localhost:8080/api/verify", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ session_id, message, signature }),
    });
    if (!verifyRes.ok) throw new Error("Verification failed");

    window.onbeforeunload = null;
  } catch (err) {
    window.onbeforeunload = null;
    console.error(err);
  }
})();

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
