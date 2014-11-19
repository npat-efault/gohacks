
SRC = cirq.go

build:: ${SRC}
	go build ${BF}

install:: ${SRC}
	go install ${IF}

vet:: ${SRC}
	go vet ${VF}

test:: ${SRC}
	go test ${TF}

clean::
	rm -f ${SRC}
	rm -f *~
	rm -f *.out
	go clean ${CF}

cirq.go : cirq.gox
	rm -f cirq.go
	./cirq_gen.sh cirq CQ New 'interface{}' > cirq.go.tmp
	mv cirq.go.tmp cirq.go
	chmod -w cirq.go
