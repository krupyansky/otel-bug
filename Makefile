down:
	docker-compose down --remove-orphans

init:
	docker-compose up -d --build && go run main.go

curl-metrics:
	curl localhost:8080/metrics

logs-otel:
	docker-compose logs otel-collector
