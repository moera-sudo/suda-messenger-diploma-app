import { useEffect, useState } from 'react';
import { api } from '../lib/api';

interface Props {
  mediaId: string;
  username: string;
  size?: number;
}

// Avatar resolves a media id to a presigned URL via the media service
// (GET /media/{id}/url) and renders it. Falls back to the username initial
// when there is no avatar or the lookup fails.
export default function Avatar({ mediaId, username, size = 28 }: Props) {
  const [url, setUrl] = useState<string>('');

  useEffect(() => {
    let cancelled = false;
    setUrl('');
    if (!mediaId) return;
    api
      .get<{ url: string }>(`/media/${mediaId}/url`)
      .then(({ data }) => {
        if (!cancelled) setUrl(data.url);
      })
      .catch(() => {
        /* fallback to initial */
      });
    return () => {
      cancelled = true;
    };
  }, [mediaId]);

  const initial = (username || '?').replace(/^@/, '').charAt(0).toUpperCase();

  return (
    <span
      style={{
        width: size,
        height: size,
        borderRadius: '50%',
        overflow: 'hidden',
        display: 'inline-flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'var(--surface-2, #1F232C)',
        color: 'var(--text-accent)',
        fontSize: Math.round(size * 0.45),
        fontWeight: 700,
        flexShrink: 0,
      }}
    >
      {url ? (
        <img
          src={url}
          alt={username}
          width={size}
          height={size}
          style={{ width: size, height: size, objectFit: 'cover' }}
          // Broken/expired presigned URL → fall back to the initial.
          onError={() => setUrl('')}
        />
      ) : (
        initial
      )}
    </span>
  );
}
