cd ./traderServer
cmake .
make
cd ../mdServer
cmake .
make
cd ../
go build main.go
./main
