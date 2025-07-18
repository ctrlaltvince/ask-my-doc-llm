import { useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';

const OAuthCallback = () => {
  const navigate = useNavigate();
  const hasRun = useRef(false);  // prevent double invocation

  useEffect(() => {
    if (hasRun.current) return;
    hasRun.current = true;

    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get('code');

    if (code) {
      if (code) {
        fetch("http://localhost:8081/callback", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ code }),
        })
            .then((res) => res.json())
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
        }
    } else {
      navigate('/', { state: { message: 'Login was canceled or failed. Please try again.' } });
    }
  }, [navigate]);

  return <div>Logging you in...</div>;
};

export default OAuthCallback;
