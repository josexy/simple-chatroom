FROM node:16-alpine AS builder

WORKDIR /web
COPY . .

RUN npm install -g cnpm --registry=https://registry.npm.taobao.org
RUN cnpm install || npm install
RUN npm run build

FROM nginx:alpine
ENV TZ=Asia/Shanghai

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

COPY --from=builder /web/build /usr/share/nginx/html

EXPOSE 80 443
