# Secrets

Файлы здесь монтируются в контейнеры как docker secrets. **Никогда** не коммитятся (см. `.gitignore`).

## Что должно быть перед первым `docker compose up`

| Файл | Как получить |
|------|--------------|
| `postgres_password.txt` | `task secrets-generate` |
| `wallet_encryption_key.txt` | `task secrets-generate` (или `openssl rand -hex 32 > ...`) |
| `message_encryption_key.txt` | `openssl rand -hex 32 > secrets/message_encryption_key.txt` (AES-256, шифрование контента сообщений DIRECT/GROUP at-rest; пусто/нет файла → шифрование выключено) |
| `jwt_secret.txt` | `task secrets-generate` |
| `gateway_signature_secret.txt` | `task secrets-generate` |
| `eth_private_key.txt` | вручную — твой treasury PK |
| `firebase-adminsdk.json` | скачать из Firebase Console → Project Settings → Service accounts |

После создания файлов: `docker compose up -d`.