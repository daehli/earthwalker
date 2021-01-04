FROM node:14-alpine

WORKDIR /usr/src/app

COPY rollup.config.js ./
COPY package*.json ./

RUN npm install

COPY ./src ./src

RUN npm run build

COPY ./public ./public

CMD [ "npm", "start" ]