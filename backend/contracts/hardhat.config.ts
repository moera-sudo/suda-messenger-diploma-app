import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import * as dotenv from "dotenv";

// Load .env from the project root (one level up from contracts/).
// If you keep a separate .env inside contracts/, change the path.
dotenv.config({ path: "../.env" });

const DEPLOYER_PK = process.env.ETH_PRIVATE_KEY ?? "";
const BESU_RPC_URL = process.env.BESU_RPC_URL ?? "http://localhost:8545";
const BESU_CHAIN_ID = parseInt(process.env.BESU_CHAIN_ID ?? "1337", 10);

const config: HardhatUserConfig = {
  solidity: {
    version: "0.8.28",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
      evmVersion: "cancun"
    },
  },
  networks: {
    // Local Hardhat in-memory network (default for `npx hardhat test`).
    hardhat: {
      chainId: 31337,
    },
    // Our Besu private chain.
    suda: {
      url: BESU_RPC_URL,
      chainId: BESU_CHAIN_ID,
      accounts: DEPLOYER_PK ? [DEPLOYER_PK] : [],
      // Besu with zeroBaseFee: gas price stays 0.
      gasPrice: 0,
    },
  },
};

export default config;