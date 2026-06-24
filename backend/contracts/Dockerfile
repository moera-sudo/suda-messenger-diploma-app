FROM node:20-alpine
WORKDIR /work

# Install deps once at image build time
COPY package*.json ./
RUN npm ci

# Source (deploy.ts, hardhat.config.ts, etc.)
COPY . .

# Idempotent deploy wrapper
RUN chmod +x deploy-init.sh
CMD ["./deploy-init.sh"]
