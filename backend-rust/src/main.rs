use actix_cors::Cors;
use actix_web::{post, web, App, HttpResponse, HttpServer, Responder};
use serde::{Deserialize, Serialize};
use aes::Aes128;
use block_modes::{BlockMode, Cbc};
use block_modes::block_padding::Pkcs7;
use hex_literal::hex;
use base64::{encode};

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
    let key = b"verysecretkey123";  
    let iv = b"uniqueinitvector";   

    let cipher = Cbc::<Aes128, Pkcs7>::new_from_slices(key, iv).unwrap();
    let encrypted_data = cipher.encrypt_vec(data.password.as_bytes());

    
    let encrypted_base64 = encode(&encrypted_data);

    HttpResponse::Ok().json(EncryptResponse { encrypted_password: encrypted_base64 })
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .wrap(Cors::permissive())  
            .service(encrypt)
    })
    .bind(("0.0.0.0",8081))?
    .run()
    .await
}
