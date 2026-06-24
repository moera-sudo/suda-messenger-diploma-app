# Suda — Deployment Guide

Operational guide for running the full Suda stack on a Linux server, accessed remotely via ZeroTier VPN.

> All commands are run from the repository root: `/home/<user>/git/projects/suda/backend/` (or wherever you cloned).

---

## 1. Prerequisites

On the **server** (Linux):

| Tool | Why | Install |
|------|-----|---------|
| Docker 24+ | Container runtime | https://docs.docker.com/engine/install/ |
| Docker Compose v2 | Multi-container orchestration (`docker compose ...`) | Bundled with modern Docker installs |
| OpenSSL | Generates random secrets | Pre-installed on most distros |
| Task CLI (optional) | Convenience commands like `task docker-up-prod` | https://taskfile.dev/installation/ |
| ZeroTier One | VPN client (server-side) | https://www.zerotier.com/download/ |
| `git` | Clone the repo | distro-dependent |

On the **client** (Android emulator / phone / desktop):

| Tool | Why |
|------|-----|
| ZeroTier One | VPN client — join the same network as the server |
| Browser / Flutter app | Reach `http://<server-zt-ip>/...` once connected |

---

## 2. First-Time Setup

### 2.1 Clone & enter

```bash
git clone <repo-url>
cd suda/backend
```

### 2.2 Generate secrets

```bash
task secrets-generate
```

This creates random 32-byte hex files in `secrets/`:
- `wallet_encryption_key.txt` — AES-256-GCM key for `tx_wallets.encrypted_private_key`
- `jwt_secret.txt` — HMAC for access / refresh tokens
- `gateway_signature_secret.txt` — HMAC the api-gateway adds to downstream requests

**Manual steps after `task secrets-generate`:**

1. **Treasury ETH private key** — paste your dev-treasury PK (64 hex chars, **no `0x` prefix, no trailing newline**):
   ```bash
   # If you already have it in main .env:
   grep '^ETH_PRIVATE_KEY=' .env | cut -d= -f2- | tr -d '\n' > secrets/eth_private_key.txt
   wc -c secrets/eth_private_key.txt   # must be 64
   ```

2. **Firebase service account** — download `firebase-adminsdk.json` from Firebase Console → Project Settings → Service accounts → Generate new private key, save as:
   ```
   secrets/firebase-adminsdk.json
   ```

### 2.3 .env file

Copy the template and fill in non-secret values (SMTP creds, contract URLs, etc.):
```bash
cp backend/docs/examples/.env.example .env
$EDITOR .env
```

Secrets (JWT, WALLET_ENCRYPTION_KEY, ETH_PRIVATE_KEY, GATEWAY_SIGNATURE_SECRET) live in `secrets/*.txt` and are **not** read from `.env` when running through compose — `readSecret()` in each service prefers `*_FILE` env vars first.

### 2.4 Start the prod stack

```bash
task docker-up-prod
```

This runs `docker compose -f docker-compose.yml up -d --build` (the `-f` skips `docker-compose.override.yml` which is dev-only).

Wait ~60 s for first start: image build, migrations, contracts deploy.

### 2.5 Verify

```bash
task docker-ps
```

Expected state:

| Service | Status |
|---------|--------|
| postgres, redis, minio, besu-node | `Up (healthy)` |
| migrate, contracts-deploy | `Exited (0)` |
| api-gateway, messenger-service, transaction-service, media-service | `Up` |
| nginx, portainer | `Up` |

Quick smoke:
```bash
curl -sI http://localhost/                                # 200, landing
curl -sI http://localhost/swagger/messenger/index.html    # 200
curl -sI http://localhost/swagger/transaction/index.html  # 200
curl -sI http://localhost:9443/                           # 200, portainer
```

---

## 3. Daily Operations

```bash
# Status / logs
task docker-ps
task docker-logs                              # all services, follow
docker compose logs -f messenger-service      # single service

# Restart one service
docker compose restart messenger-service

# Rebuild after Go code change
task docker-rebuild SVC=transaction-service

# Connect to Postgres
docker compose exec postgres psql -U postgres suda

# Run a fresh migration that you just added
docker compose run --rm migrate

# Stop everything (data preserved)
task docker-down

# DESTRUCTIVE — wipes postgres / besu / minio data, forces re-deploy of contracts
task docker-down-volumes
```

---

## 4. Portainer (Web UI for Docker)

### First login

