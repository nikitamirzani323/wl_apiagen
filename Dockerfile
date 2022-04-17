FROM golang:alpine AS apimaster

WORKDIR /appbuilds

COPY . .

RUN go mod tidy
RUN go build -o binary


FROM alpine:latest as clientapirelease
WORKDIR /app
RUN apk add tzdata
COPY --from=apimaster /appbuilds/binary .
COPY --from=apimaster /appbuilds/.env /app/.env
ENV PORT=1013
ENV DB_USER="admindb"
ENV DB_PASS="asd123QWE"
ENV DB_HOST="128.199.124.131"
ENV DB_PORT="5432"
ENV DB_NAME="admindb"
ENV DB_SCHEMA="db_wl"
ENV DB_DRIVER="postgres"
ENV DB_REDIS_HOST="128.199.124.131"
ENV DB_REDIS_PORT="6379"
ENV DB_REDIS_PASSWORD="asdQWE123!@#"
ENV DB_REDIS_NAME="7"
ENV JWT_SECRET_KEY="secretbabaliong"
ENV JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT="1440"
ENV TZ=Asia/Jakarta

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENTRYPOINT [ "./binary" ]