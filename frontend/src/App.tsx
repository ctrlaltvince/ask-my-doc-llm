import { useEffect, useState } from "react";
import { Routes, Route, useLocation } from "react-router-dom";
import OAuthCallback from "./OAuthCallback";
import AskQuestion from "./AskQuestion";
import UploadFile from "./UploadFile";
import { BACKEND_URL } from "./config";

const Home = () => {
  const clientID = import.meta.env.VITE_CLIENT_ID;
  const redirectURL = import.meta.env.VITE_REDIRECT_URL;
  const cognitoDomain = import.meta.env.VITE_COGNITO_DOMAIN;
  const [message, setMessage] = useState<string | null>(null); // Login message
  const [profileInfo, setProfileInfo] = useState<string | null>(null); // Profile display
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

    fetch(`${BACKEND_URL}/profile`, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
    })
      .then((res) => {
        if (!res.ok) throw new Error("Unauthorized");
        return res.json();
      })
      .then((data) => {
        setProfileInfo(`Profile email: ${data.email}`);
      })
      .catch((err) => {
        setProfileInfo("Failed to fetch profile: " + err.message);
      });
  };

  return (
    <div
      style={{
        backgroundImage: 'url("/background.png")',
        backgroundSize: "cover",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
        minHeight: "100vh",
        width: "100vw",
      }}
    >
      <div
        style={{
          backgroundColor: "rgba(0, 0, 0, 0.6)",
          minHeight: "100vh",
          width: "100vw",
          position: "relative",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          padding: "2rem",
          color: "white",
        }}
      >
        {/* Top-left: Profile button + info */}
        {token && (
          <div
            style={{
              position: "absolute",
              top: "1rem",
              left: "1rem",
              display: "flex",
              flexDirection: "column",
              alignItems: "flex-start",
              gap: "0.5rem",
            }}
          >
            <button onClick={fetchProfile}>Get Profile</button>
            {profileInfo && (
              <span style={{ fontSize: "0.9rem" }}>{profileInfo}</span>
            )}
          </div>
        )}

        <h1>Upload & Ask AI | Ask My Doc LLM</h1>
        {!token && <button onClick={handleLogin}>Login with Google</button>}
        {message && <p>{message}</p>}
        {token && <UploadFile onUploadSuccess={() => setFileUploaded(true)} />}
        {token && fileUploaded && <AskQuestion />}
      </div>
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
