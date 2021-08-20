migrate:
	 migrate -database 'postgres://postgres:password@127.0.0.1:${POSTGRES_PORT}/ledger_auth?sslmode=disable' -path sql/migrations up

migrate-down:
	 migrate -database 'postgres://postgres:password@127.0.0.1:${POSTGRES_PORT}/ledger_auth?sslmode=disable' -path sql/migrations down

test:
	docker-compose -f docker-compose.test.yml up --abort-on-container-exit --build
	docker-compose -f docker-compose.test.yml down --volumes

report:
	go tool cover -html=./reports/coverage.txt
