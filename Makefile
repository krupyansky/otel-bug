init:
	docker-compose up -d --build && go run main.go

curl-metrics:
	curl localhost:8080/metrics
