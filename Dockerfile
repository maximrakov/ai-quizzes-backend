FROM golang:1.24.7-alpine3.22 AS builder

WORKDIR /opt/ai-quizzes-backend

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o ./bin/ai-quizzes-backend ./cmd/


FROM alpine:3.22

RUN apk add --no-cache ca-certificates
RUN addgroup -S app && adduser -S app -G app

WORKDIR /opt/ai-quizzes-backend
COPY --from=builder /opt/ai-quizzes-backend/bin/ai-quizzes-backend ./ai-quizzes-backend
EXPOSE 8080

USER app

CMD ["./ai-quizzes-backend"]