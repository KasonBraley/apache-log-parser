FROM node:16-alpine AS base

FROM base as deps
# Create app directory
WORKDIR /app
# Copy dependency definitions
COPY package*.json ./
# Install app dependencies
RUN yarn install
# If you are building your code for production
# RUN npm ci --only=production

FROM base
WORKDIR /app
# Copy source code
COPY . .
COPY --from=deps /app/node_modules ./node_modules/

RUN yarn build

EXPOSE 3000

# Specify what container executes
CMD ["yarn", "start"]