1. Open in your browser: `http://<server-ip>:9443/` (use `localhost` if on the server itself, or ZT-IP from outside)
2. Portainer will prompt to create an **admin user**. Username + password ≥ 12 characters.
3. Choose **Get Started** to manage the local Docker engine (Portainer talks to the socket mounted at `/var/run/docker.sock`).

### What you can do

- **Containers** — view logs, stats, environment, restart, exec into shell
- **Stacks** — see `backend` stack (auto-detected from compose)
- **Volumes** — inspect / browse files inside named volumes
- **Networks** — `backend_app-network` with all service IPs

### Security note

Portainer has root access to Docker via the socket mount. Keep port `9443` **only** behind ZeroTier (see §5.3 Firewall).

---

## 5. ZeroTier Networking

The server has no public IP. All access goes through ZeroTier VPN — your devices connect to the same private network as the server and reach it by a virtual IP (typically `192.168.192.x` or `10.x.x.x`).

### 5.1 ZeroTier Central — Legacy vs New

ZeroTier has two web portals:

- **Legacy Central** — https://my.zerotier.com — original portal, still fully supported, free up to 25 nodes. **Use this one.**
- **New Central / ZeroTier Cloud** — https://www.zerotier.com/cloud/ — newer commercial product with SSO and enterprise features. Not needed for diploma / pet projects.

If your network already exists in Legacy Central — no need to migrate.

### 5.2 Server: install ZeroTier and join the network

**Arch / Manjaro:**
```bash
sudo pacman -S zerotier-one
sudo systemctl enable --now zerotier-one
```

**Debian / Ubuntu:**
```bash
curl -s https://install.zerotier.com | sudo bash
# systemd unit is set up automatically
```

**Fedora / RHEL:**
```bash
curl -s https://install.zerotier.com | sudo bash
```

Join your network (replace `<NETWORK_ID>` with the 16-hex-char ID from my.zerotier.com):
```bash
sudo zerotier-cli join <NETWORK_ID>
sudo zerotier-cli status
sudo zerotier-cli listnetworks
```

### 5.3 Authorize the server in Legacy Central

1. Open https://my.zerotier.com/network/<NETWORK_ID> in browser.
2. Scroll to **Members** section. A new entry shows up with **Address** matching `zerotier-cli status` of the server.
3. Tick the **Auth** checkbox.
4. (Optional) Set a **Name** like `suda-server` to recognize it later.
5. Back on the server, after ~10 s:
   ```bash
   sudo zerotier-cli listnetworks
   ```
   Status should change from `REQUESTING_CONFIGURATION` to `OK` and an IP appears.

Find the assigned IP:
```bash
ip addr show | grep -A1 'zt'
# OR specifically:
sudo zerotier-cli get <NETWORK_ID> ip4
```

Example output:
```
ZTV4_IP_0=192.168.192.5
```

**Remember this IP** — that's how clients reach the server.

### 5.4 Server firewall (if any)

The server already opens ports `:80` and `:9443` on `0.0.0.0` (all interfaces). If a host firewall is active, allow ZeroTier interface:

**ufw (Debian / Ubuntu):**
```bash
sudo ufw allow in on zt+ to any port 80
sudo ufw allow in on zt+ to any port 9443
sudo ufw allow in on zt+ to any port 22
sudo ufw reload
```

(`zt+` is a wildcard — every ZT network creates an interface named like `ztabcd1234`.)

**firewalld (Fedora / RHEL):**
```bash
ZT_IFACE=$(ip -o link show | awk -F': ' '/zt/{print $2; exit}')
sudo firewall-cmd --permanent --zone=trusted --add-interface=$ZT_IFACE
sudo firewall-cmd --reload
```

**No firewall (default Arch):** nothing to do.

### 5.5 Client: Android Emulator on Windows

The Android emulator on Windows runs in its own QEMU VM with a separate network stack. **It does NOT see a ZT client installed on the Windows host.** Install ZT *inside* the emulator like a regular Android app:

1. Launch the emulator (Pixel 6 / Android 14+ recommended; needs Play Services).
2. Open **Google Play** in the emulator.
3. Search for `ZeroTier One` (publisher: ZeroTier, Inc.). Install.
4. Open the app → tap `+` (Add Network) → enter `<NETWORK_ID>`.
5. Android shows a VPN-permission dialog → **OK / Allow**.
6. Authorize the new member in https://my.zerotier.com/network/<NETWORK_ID> (Members → tick Auth, name it `android-emu`).
7. Inside the app you'll see a Managed IP — note it.

Verify from inside the emulator:
- Open Chrome → `http://<server-zt-ip>/swagger/messenger/index.html` — should load Swagger UI.

