
############################
# STEP 1 build executable binary
############################

FROM golang:1.17 as builder

WORKDIR /src
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o app -a -ldflags '-linkmode external -extldflags "-static"' .


############################
# STEP 2 build a small image
############################

FROM scratch

# Import from builder.

# Copy our static executable.
COPY --from=builder /src/app /app

# Tests only
COPY --from=builder /src/db /db

# Port on which the service will be exposed.
EXPOSE 80/tcp

ENTRYPOINT ["/app"]