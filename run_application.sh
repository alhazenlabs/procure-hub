rm runner

go mod tidy
export SERVICE_NAME=procure-hub
export INSECURE_MODE=true 
export OTEL_EXPORTER_OTLP_ENDPOINT=0.0.0.0:4317
cd backend && ./pkg/swag init && cd -
go build -o runner  backend/cmd/*

./runner