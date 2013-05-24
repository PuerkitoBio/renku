bench:
	renku -d testdata/ -o /dev/null &
	sleep 1s
	siege -t 30s -b -H "Accept-Encoding: *" --log="./bench/siege.log" http://localhost:9000/my.md
	pkill -9 renku

.PHONY: bench

