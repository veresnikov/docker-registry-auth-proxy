FROM debian:11-slim
RUN apt-get update && \
    apt-get install -y --no-install-suggests --no-install-recommends ca-certificates && \
    apt-get install -y curl && \
    apt-get clean && \
    groupadd -g 1001 serviceuser && \
    useradd -u 1001 -r -g 1001 -s /sbin/nologin -c "go service user" serviceuser

USER serviceuser

EXPOSE 8000

ADD conf/example.config.json /app/.authproxy/config.json
ADD bin/authproxy /app/authproxy

ENTRYPOINT ["/app/authproxy"]
CMD ["serve"]