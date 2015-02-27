export GOPATH=$(pwd)

go install driver
go install elev_handler

go run main.go
