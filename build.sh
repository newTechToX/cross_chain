rm ./bin -fr
mkdir ./bin
go build -o ./bin/collector ./cmd/collector/*
go build -o ./bin/matcher ./cmd/matcher/*
go build -o ./bin/detector ./cmd/detector/*