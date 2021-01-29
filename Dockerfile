FROM golang:1.15.3-alpine AS builder

LABEL maintainer="Erick Zurita <erickzuria@gmail.com>"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build app
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/app /app/cmd/draid

FROM alpine:3.11.2
WORKDIR /root
COPY --from=builder /app/app .
COPY --from=builder /app/docs ./docs
RUN mkdir credentials

ENV PORT=8000
ENV SIGNING_STRING='SECRET'
ENV DATABASE_URI='mongodb://127.0.0.1:27017'
ENV SERVER_HOST=http://localhost:${PORT}
ENV CLOUD_MESSAGING_KEY=''
ENV FIREBASE_CREDENTIALS_PATH=''

EXPOSE ${PORT}
CMD [ "/root/app" ]
