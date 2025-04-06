## testing
unit test:
go test ./...

go test -v ./<pathToPackage>
go test -v ./testex/arithmetic

benchmark:
go test <pathToPackage> -bench=<pattern> ". or matching pattern"
go test ./testex/arithmetic -bench=. -benchmem

go test -bench .
go test -bench . -benchmem