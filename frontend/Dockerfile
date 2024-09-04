# for production

FROM node:18.17.0-alpine3.18 AS builder

WORKDIR /app

COPY package.json ./

RUN yarn install

COPY . .

RUN yarn build

FROM nginx:1.21.4-alpine

COPY --from=builder /app/out /usr/share/nginx/html

RUN rm /etc/nginx/conf.d/default.conf
COPY --from=builder /app/nginx/default.conf /etc/nginx/conf.d/default.conf

# RUN rm /etc/nginx/nginx.conf

# COPY --from=builder /app/nginx/default.conf /etc/nginx/conf.d/default.conf

# COPY --from=builder /app/nginx/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
