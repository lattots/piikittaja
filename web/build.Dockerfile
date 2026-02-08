FROM node:20-alpine AS builder
WORKDIR /src

RUN npm init -y && npm install --save-dev esbuild chart.js

COPY . .

RUN npx esbuild js/index.ts --bundle --outfile=index.js --minify --format=esm

RUN mkdir /dist && \
	mkdir /dist/assets && \
	mv index.js /dist/index.js && \
	cp html/index.html /dist/index.html && \
	cp html/login.html /dist/login.html && \
	cp assets/* /dist/assets && \
	cp css/style.css /dist/style.css

FROM alpine:latest
WORKDIR /app

COPY --from=builder /dist /app/holding

CMD sh -c "cp -rp /app/holding/. /app/build/ && echo 'Build assets copied to volume'"
