# Build image
##################################################
FROM golang:1.16-alpine3.13 AS build

WORKDIR /app

# Download dependencies and cache them
COPY go.* ./
RUN go mod download

# Build the program
# TODO: Consider using CGO_ENABLED=0 to avoid dependency on libc
COPY . .
RUN go build \
	-v \
	-ldflags "-extldflags '-static'" \
	-o bin/developer-bot \
	./main.go

# Runtime image
##################################################
FROM alpine:3.13

# Metadata
LABEL maintainer="denniskr@stud.ntnu.no"
EXPOSE 80

# Copy final binariy from build image,
# and run it by default when the container starts
COPY --from=build /app/bin/developer-bot /app/bin/developer-bot
CMD ["/app/bin/developer-bot"]
