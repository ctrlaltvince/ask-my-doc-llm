import { useEffect, useState } from "react";

function App() {
  const clientID = "39u7iped9gp9cfnfutjp1ras8b";
  const redirectURL = "http://localhost:5173";
  const [message, setMessage] = useState<string | null>(null);
  const [token, setToken] = useState<string | null>(null);

  const handleLogin = () => {
    console.log("Login button clicked");
    window.location.href = `https://us-west-1rdclhxshd.auth.us-west-1.amazoncognito.com/login?client_id=${clientID}&response_type=code&scope=openid+email+profile&redirect_uri=${encodeURIComponent(redirectURL)}`;
  };

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const code = params.get("code");

    if (code) {
      fetch("http://localhost:8081/callback", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ code }),
      })
        .then((res) => res.json())
        .then((data) => {
          setToken(data.id_token); // backend should send id_token too
          setMessage(`Login successful! Welcome, ${data.email}`);
        })
        .catch((err) => {
          setMessage("Login failed: " + err.message);
        });
    }
  }, []);

  const fetchProfile = () => {
    if (!token) return;

    fetch("http://localhost:8081/profile", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((res) => {
        if (!res.ok) throw new Error("Unauthorized");
        return res.json();
      })
      .then((data) => {
        setMessage(`Profile email: ${data.email}`);
      })
      .catch((err) => {
        setMessage("Failed to fetch profile: " + err.message);
      });
  };

  return (
    <div>
      <h1>Ask My Doc LLM</h1>
      {!token && <button onClick={handleLogin}>Login with Google</button>}
      {token && <button onClick={fetchProfile}>Get Profile</button>}
      {message && <p>{message}</p>}
    </div>
  );
}

export default App;
