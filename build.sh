rm ./bin -fr
mkdir ./bin
go build -o ./bin/collector ./cmd/collector/main.go
go build -o ./bin/matcher ./cmd/matcher/*
go build -o ./bin/detector ./cmd/detector/*
go build -o ./bin/pro ./cmd/collector/pro.go