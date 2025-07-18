import React, { useState } from 'react';

export default function UploadFile({ onUploadSuccess }: { onUploadSuccess: () => void }) {
  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleUpload() {
    if (!file) return;
    setLoading(true);
    setStatus('');
    try {
      const token = sessionStorage.getItem("id_token");
      const formData = new FormData();
      formData.append('file', file);

      const res = await fetch('http://localhost:8081/upload', {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      });

      if (!res.ok) throw new Error('Upload failed');

      setStatus('File uploaded successfully!');
      onUploadSuccess(); // ✅ Notify parent
    } catch (e) {
      setStatus('Error uploading file');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      <h2>Upload a file</h2>
      <input type="file" onChange={e => setFile(e.target.files?.[0] || null)} />
      <button onClick={handleUpload} disabled={!file || loading}>
        {loading ? 'Uploading...' : 'Upload'}
      </button>
      {status && <p>{status}</p>}
    </div>
  );
}
