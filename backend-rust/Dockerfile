FROM rust:latest as builder
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/app/target \
    --mount=type=cache,target=/usr/local/cargo/registry \
    cargo build --release

FROM debian:buster-slim
COPY --from=builder /app/target/release/backend-rust /usr/local/bin/app
EXPOSE 8081
CMD ["/backend-rust"]
