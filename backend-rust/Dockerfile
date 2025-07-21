FROM rust:1.72 as builder
WORKDIR /app
COPY . .
RUN cargo build --release

FROM debian:buster-slim
COPY --from=builder /app/target/release/your_binary_name /usr/local/bin/app
EXPOSE 8081
CMD ["app"]
