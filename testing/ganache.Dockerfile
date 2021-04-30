FROM node:alpine

RUN npm install -g ganache-cli

ENTRYPOINT [ "ganache-cli" ]