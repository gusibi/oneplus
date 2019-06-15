module html-server

go 1.12

require (
	github.com/aws/aws-lambda-go v1.11.1
	github.com/awslabs/aws-lambda-go-api-proxy v0.3.0
	github.com/go-openapi/loads v0.19.0
	github.com/go-sql-driver/mysql v1.4.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/julienschmidt/httprouter v1.2.0
	github.com/tencentyun/scf-go-lib v0.0.0-20190226080708-b26c1e4b5808
	golang.org/x/net v0.0.0-20181005035420-146acd28ed58
	golang.org/x/text v0.3.0
)

replace (
	golang.org/x/net v0.0.0-20181005035420-146acd28ed58 => github.com/golang/net v0.0.0-20181005035420-146acd28ed58
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0

)
