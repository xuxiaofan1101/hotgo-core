#!/usr/bin/env bash
set -euo pipefail

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PASSWORD="${CERT_PASSWORD:-changeit}"

cd "$BASE_DIR"

rm -f ca.crt ca.key ca.srl client.crt client.csr client.key client.keystore.p12 client.truststore.p12 \
  server.crt server.csr server.ext server.key server.keystore.p12 server.truststore.p12

openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 \
  -subj "/CN=hotgo-kafka-local-ca" \
  -out ca.crt

cat > server.ext <<'EOF'
subjectAltName=DNS:kafka,DNS:localhost,IP:127.0.0.1
extendedKeyUsage=serverAuth
EOF

openssl genrsa -out server.key 2048
openssl req -new -key server.key -subj "/CN=kafka" -out server.csr
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out server.crt -days 3650 -sha256 -extfile server.ext
openssl pkcs12 -export -in server.crt -inkey server.key -certfile ca.crt \
  -name kafka-server -out server.keystore.p12 -passout "pass:${PASSWORD}"
keytool -importcert -noprompt -alias hotgo-kafka-ca -file ca.crt \
  -keystore server.truststore.p12 -storetype PKCS12 -storepass "$PASSWORD"

openssl genrsa -out client.key 2048
openssl req -new -key client.key -subj "/CN=hotgo" -out client.csr
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out client.crt -days 3650 -sha256
openssl pkcs12 -export -in client.crt -inkey client.key -certfile ca.crt \
  -name hotgo-client -out client.keystore.p12 -passout "pass:${PASSWORD}"
keytool -importcert -noprompt -alias hotgo-kafka-ca -file ca.crt \
  -keystore client.truststore.p12 -storetype PKCS12 -storepass "$PASSWORD"

echo "Kafka TLS certificates generated in $BASE_DIR"
