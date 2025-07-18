import { useEffect, useState } from "react";

function App() {
  const [message, setMessage] = useState<string | null>(null);

  const handleLogin = () => {
    const domain = "https://us-west-1rdclhxshd.auth.us-west-1.amazoncognito.com";
    const clientId = "39u7iped9gp9cfnfutjp1ras8b";
    const redirectUri = "http://localhost:5173";
    const responseType = "code";
    const scopes = "openid email profile";

    const loginUrl = `${domain}/oauth2/authorize?response_type=code&client_id=${clientId}&redirect_uri=${encodeURIComponent(
      redirectUri
    )}&scope=${encodeURIComponent(scopes)}`;

    window.location.href = loginUrl;
  };

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const code = params.get("code");

    if (code) {
      // Send code to backend at port 8081
      fetch("http://localhost:8081/callback", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ code }),
      })
        .then((res) => res.json())
        .then((data) => {
          setMessage(`Login successful! Welcome, ${data.email}`);
        })
        .catch((err) => {
          setMessage("Login failed: " + err.message);
        });
    }
  }, []);

  return (
    <div>
      <h1>Ask My Doc LLM</h1>
      <button onClick={handleLogin}>Login with Google</button>
      {message && <p>{message}</p>}
    </div>
  );
}

export default App;
