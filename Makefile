migrate:
	 migrate -database 'mysql://root:password@tcp(127.0.0.1:32771)/ledger_auth' -path sql/migrations up

migrate-down:
	 migrate -database 'mysql://root:password@tcp(127.0.0.1:32771)/ledger_auth' -path sql/migrations down
