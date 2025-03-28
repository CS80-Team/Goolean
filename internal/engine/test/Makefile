GOCMD = go
GOTEST = $(GOCMD) test -bench=. -benchmem

bench:
	$(GOTEST)

bench_cpu:
	$(GOTEST) -cpuprofile cpu.prof

bench_mem:
	$(GOTEST) -memprofile mem.prof

cpu_profile:
	go tool pprof cpu.prof

mem_profile:
	go tool pprof mem.prof

clean:
	rm -f cpu.prof mem.prof
	$(GOCMD) clean -testcache

.PHONY: gen
gen:
	@echo "Generating queries..."
	rm -rf gen/queries.text
	go run gen/generate_queries.go