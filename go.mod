module github.com/bradenrayhorn/ledger-auth

go 1.16

require (
	github.com/bradenrayhorn/ledger-protos v0.0.0-20210516151140-f082cd637e1e
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/zap v0.0.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis/v8 v8.8.2
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.2.0
	github.com/jmoiron/sqlx v1.3.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
	google.golang.org/genproto v0.0.0-20210513213006-bf773b8c8384 // indirect
	google.golang.org/grpc v1.37.1
)
