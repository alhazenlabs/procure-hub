rm runner

cd backend && ./pkg/swag init && cd -
go mod tidy
export SERVICE_NAME=procure-hub
export INSECURE_MODE=true 
export OTEL_EXPORTER_OTLP_ENDPOINT=0.0.0.0:4317
go build -o runner  backend/cmd/*

# docker run --name procure-hub-postgres -e POSTGRES_PASSWORD=admin@123 -p 5432:5432 -d postgres
./runner