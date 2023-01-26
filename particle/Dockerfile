FROM golang:1.18-alpine AS builder

# Install git
RUN apk update && apk add --no-cache git make

# Where our file will be in the docker container 
WORKDIR /opt/app

# Copy the source from the current directory to the working Directory inside the container 
# Source also contains go.mod and go.sum which are dependency files
COPY go.mod .
COPY go.sum .

# Get Dependencies
RUN go mod download

COPY . .

RUN make

FROM alpine:latest AS runner
RUN apk --no-cache add ca-certificates bash curl
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip


FROM runner AS app
WORKDIR /app
COPY --from=builder /opt/app/bin/app .
CMD [ "./app" ]


