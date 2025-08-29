# FROM golang:1.23.0-bookworm

# COPY . /usr/src/app

# WORKDIR /usr/src/app

# RUN go build -o main .

# EXPOSE 8000

# CMD ["./main"]

# docker build -t maze-conquest-api .
# docker run --name maze-conquest-api --publish 8000:8000 maze-conquest-api MODE=prod

FROM golang:1.23.0-bookworm AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o /myapp .

# FROM scratch
# FROM gcr.io/distroless/base-debian10
FROM alpine:latest

COPY --from=build /myapp /myapp
COPY --from=build /app/public /public

# While in development and want to try docker, then use keys.json
COPY --from=build /app/build-keys.json /build-keys.json
# COPY --from=build /app/keys.json /keys.json

# Set the MODE environment variable
ENV MODE=prod

ENTRYPOINT ["/myapp"]

# gcloud builds submit --tag gcr.io/maze-conquest-api/maze-conquest-api .