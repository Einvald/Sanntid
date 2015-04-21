export GOPATH=$(pwd)

go install driver
go install elev_handler
go install networkModule

go run main.go
