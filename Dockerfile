FROM --platform=${BUILDPLATFORM} golang:1.25 AS build
ARG TARGETOS
ARG TARGETARCH

WORKDIR /build
COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o /dist/lima_ddns

FROM gcr.io/distroless/static-debian12
COPY --from=build /dist/lima_ddns /lima_ddns
ENTRYPOINT ["/lima_ddns"]
