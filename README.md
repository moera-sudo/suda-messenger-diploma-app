# Suda Messenger — Web3-Integrated Chat System

A cross-platform microservice messaging system featuring real-time encrypted chats, media storage, push notifications, and an embedded Web3 custodial wallet running on a private Hyperledger Besu EVM blockchain. 

---

## 1. Monorepo Directory Structure

```text
├── backend/
│   ├── api-gateway/          # Echo gateway, CORS, JWT verification, HMAC downstream signing
│   ├── messenger-service/    # Main business domain (auth, chats, melody WS hub, encryption)
│   ├── media-service/        # S3-compatible metadata and ACL pipeline
│   ├── transaction-service/  # Custodial Web3 operations, Besu JSON-RPC client, and background Observer
│   └── proto/                # Shared gRPC/protobuf definitions
├── contracts/                # Solidity source code & Hardhat deployment artifacts
├── migrations/               # PostgreSQL schema migrations (migrate/migrate container)
├── web/wallet/               # WebView React wallet SPA (Vite + TS + Zustand)
├── besu/network/             # Hyperledger Besu genesis block & QBFT validator config
├── nginx/                    # Production reverse proxy
└── client/                   # Flutter (Material 3) mobile client application (Android-first)
```

---

## 2. High-Level Architecture

The project employs a secure, low-latency microservice architecture communicating via internal REST and gRPC endpoints.

```text
                        Client (Mobile / WebView)
                                    │
                                    ▼
                                Nginx (:80)
                                    │
                                    ▼
                           API Gateway (:8080)  ─── [JWT & HMAC Signature]
                                    │
          ┌─────────────────────────┼─────────────────────────┐
          ▼                         ▼                         ▼
   Messenger (:8081)       Transactions (:8082)         Media (:8084)
   (Melody WebSockets)     (Besu Node Integration)      (MinIO S3 Store)
          │                         │                         │
          └─────────────────────────┼─────────────────────────┘
                                    ▼
                         PostgreSQL / Redis Cache
```

### gRPC Dependencies
*   **`Messenger ──( gRPC )──> Transactions`**: Checks gating rules (`CheckTokenGating`), and provisions channel treasury wallets (`CreateWalletForChannel`).
*   **`Transactions ──( gRPC )──> Messenger`**: Resolves usernames (`ResolveUsername`), authorizes administrators (`CheckChannelPermission`), checks chat membership (`CheckChatMembership`), and triggers WS events (`NotifyUserEvent`).
*   **`Media ──( gRPC )──> Messenger`**: Authenticates private file access via chat membership rules (`CheckEntityAccess`).

---

## 3. Tech Stack

*   **Backend:** Go (Echo, pgx, melody WS, go-ethereum), gRPC (Protobuf)
*   **Frontend:** Flutter (Dart, Bloc/Cubit) & React (Vite, TS, Zustand)
*   **Database & Cache:** PostgreSQL 16 (raw SQL, pgx pools), Redis 7 (status tracking, pub/sub, locks)
*   **Object Storage:** MinIO / AWS S3 Compatible Storage
*   **Blockchain Network:** Hyperledger Besu (EVM private chain, QBFT consensus, 2-second block period, zero-gas fees)

---

## 4. Key Technical Mechanisms

### At-Rest Encryption
Messages in `DIRECT` and `GROUP` chats are encrypted on the database level using **AES-256-GCM** via a dedicated `ContentCipher` platform layer. Plaintext indexing is performed on a dedicated search vector prior to encryption to ensure full-text search remains functional without data exposure.

### Custodial Web3 Wallets
*   **User Wallet:** Automatically generated on email verification. The private key is encrypted via **AES-256-GCM** using the `WALLET_ENCRYPTION_KEY` and stored securely. Users interact with the blockchain through OAuth/JWT auth.
*   **On-Chain Observer:** A background polling loop in `transaction-service` ticks every 2 seconds to index block transfers. To guarantee **idempotency**, operations are committed using `ON CONFLICT (tx_hash, log_index) DO NOTHING`. If a transaction is already indexed, the observer skips duplicate business logic.

