FROM node:18.19.1 AS builder

WORKDIR /app

ARG NEXT_PUBLIC_BACKEND_URL
ENV NEXT_PUBLIC_BACKEND_URL=${NEXT_PUBLIC_BACKEND_URL}

COPY package.json package-lock.json ./
RUN npm install

COPY . .

RUN npm run build

FROM node:18.19.1

WORKDIR /app

COPY --from=builder /app ./

ENV PORT=3001
EXPOSE 3001

CMD ["npm", "run", "start"]
