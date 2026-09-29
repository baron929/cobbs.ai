FROM --platform=$BUILDPLATFORM node:22 AS front
WORKDIR /web
COPY ./web .
RUN yarn install --frozen-lockfile --network-timeout 1000000 && yarn run build


FROM --platform=$BUILDPLATFORM node:22 AS pptx_worker
WORKDIR /pptx-worker
COPY ./tool/pptx-worker/package.json ./tool/pptx-worker/package-lock.json ./
COPY ./tool/pptx-worker/worker.mjs ./
RUN npm ci && npm run build


FROM --platform=$BUILDPLATFORM golang:1.25 AS back
ARG VERSION
ARG COMMIT
ARG BUILD_DATE
WORKDIR /go/src/cobbs.ai
COPY . .
RUN chmod +x ./build.sh
RUN VERSION="${VERSION}" COMMIT="${COMMIT}" BUILD_DATE="${BUILD_DATE}" ./build.sh


FROM alpine:latest AS standard
ARG VERSION
ARG COMMIT
ARG BUILD_DATE
LABEL MAINTAINER="https://github.com/baron929/cobbs.ai"
LABEL org.opencontainers.image.source="https://github.com/baron929/cobbs.ai"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.revision="${COMMIT}"
LABEL org.opencontainers.image.created="${BUILD_DATE}"
ARG USER=casibase
ARG TARGETOS
ARG TARGETARCH
ENV BUILDX_ARCH="${TARGETOS:-linux}_${TARGETARCH:-amd64}"

RUN apk add curl
RUN apk add nodejs
RUN apk add ca-certificates && update-ca-certificates

RUN adduser -D $USER -u 1000 \
    && mkdir logs \
    && mkdir files \
    && chown -R $USER:$USER logs \
    && chown -R $USER:$USER files

USER 1000
WORKDIR /
COPY --from=back --chown=$USER:$USER /go/src/cobbs.ai/server_${BUILDX_ARCH} ./server
COPY --from=back --chown=$USER:$USER /go/src/cobbs.ai/data ./data
COPY --from=back --chown=$USER:$USER /go/src/cobbs.ai/conf/app.conf ./conf/app.conf
COPY --from=back --chown=$USER:$USER /go/src/cobbs.ai/skills ./skills
COPY --from=front --chown=$USER:$USER /web/build ./web/build
COPY --from=pptx_worker --chown=$USER:$USER /pptx-worker/worker.bundle.mjs ./pptx-worker/worker.bundle.mjs
ENV RUNNING_IN_DOCKER=true

ENTRYPOINT ["/server"]


FROM debian:latest AS db
RUN apt update \
    && apt install -y \
        mariadb-server \
        mariadb-client \
    && rm -rf /var/lib/apt/lists/*


FROM db AS allinone
ARG VERSION
ARG COMMIT
ARG BUILD_DATE
LABEL MAINTAINER="https://github.com/baron929/cobbs.ai"
LABEL org.opencontainers.image.source="https://github.com/baron929/cobbs.ai"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.revision="${COMMIT}"
LABEL org.opencontainers.image.created="${BUILD_DATE}"
ARG TARGETOS
ARG TARGETARCH
ENV BUILDX_ARCH="${TARGETOS:-linux}_${TARGETARCH:-amd64}"

RUN apt update && apt install -y ca-certificates nodejs && update-ca-certificates

WORKDIR /
COPY --from=back /go/src/cobbs.ai/server_${BUILDX_ARCH} ./server
COPY --from=back /go/src/cobbs.ai/data ./data
COPY --from=back /go/src/cobbs.ai/docker-entrypoint.sh /docker-entrypoint.sh
COPY --from=back /go/src/cobbs.ai/conf/app.conf ./conf/app.conf
COPY --from=back /go/src/cobbs.ai/skills ./skills
COPY --from=front /web/build ./web/build
COPY --from=pptx_worker /pptx-worker/worker.bundle.mjs ./pptx-worker/worker.bundle.mjs
ENV RUNNING_IN_DOCKER=true

ENTRYPOINT ["/bin/bash"]
CMD ["/docker-entrypoint.sh"]
