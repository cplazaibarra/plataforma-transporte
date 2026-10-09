FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY server /app/server
COPY simulator /app/simulator
COPY migrations /app/migrations

ENV PORT=8095
ENV DB_HOST=lastmile_postgres
ENV DB_PORT=5432
ENV DB_USER=lastmile_user
ENV DB_PASSWORD=lastmile_password
ENV DB_NAME=lastmile_db

EXPOSE 8095
CMD ["/app/server"]
