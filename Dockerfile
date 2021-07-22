FROM golang:latest

WORKDIR D:/GIT/sender
#RUN ls

# Fetch dependencies
COPY go.mod ./
RUN go mod download

# Build
COPY ./ ./
RUN go build -o sender cmd/main.go


# Create final image
CMD ["./sender"]