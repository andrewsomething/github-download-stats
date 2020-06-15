FROM golang:alpine AS build

ARG COMMIT
ARG VERSION

RUN apk --no-cache add build-base gcc
ADD . /src
RUN cd /src && \
	go build -mod=vendor \
	-ldflags="-X main.version=${VERSION} -X main.commit=${COMMIT}"


FROM alpine
WORKDIR /app
COPY --from=build /src/github-download-stats /app/
ENTRYPOINT ["./github-download-stats"]