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
    const updRes = await fetch(`/api/session/${sid}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ status: "pending_wallet" }),
    });
    if (!updRes.ok) throw new Error("Failed to update session");

    try {
      [address] = await ethereum.request({ method: "eth_requestAccounts" });
    } catch (err) {
      deleteSession();
      return;
    }
  }

  console.log("Connected wallet address:", address);
  const updRes = await fetch(`/api/session/${sid}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ wallet: address }),
  });
  if (!updRes.ok) throw new Error("Failed to update session");

  // const updRes = await fetch(`http://localhost:8080/api/session/${sid}`, {
  //   method: "PATCH",
  //   headers: { "Content-Type": "application/json" },
  //   body: JSON.stringify({ status: "pending_wallet" }),
  // });
  // if (!updRes.ok) throw new Error("Failed to update session");

  // await sleep(12000);

  // try {
  //   const [address] = await ethereum.request({ method: "eth_requestAccounts" });

  //   const authRes = await fetch("http://localhost:8080/api/auth", {
  //     method: "POST",
  //     headers: { "Content-Type": "application/json" },
  //     body: JSON.stringify({ session_id, address }),
  //   });
  //   if (!authRes.ok) throw new Error("Auth failed");
  //   const { message } = await authRes.json();

  //   window.onbeforeunload = () => "Please wait—verifying your signature...";
  //   const signature = await ethereum.request({
  //     method: "personal_sign",
  //     params: [message, address],
  //   });

  //   const verifyRes = await fetch("http://localhost:8080/api/verify", {
  //     method: "POST",
  //     headers: { "Content-Type": "application/json" },
  //     body: JSON.stringify({ session_id, message, signature }),
  //   });
  //   if (!verifyRes.ok) throw new Error("Verification failed");

  //   window.onbeforeunload = null;
  // } catch (err) {
  //   window.onbeforeunload = null;
  //   console.error(err);
  // }
})();

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
