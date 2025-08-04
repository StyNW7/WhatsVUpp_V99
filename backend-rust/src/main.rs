use actix_cors::Cors;
use actix_web::{post, web, App, HttpResponse, HttpServer, Responder};
use serde::{Deserialize, Serialize};
use aes::Aes128;
use block_modes::{BlockMode, Cbc};
use block_modes::block_padding::Pkcs7;
use base64::encode;
use actix_web_prom::PrometheusMetricsBuilder;
use prometheus::{Encoder, TextEncoder, register_gauge, register_int_gauge, IntGauge, Gauge};
use lazy_static::lazy_static;
use sysinfo::System;

lazy_static! {
    static ref MEMORY_USAGE: IntGauge = register_int_gauge!(
        "app_memory_usage_bytes",
        "Resident memory size in bytes"
    ).unwrap();

    static ref START_TIME: Gauge = register_gauge!(
        "app_start_time_seconds",
        "App start time in seconds since Unix epoch"
    ).unwrap();
}

#[derive(Deserialize)]
struct EncryptRequest {
    password: String,
}

#[derive(Serialize)]
struct EncryptResponse {
    encrypted_password: String,
}

#[post("/encrypt")]
async fn encrypt(data: web::Json<EncryptRequest>) -> impl Responder {
    let key = b"verysecretkey123";  // 16 bytes
    let iv = b"uniqueinitvector";   // 16 bytes

    let cipher = Cbc::<Aes128, Pkcs7>::new_from_slices(key, iv).unwrap();
    let encrypted_data = cipher.encrypt_vec(data.password.as_bytes());

    let encrypted_base64 = encode(&encrypted_data);

    HttpResponse::Ok().json(EncryptResponse { encrypted_password: encrypted_base64 })
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Build Prometheus middleware
    let start_time = std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .unwrap()
        .as_secs_f64();
    START_TIME.set(start_time);

    let prometheus = PrometheusMetricsBuilder::new("api")
        .endpoint("/metrics") // Expose Prometheus metrics at `/metrics`
        .build()
        .unwrap();

    let mut sys = System::new_all();
    actix_web::rt::spawn(async move {
        loop {
            sys.refresh_memory();
            MEMORY_USAGE.set(sys.used_memory() as i64 * 1024); // bytes
            tokio::time::sleep(tokio::time::Duration::from_secs(10)).await;
        }
    });

    HttpServer::new(move || {
        App::new()
            .wrap(prometheus.clone())     // Add Prometheus middleware
            .wrap(Cors::permissive())     // Enable CORS
            .service(encrypt)             // Add `/encrypt` endpoint
    })
    .bind(("0.0.0.0", 8081))?
    .run()
    .await
}
