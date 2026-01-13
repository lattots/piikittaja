FROM node:20-alpine AS builder
WORKDIR /src

RUN npm init -y && npm install --save-dev esbuild

COPY . .

RUN npx esbuild js/UserTable.ts --bundle --outfile=index.js --minify --format=esm

RUN mkdir /dist && \
	mv index.js /dist/index.js && \
	cp html/index.html /dist/index.html && \
	cp css/style.css /dist/style.css

FROM alpine:latest
WORKDIR /app

COPY --from=builder /dist /app/holding

CMD sh -c "cp -rp /app/holding/. /app/build/ && echo 'Build assets copied to volume'"
