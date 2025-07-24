import { useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { BACKEND_URL } from "./config";

const OAuthCallback = () => {
  const navigate = useNavigate();
  const hasRun = useRef(false);  // prevent double invocation

  useEffect(() => {
    if (hasRun.current) return;
    hasRun.current = true;

    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get('code');

    if (code) {
      fetch(`${BACKEND_URL}/callback`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ code }),
      })
        .then((res) => {
          if (!res.ok) throw new Error("OAuth callback failed");
          return res.json();
        })
        .then((data) => {
          sessionStorage.setItem("access_token", data.access_token);
          sessionStorage.setItem("id_token", data.id_token);
          navigate("/", {
            state: {
              message: `Login successful! Welcome, ${data.email}`,
            },
          });
        })
        .catch((err) => {
          navigate("/", {
            state: {
              message: "Login failed: " + err.message,
            },
          });
        });
    } else {
      navigate('/', { state: { message: 'Login was canceled or failed. Please try again.' } });
    }
  }, [navigate]);

  return <div>Logging you in...</div>;
};

export default OAuthCallback;