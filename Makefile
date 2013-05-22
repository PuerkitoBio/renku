bench:
	renku -d testdata/ &
	siege -t 30s -b -H "Accept-Encoding: *" --log="./bench/siege.log" http://localhost:9000/my.md

.PHONY: bench

