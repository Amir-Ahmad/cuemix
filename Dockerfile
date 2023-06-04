FROM cuelang/cue:0.5.0 as cue

FROM alpine:3.18.0

ARG TARGETARCH

COPY "dist/cuemix_linux_${TARGETARCH}" /usr/bin/cuemix
COPY --from=cue /usr/bin/cue /usr/bin/cue

RUN apk add --no-cache git

WORKDIR /app

ENTRYPOINT ["cuemix"]
