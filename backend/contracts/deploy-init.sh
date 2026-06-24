#!/bin/sh
set -e

ENV_OUT=/secrets/contracts.env
CHAIN_ID="${BESU_CHAIN_ID:-1337}"

# Idempotency: skip if env file already populated
if [ -f "$ENV_OUT" ] && grep -q "SUDA_TOKEN_ADDRESS=0x" "$ENV_OUT"; then
  echo "[deploy-init] $ENV_OUT exists — skipping deploy."
  echo "[deploy-init] Delete it (or 'docker compose down -v') to force re-deploy."
  exit 0
fi

echo "[deploy-init] Deploying contracts to $BESU_RPC_URL (chain $CHAIN_ID)..."
npx hardhat run scripts/deploy.ts --network suda

# Read JSON output and emit KEY=VAL env-file
node -e "
  const d = require('/work/deployments/${CHAIN_ID}.json');
  const c = d.contracts || d;
  const lines = [
    'SUDA_TOKEN_ADDRESS='        + c.SudaToken,
    'SUDA_NFT_ADDRESS='          + c.SudaNFT,
    'SUDA_MARKETPLACE_ADDRESS='  + c.SudaMarketplace,
    'SUDA_ESCROW_ADDRESS='       + c.SudaEscrow,
    'SUDA_FUNDRAISING_ADDRESS='  + c.SudaFundraising,
  ].join('\n') + '\n';
  require('fs').writeFileSync('${ENV_OUT}', lines);
"

echo "[deploy-init] wrote $ENV_OUT:"
cat "$ENV_OUT"
