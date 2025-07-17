import { useState, useEffect } from 'react'
import './App.css'

function App() {
  const [token, setToken] = useState<string | null>(null)
  const [file, setFile] = useState<File | null>(null)
  const [question, setQuestion] = useState("")
  const [answer, setAnswer] = useState("")

  useEffect(() => {
    const hash = window.location.hash
    if (hash.includes("id_token") || hash.includes("access_token")) {
      const params = new URLSearchParams(hash.slice(1))
      const accessToken = params.get("access_token")
      if (accessToken) {
        setToken(accessToken)
      }
    }
  }, [])  // empty dependency array means this runs once on mount

  const login = () => {
    // Replace with your Cognito Hosted UI domain
    window.location.href = 'https://your-cognito-domain.auth.us-east-1.amazoncognito.com/login?client_id=XXX&response_type=token&redirect_uri=http://localhost:5173'
  }

  const upload = async () => {
    if (!file || !token) return
    const formData = new FormData()
    formData.append('file', file)

    const res = await fetch('http://localhost:8081/upload', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: formData,
    })
    const json = await res.json()
    console.log(json)
  }

  const ask = async () => {
    const res = await fetch('http://localhost:8081/ask', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ question }),
    })
    const json = await res.json()
    setAnswer(json.answer)
  }

  return (
    <div style={{ maxWidth: 600, margin: "auto", padding: "1rem" }}>
      <h1>🧠 LLM Doc QA</h1>

      {!token && (
        <button onClick={login}>Login with Google</button>
      )}

      {token && (
        <>
          <div>
            <input type="file" onChange={e => setFile(e.target.files?.[0] || null)} />
            <button onClick={upload}>Upload</button>
          </div>

          <div style={{ marginTop: 20 }}>
            <input
              type="text"
              value={question}
              onChange={e => setQuestion(e.target.value)}
              placeholder="Ask a question..."
              style={{ width: "100%" }}
            />
            <button onClick={ask}>Ask</button>
          </div>

          <div style={{ marginTop: 20 }}>
            <strong>Answer:</strong>
            <div>{answer}</div>
          </div>
        </>
      )}
    </div>
  )
}

export default App

