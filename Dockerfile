FROM golang:1.20 AS build-env
COPY . /app
WORKDIR /app
RUN go get -d ./... && \
    CGO_ENABLED=0 GOOS=linux go build -o main .

FROM gcr.io/distroless/base
WORKDIR /app/
COPY --from=build-env /app/main .
ENV DISCORD_APP_ID="1089018173325586462"
ENV DISCORD_PUBLIC_KEY="26a32d6739233e80efbb2179fc225c86cbf2fe5c0ce7f0b8c474f185e3b0e4c2"
CMD ["./main"]
