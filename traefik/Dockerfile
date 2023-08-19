FROM traefik:v2.2.11

RUN mkdir -p /etc/traefik/acme \
  && touch /etc/traefik/acme/acme.json \
  && chmod 600 /etc/traefik/acme/acme.json

COPY ./traefik/traefik.yaml /etc/traefik