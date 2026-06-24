import { ethers } from "hardhat";
import * as fs from "fs";
import * as path from "path";

/**
 * Deploys all 5 Suda contracts to the configured network in dependency order:
 *   1. SudaToken         (no deps)
 *   2. SudaNFT           (no deps)
 *   3. SudaMarketplace   (needs SudaToken)
 *   4. SudaEscrow        (needs SudaToken)
 *   5. SudaFundraising   (needs SudaToken)
 *
 * After deployment, writes addresses to:
 *   - contracts/deployments/<chainId>.json     (machine-readable)
 *   - prints copy-paste-ready .env block       (for the project root)
 *
 * Run with:
 *   npx hardhat run scripts/deploy.ts --network suda
 */
async function main() {
  const [deployer] = await ethers.getSigners();
  const network = await ethers.provider.getNetwork();

  console.log("─".repeat(60));
  console.log("Deploying Suda contracts");
  console.log("─".repeat(60));
  console.log("Network:        ", network.name, `(chainId ${network.chainId})`);
  console.log("Deployer:       ", deployer.address);

  const balance = await ethers.provider.getBalance(deployer.address);
  console.log("Deployer balance:", ethers.formatEther(balance), "ETH");
  console.log();

  // 1. SudaToken
  console.log("[1/5] Deploying SudaToken...");
  const SudaToken = await ethers.getContractFactory("SudaToken");
  const sudaToken = await SudaToken.deploy(deployer.address);
  await sudaToken.waitForDeployment();
  const sudaTokenAddr = await sudaToken.getAddress();
  console.log("      → SudaToken at", sudaTokenAddr);

  // 2. SudaNFT
  console.log("[2/5] Deploying SudaNFT...");
  const SudaNFT = await ethers.getContractFactory("SudaNFT");
  const sudaNFT = await SudaNFT.deploy(deployer.address);
  await sudaNFT.waitForDeployment();
  const sudaNFTAddr = await sudaNFT.getAddress();
  console.log("      → SudaNFT at", sudaNFTAddr);

  // 3. SudaMarketplace (needs SudaToken)
  console.log("[3/5] Deploying SudaMarketplace...");
  const SudaMarketplace = await ethers.getContractFactory("SudaMarketplace");
  const sudaMarketplace = await SudaMarketplace.deploy(deployer.address, sudaTokenAddr);
  await sudaMarketplace.waitForDeployment();
  const sudaMarketplaceAddr = await sudaMarketplace.getAddress();
  console.log("      → SudaMarketplace at", sudaMarketplaceAddr);

  // 4. SudaEscrow (needs SudaToken)
  console.log("[4/5] Deploying SudaEscrow...");
  const SudaEscrow = await ethers.getContractFactory("SudaEscrow");
  const sudaEscrow = await SudaEscrow.deploy(deployer.address, sudaTokenAddr);
  await sudaEscrow.waitForDeployment();
  const sudaEscrowAddr = await sudaEscrow.getAddress();
  console.log("      → SudaEscrow at", sudaEscrowAddr);

  // 5. SudaFundraising (needs SudaToken)
  console.log("[5/5] Deploying SudaFundraising...");
  const SudaFundraising = await ethers.getContractFactory("SudaFundraising");
  const sudaFundraising = await SudaFundraising.deploy(deployer.address, sudaTokenAddr);
  await sudaFundraising.waitForDeployment();
  const sudaFundraisingAddr = await sudaFundraising.getAddress();
  console.log("      → SudaFundraising at", sudaFundraisingAddr);

  console.log();
  console.log("─".repeat(60));
  console.log("All contracts deployed successfully");
  console.log("─".repeat(60));

  // ── Persist addresses to JSON ───────────────────────────
  const deployment = {
    network:    network.name,
    chainId:    Number(network.chainId),
    deployer:   deployer.address,
    deployedAt: new Date().toISOString(),
    contracts: {
      SudaToken:        sudaTokenAddr,
      SudaNFT:          sudaNFTAddr,
      SudaMarketplace:  sudaMarketplaceAddr,
      SudaEscrow:       sudaEscrowAddr,
      SudaFundraising:  sudaFundraisingAddr,
    },
  };

  const deploymentsDir = path.join(__dirname, "..", "deployments");
  if (!fs.existsSync(deploymentsDir)) {
    fs.mkdirSync(deploymentsDir, { recursive: true });
  }
  const deploymentFile = path.join(deploymentsDir, `${network.chainId}.json`);
  fs.writeFileSync(deploymentFile, JSON.stringify(deployment, null, 2));
  console.log();
  console.log("Saved:", deploymentFile);

  // ── Print .env block ────────────────────────────────────
  console.log();
  console.log("─".repeat(60));
  console.log("Add (or update) these lines in your project root .env:");
  console.log("─".repeat(60));
  console.log(`SUDA_TOKEN_ADDRESS=${sudaTokenAddr}`);
  console.log(`SUDA_NFT_ADDRESS=${sudaNFTAddr}`);
  console.log(`SUDA_MARKETPLACE_ADDRESS=${sudaMarketplaceAddr}`);
  console.log(`SUDA_ESCROW_ADDRESS=${sudaEscrowAddr}`);
  console.log(`SUDA_FUNDRAISING_ADDRESS=${sudaFundraisingAddr}`);
  console.log("─".repeat(60));
}

main().catch((err) => {
  console.error(err);
  process.exitCode = 1;
});