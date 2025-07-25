import { useEffect, useState } from "react";
import { Routes, Route, useLocation } from "react-router-dom";
import OAuthCallback from "./OAuthCallback";
import AskQuestion from "./AskQuestion";
import UploadFile from "./UploadFile";
import { BACKEND_URL } from "./config";



const Home = () => {
  const clientID = "39u7iped9gp9cfnfutjp1ras8b";
  const redirectURL = "https://askmydoc.dev/oauth/callback";
  const cognitoDomain = "https://us-west-1rdclhxshd.auth.us-west-1.amazoncognito.com";
  const [message, setMessage] = useState<string | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [fileUploaded, setFileUploaded] = useState(false);
  const location = useLocation();

  useEffect(() => {
    if (location.state?.message) {
      setMessage(location.state.message);
    }
  }, [location.state]);

  useEffect(() => {
    const storedToken = sessionStorage.getItem("id_token");
    const incomingMessage = location.state?.message;

    if (storedToken) {
      setToken(storedToken);
      // Only show "Welcome back!" if no message came from the redirect
      setMessage(incomingMessage || "Login successful! Welcome back!");
    } else if (incomingMessage) {
      setMessage(incomingMessage);
    }
  }, [location.state]);

  const handleLogin = () => {
    const loginUrl = `${cognitoDomain}/login?client_id=${clientID}&response_type=code&scope=openid+email+profile&redirect_uri=${encodeURIComponent(redirectURL)}`;
    window.location.href = loginUrl;
  };

  const fetchProfile = () => {
    if (!token) return;

    fetch(`${BACKEND_URL}/api/profile`, {
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
      {token && <UploadFile onUploadSuccess={() => setFileUploaded(true)} />}
      {token && fileUploaded && <AskQuestion />}
    </div>
  );
};

function App() {
  return (
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/oauth/callback" element={<OAuthCallback />} />
      </Routes>
  );
}

export default App;