### Token-Gated Subscriptions
Channels can define gating rules based on minimum SUDA token balance or a subscription price. When joining a paid channel:
1.  The Messenger service requests `CheckTokenGating` via gRPC.
2.  If gating is required, Messenger invokes `ChargeChannelSubscription`.
3.  The Transactions service signs and broadcasts an on-chain transfer of the subscription fee from the user's wallet to the channel's treasury wallet.
4.  Once the transaction is confirmed, the user is granted the `SUBSCRIBER` role.

---

## 5. API & WebSocket Map

Through the API gateway, all requests use the prefix: `/api/v1/[service]`.

| Service | Prefix | Key Endpoints |
|---|---|---|
| **Messenger** | `/api/v1/messenger` | `/auth/*` (register, login, verify), `/user/me`, `/chats`, `/messages`, `/channels/*` (subscribe, posts, comments), `/friends/*` |
| **Transactions** | `/api/v1/tx` | `/wallet/me`, `/wallet/transfer` (P2P), `/wallet/channel/{channel_id}/treasury`, `/purchase/*`, `/donate` |
| **Media** | `/api/v1/media` | `/media/upload/init`, `/media/{id}/confirm`, `/media/{id}/url` |

### WebSocket Events (Delivered to client on `/ws`):
*   `NEW_MESSAGE`, `MESSAGE_EDITED`, `MESSAGE_DELETED`, `MESSAGES_READ`, `TYPING`
*   `SUDA_RECEIVED`, `SUDA_SENT`, `DONATION_SENT`, `PURCHASE_COMPLETED` (handled by the Observer)

---


## Technical Deep Dive

### 1. Private EVM Network (Hyperledger Besu)
The blockchain layer runs on **Hyperledger Besu v25.12.0**, configured as a private EVM-compatible consortium chain. 

*   **QBFT Consensus (Quorum Byzantine Fault Tolerance):** Suda uses QBFT rather than traditional Proof-of-Work or Proof-of-Stake. QBFT provides **immediate finality** — once a block is minted, it cannot be reorganized, and forks are cryptographically impossible. This is critical for custodial financial operations where ledger consistency must be absolute [2].
*   **Gas Configuration (Zero-Gas Model):** The network configuration overrides default EVM mechanics with `zeroBaseFee=true` and `--min-gas-price=0`. Since gas costs are zero, the application can sign and execute contract transactions on behalf of users without managing gas stations or exposing gas mechanics to the frontend.
*   **Performance Metrics:** The block period is set to `2 seconds` (`blockPeriodSeconds=2`), facilitating fast transaction confirmation cycles suitable for real-time mobile app feedback.

---

### 2. Cryptographic Storage & Concurrency (Nonce Management)
Since the Web3 system operates as a custodial service, securing private keys and managing concurrent transaction throughput are handled at the platform level.

*   **Encrypted Key-Rest Model:** When a user's wallet is created, the system generates a standard `secp256k1` Elliptic Curve key pair. The private key is encrypted using **AES-256-GCM** with a master key supplied via Docker Secrets/KMS. Each database record in `tx_wallets` includes a `key_version` column to facilitate zero-downtime key rotation policies in the future.
*   **Concurrency & Nonce Management:** In standard Web3 clients, sending concurrent transactions from a single address often causes race conditions, leading to "nonce too low" errors. Suda solves this using a dedicated **`NonceManager`**:
    *   It maintains an in-memory thread-safe state lock per address using Go `sync.Mutex` [2].
    *   When a transaction is initiated, the manager locks the sender's address, fetches the pending nonce from the blockchain node, signs the transaction, broadcasts it, and increments the local nonce before releasing the lock [2].

---

