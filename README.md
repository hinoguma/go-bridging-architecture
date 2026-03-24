# go-connectable-architecture




# Folder Structure

```text
go-connectable-architecture/
├── .gitattribute
├── .gitignore
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   └── lambda_apigw/
│       └── main.go
├── config/
└── internal/
    ├── appinfra/
    │   ├── aws/
    │   └── postgres/
    ├── appinout/
    │   └── lambdaapigw/
    ├── appinout_applogic_connector/
    │   └── lambdaapigw_connector/
    ├── applogic/
    │   ├── domain/
    │   │   ├── model/
    │   │   ├── repository/
    │   │   └── service/
    │   └── usecase/
    ├── applogic_appinfra_connector/
    ├── crosscutting/
    └── setup/
        └── registry/
