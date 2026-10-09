FROM node:20-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

ENV NEXT_PUBLIC_API_URL=http://localhost:8095/api/v1
ENV NEXT_PUBLIC_WS_URL=ws://localhost:8095/ws

EXPOSE 3000
CMD ["npm", "run", "dev"]