### 5.6 Client: Real Android phone

Same flow as 5.5:
1. Play Store → `ZeroTier One` → Install.
2. Add Network → `<NETWORK_ID>` → Allow VPN.
3. Authorize in Central, name it (e.g. `pixel-me`).
4. Open browser → `http://<server-zt-ip>/...`.

### 5.7 Client: Windows desktop (optional, for direct dev)

If you want to hit the server from your Windows machine directly (without going through the emulator):

1. Download ZeroTier One for Windows: https://www.zerotier.com/download/
2. Install. Tray icon appears.
3. Tray icon → Join Network → enter `<NETWORK_ID>`.
4. Authorize in Central.
5. In PowerShell:
   ```powershell
   ipconfig | Select-String "ZeroTier"
   curl http://<server-zt-ip>/swagger/messenger/index.html
   ```

### 5.8 Client: Linux / macOS desktop (optional)

Same package as the server: `pacman -S zerotier-one` / `brew install --cask zerotier-one`. Then `sudo zerotier-cli join <NETWORK_ID>` + authorize in Central.

---

## 6. SSH

The server's `sshd` runs on the host OS, not in Docker. Once both sides are in the same ZT network, SSH works over the ZT-IP:

```bash
ssh user@<server-zt-ip>
```

If sshd is bound only to `127.0.0.1`, edit `/etc/ssh/sshd_config`:
```
ListenAddress 0.0.0.0
```
Then `sudo systemctl restart sshd`.

Public-key auth is recommended (`ssh-copy-id user@<server-zt-ip>`).

---

## 7. Troubleshooting

### `task docker-up-prod` fails on `contracts-deploy` exit 1
- Check `docker compose logs contracts-deploy`.
- Most common: `secrets/eth_private_key.txt` has wrong format. Must be 64 hex chars, no `0x`, no trailing newline. Verify: `wc -c secrets/eth_private_key.txt` → exactly `64`.
- Besu not minting blocks: `docker compose logs besu-node` should show `Created block #N`. If it stays at block 0 — validator-key mount is broken (see compose `--node-private-key-file`).

### transaction-service in `Restarting` loop
- `docker compose logs transaction-service`.
- `SudaToken at 0x... has no code` → contracts not deployed yet. Wait for `contracts-deploy` to finish, then `docker compose up -d transaction-service`.
- `failed to connect to 127.0.0.1:5432` → env var `DATABASE_URL` is missing (only happens after env var rename — re-check compose).

### `messenger-service` errors on Firebase init
- `secrets/firebase-adminsdk.json` exists? Valid JSON? `python -m json.tool secrets/firebase-adminsdk.json`.
- Service reads `FIREBASE_CREDENTIALS=/run/secrets/firebase-adminsdk.json` (mounted via bind).

### nginx `502 Bad Gateway` on `/api/...`
- `docker compose ps api-gateway` — is it Up?
- If not — logs of api-gateway. Otherwise nginx upstream IP changed (`docker compose restart nginx`).

### Migration failed mid-way ("Dirty database version X")
- Check `docker compose logs migrate`.
- Fix the offending SQL in `migrations/` if needed.
- Force version back: `task docker-migrate-force version=<N-1>` (uses host `migrate` CLI against the docker postgres).
- Re-run: `docker compose run --rm migrate`.

### Portainer not reachable on `:9443`
- It's prod-only; on `task docker-up` (dev mode) it's disabled via profile. Use `task docker-up-prod`.
- Another app using `:9443`? `ss -ltn | grep 9443`. Change port in compose if needed.

### "Volume in use" when trying `docker volume rm`
- Stop AND remove the container that holds it: `docker compose rm -fsv <service>`.

---

## 8. Backups (manual, optional)

Not automated in current setup. Periodic dump:

```bash
mkdir -p backups
docker compose exec -T postgres pg_dump -U postgres suda > backups/suda-$(date +%F).sql
```

Restore into a freshly recreated stack:
```bash
docker compose down -v
docker compose -f docker-compose.yml up -d postgres
sleep 5
docker compose exec -T postgres psql -U postgres suda < backups/suda-2026-05-27.sql
docker compose -f docker-compose.yml up -d
```

---

## 9. Updating the stack

```bash
git pull
task docker-up-prod    # rebuilds any changed images, restarts only updated containers
```

Volumes / secrets / `secrets/contracts.env` are preserved. If migrations were added — they apply automatically on next start (the `migrate` init container is idempotent).
