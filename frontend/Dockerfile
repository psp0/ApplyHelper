# for production

FROM node:18.17.0-alpine3.18 as builder

WORKDIR /app

COPY package.json ./

RUN yarn install

COPY . .

RUN yarn build

FROM nginx:1.21.4-alpine

COPY --from=builder /app/out /usr/share/nginx/html

RUN rm /etc/nginx/conf.d/default.conf
COPY --from=builder /app/nginx/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
