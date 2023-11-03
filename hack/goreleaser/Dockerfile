FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY keycloak-operator /usr/local/bin
USER 65532:65532
ENTRYPOINT ["/usr/local/bin/keycloak-operator"]