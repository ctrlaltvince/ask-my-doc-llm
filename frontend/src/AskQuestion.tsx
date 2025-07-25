import { useState } from 'react';
import { BACKEND_URL } from "./config";

export default function UploadFile({ onUploadSuccess }: { onUploadSuccess: () => void }) {
  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState('');
  const [loading, setLoading] = useState(false);
  const [s3Key, setS3Key] = useState<string | null>(null);

  async function handleUpload() {
    if (!file) return;
    setLoading(true);
    setStatus('');
    setS3Key(null);

    try {
      const token = sessionStorage.getItem("id_token");
      const formData = new FormData();
      formData.append('file', file);

      const res = await fetch(`${BACKEND_URL}/api/upload`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      });

      if (!res.ok) throw new Error('Upload failed');

      const data = await res.json();
      setStatus('File uploaded successfully!');
      setS3Key(data.key);
      onUploadSuccess();
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
      {s3Key && (
        <p>
          <strong>S3 Key:</strong> <code>{s3Key}</code>
        </p>
      )}
    </div>
  );
}
