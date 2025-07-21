 # ---- Builder Stage ----
# Use the official Rust image as a builder
FROM rust:1.79 as builder

# Install build essentials for dependencies like OpenSSL
RUN apt-get update && apt-get install -y build-essential pkg-config libssl-dev

WORKDIR /app

# Copy your manifests and build dependencies first to leverage Docker cache
COPY Cargo.toml Cargo.lock ./
# Create a dummy src file to build dependencies before copying your source
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
# Now that dependencies are built and cached, remove the dummy file
RUN rm -rf src

# Copy your actual source code
COPY ./src ./src

# Build your application
# Using --locked ensures it uses the versions from Cargo.lock
RUN cargo build --release --locked


# ---- Final Stage ----
# Use a small, secure base image for the final product
FROM debian:buster-slim as runtime

# Install OpenSSL libraries needed at runtime
RUN apt-get update && apt-get install -y libssl-dev ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the compiled binary from the builder stage
# Note the source path matches the binary name from Cargo.toml
COPY --from=builder /app/target/release/backend-rust /usr/local/bin/

# Set the binary as the container's command
CMD ["/usr/local/bin/backend-rust"]