FROM golang:1.27.1-alpine3.24 AS build_stage
LABEL authors="artemveber"

WORKDIR /app
COPY go.mod go.sum ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /pet_restapi ./cmd/api/main.go

FROM gcr.io/distroless/base-debian13 AS runner

WORKDIR /

COPY --from=build_stage /pet_restapi /pet_restapi

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/pet_restapi"]