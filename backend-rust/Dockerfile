FROM rust:1.72 as builder
WORKDIR /app
COPY . .
RUN cargo build --release

FROM debian:buster-slim
COPY --from=builder /app/target/release/backend-rust /usr/local/bin/app
EXPOSE 8081
CMD ["app"]
