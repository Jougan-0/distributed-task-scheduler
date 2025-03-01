# Use Node.js for building the Next.js app
FROM node:18.19.1 AS builder

WORKDIR /app

# Copy package.json and install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy project files
COPY . .

# Build the Next.js app
RUN npm run build

# Use a minimal base image to serve the UI
FROM node:18.19.1

WORKDIR /app

# Copy built files from the builder
COPY --from=builder /app ./

# Set environment variables
ENV PORT=3001
EXPOSE 3001

# Start the Next.js server
CMD ["npm", "run", "start"]
