
SERVERTEST := 1

default: test
	
.PHONY: test
	go test .\pkg\server