### 3. Resilient & Idempotent Blockchain Observer
The transaction indexing engine operates as an asynchronous background routine (`Observer`) in `transaction-service` that polls the Besu node [2].

*   **Double-Processing Protection:** To survive system crashes and node disconnects without duplicating transaction side effects (e.g., doubling a donation or minting twice), the database schema defines a unique constraint on `tx_transactions(tx_hash, log_index)`. 
*   **State Machine Lifecycle:**
    ```text
    [Client Action] ──> Write Audit Trail ──> Broadcast Tx ──> Insert tx_pending
                                                                     │
    ┌─────────────────────────── [Observer Loop] ────────────────────┘
    ▼
    Poll Blocks ──> Filter Logs ──> Insert tx_transactions (ON CONFLICT DO NOTHING)
                                             │
                       ┌─────────────────────┴─────────────────────┐
                       ▼ (If inserted = true)                      ▼ (If inserted = false)
             Trigger Domain Side Effects & WS                     Skip (Already Processed)
    ```
*   **Cursor Persistence:** The block cursor (`tx_observer_state`) is updated and persisted to PostgreSQL only *after* all logs in a batch are processed, ensuring the engine can resume from its exact prior block state after an unexpected shutdown [2].

---

### 4. Database Search Under Message Encryption
To balance user privacy with application usability, Suda implements a hybrid storage strategy for messages in `DIRECT` and `GROUP` chats:

*   **At-Rest Encryption:** Raw message content is processed through the `aesContentCipher` platform package, prepended with a version sentinel (`enc:v1:`), and written to `messenger_messages.content` as a Base64-encoded AES-256-GCM ciphertext.
*   **Searchability Vector:** Standard full-text search under encryption is solved by constructing the PostgreSQL search vector (`to_tsvector`) **prior to encryption** in the application service layer. The plaintext vector is stored in a separate indexed column (`search_vector`). This allows the PostgreSQL GIN index to execute full-text search operations normally, while the actual message content remains encrypted at rest [2].

---

## 6. Deployment & Quick Start

The deployment environment relies on a host Linux machine running **ZeroTier One** for secure access over a private VPN.

### First-Time Configuration
1.  **Generate Secrets:**
    ```bash
    task secrets-generate
    ```
    *This generates random 32-byte hexadecimal files in `secrets/` (`wallet_encryption_key.txt`, `jwt_secret.txt`, `gateway_signature_secret.txt`).*
2.  **Provide Manual Keys:**
    *   Save your EVM treasury private key (64 hex characters, no `0x` prefix) inside `secrets/eth_private_key.txt`.
    *   Download your Firebase Admin SDK service account key as `secrets/firebase-adminsdk.json`.
3.  **Environment Variables:**
    ```bash
    cp backend/docs/examples/.env.example .env
    ```

### Running the Stack
*   **Development Mode (Exposes all ports on host):**
    ```bash
    task docker-up
    ```
*   **Production Mode (Restricted ports, routes proxied via Nginx, Portainer enabled on `:9443`):**
    ```bash
    task docker-up-prod
    ```

### Troubleshooting: "SudaToken has no code" Error
If you destroy or reset the blockchain volume (`besu_data`), your previously saved contract addresses inside `secrets/contracts.env` will become stale. Follow this precise state-reset sequence:
1.  Stop the stack: `task docker-down`
2.  Delete the stale contract address file on the host: `rm -f ./secrets/contracts.env`
3.  Boot only the database and blockchain services:
    ```bash
    docker compose up -d postgres redis besu-node minio
    ```
4.  Run migrations and trigger the one-shot Solidity compiler/deployer:
    ```bash
    docker compose up migrate contracts-deploy
    ```
    *Ensure both containers exit with status `0` before continuing.*
5.  Launch the rest of the stack (force-recreating the container to reload the new `contracts.env` variables):
    ```bash
    docker compose up -d --force-recreate
    ```

