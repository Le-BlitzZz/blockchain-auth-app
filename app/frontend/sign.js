(async function () {
  const sid = window.location.pathname.replace(/^\/+/, "");

  const deleteSession = async () => {
    const delRes = await fetch(`/api/session/${sid}`, {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
    });
    if (!delRes.ok) {
      throw new Error("Failed to delete session");
    }
  };

  const es = new EventSource(`/api/session/${sid}/stream`);
  es.addEventListener("status", async (e) => {
    console.log("Received status:", e.data);
    if (e.data === "gone") {
      es.close();
      window.close();
    }
  });

  window.addEventListener("beforeunload", () => {
    es.close();
  });

  await sleep(1000);

  if (!window.ethereum) {
    deleteSession();
  }

  let address;
  const accounts = await ethereum.request({ method: "eth_accounts" });

  if (accounts.length) {
    address = accounts[0];
  } else {
    const updPendingWalletRes = await fetch(`/api/session/${sid}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ status: "pending_wallet" }),
    });
    if (!updPendingWalletRes.ok) throw new Error("Failed to update session");

    try {
      [address] = await ethereum.request({ method: "eth_requestAccounts" });
    } catch (err) {
      deleteSession();
      return;
    }
  }

  const updWalletRes = await fetch(`/api/session/${sid}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ wallet: address }),
  });
  if (!updWalletRes.ok) throw new Error("Failed to update session");
  const { message } = await updWalletRes.json();

  const signature = await ethereum.request({
    method: "personal_sign",
    params: [message, address],
  });

  const updSignatureRes = await fetch(`http://localhost:8080/api/session/${sid}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ signature }),
  });
  if (!updSignatureRes.ok) throw new Error("Verification failed");

  es.close();
  window.close();
})();

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
