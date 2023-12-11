FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY keycloak-operator /bin
USER 65532:65532
ENTRYPOINT ["/bin/keycloak-operator"]