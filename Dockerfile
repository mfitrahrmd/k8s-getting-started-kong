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

WORKDIR /app

COPY ./plugins/go-log ./go-log

WORKDIR /app/go-log

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build

FROM kong:latest

USER root

RUN mkdir -p /go-plugins /go-log

COPY --from=kong-plugin-builder /app/rbac/kong-plugin-rbac /go-plugins/

COPY --from=kong-plugin-builder /app/go-wait/kong-plugin-go-wait /go-plugins/

COPY --from=kong-plugin-builder /app/go-log/kong-plugin-go-log /go-plugins/

COPY --from=kong-plugin-builder /app/rbac/model.conf /etc/kong/casbin/
COPY --from=kong-plugin-builder /app/rbac/policy.csv /etc/kong/casbin/

RUN chmod -R 755 /go-plugins /go-log

RUN chown -R kong:kong /go-log

# configuration files

COPY ./kong.conf /etc/kong/kong.conf
COPY ./kong.yml /kong/declarative/kong.yml

USER kong

ENTRYPOINT [ "/docker-entrypoint.sh" ]

EXPOSE 8000
EXPOSE 8001
EXPOSE 8002
EXPOSE 8443
EXPOSE 8444

STOPSIGNAL SIGQUIT
HEALTHCHECK --interval=10s --timeout=10s --retries=10 CMD kong health

CMD [ "kong", "docker-start" ]