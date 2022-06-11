FROM alpine:latest

LABEL maintainer="Cédric Roussel <cedric.roussel@lilo.org>"

RUN apk add --no-cache libc6-compat

USER nobody

COPY server /app/

EXPOSE 3000

ARG GIT_HASH
ENV GIT_HASH=$GIT_HASH

CMD [ "/app/server" ]
