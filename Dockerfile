FROM golang:1.20-alpine AS build-env

RUN mkdir /app
WORKDIR /app

COPY . .
RUN go mod tidy && go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o validations

EXPOSE 8091

FROM scratch
COPY ./ ./
COPY --from=build-env /app/validations validations
CMD ["./validations"]