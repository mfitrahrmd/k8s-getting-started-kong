FROM golang:alpine AS kong-plugin-builder

WORKDIR /app

COPY ./plugins/rbac/ ./rbac/

WORKDIR /app/rbac

RUN go mod download

RUN go build

WORKDIR /app

COPY ./plugins/go-wait ./go-wait

WORKDIR /app/go-wait

RUN go mod download

RUN go build

FROM kong:latest

USER root

WORKDIR /app

COPY --from=kong-plugin-builder /app/rbac/kong-plugin-rbac /usr/local/bin/

COPY --from=kong-plugin-builder /app/go-wait/kong-plugin-go-wait /usr/local/bin/

COPY --from=kong-plugin-builder /app/rbac/model.conf /etc/kong/casbin/
COPY --from=kong-plugin-builder /app/rbac/policy.csv /etc/kong/casbin/

# configuration files

COPY ./kong.conf /etc/kong/kong.conf
COPY ./kong.yml /kong/declarative/kong.yml

USER kong

# kong databse bootstrap

RUN kong migrations bootstrap && kong migrations up && kong migrations finish

ENTRYPOINT [ "/docker-entrypoint.sh" ]

CMD [ "kong", "docker-start" ]

EXPOSE 8000
EXPOSE 8001
EXPOSE 8002
EXPOSE 8443
EXPOSE 8444

STOPSIGNAL SIGQUIT
HEALTHCHECK --interval=10s --timeout=10s --retries=10 CMD kong health
