

up: tidy
	go run .
tidy:
	go mod tidy


heyset:
	curl -X POST http://127.0.0.1:8081/promote/iphone/20

hey50: heyset
	hey -n 100 -c 20 http://127.0.0.1:8081/promote/iphone


hey200: heyset
	hey -n 1000 -c 300 http://127.0.0.1:8081/promote/iphone