# Synthetic benchmark for concurrency approaches
- Conc package - https://github.com/sourcegraph/conc
- Builtin sync package - https://pkg.go.dev/sync

# Benchmarking was implemented with these conditions
- conc package callbacks do not provide context, so anonymous functions call job with context under the hood
- conc package provides both Wait and WaitAndRecover mechanisms, so we test both approaches
- builtin wait group does not provide recovery, so we defer this call

# Iterations and work emulation
Under the const block you can see the number of workers, timeout for each worker (worker timeout is increased with each iteration) and timeout for cancel function.
In example:
```
const (
	workerLimit       = 1000
	workerMinDuration = time.Nanosecond * 10
	contextTimeout    = workerMinDuration * workerLimit
	cancelTimeout     = contextTimeout - workerMinDuration
)
```

# Results
```
goos: darwin
goarch: arm64
pkg: github.com/aliocode/conc-example
BenchmarkWithConcWgNoPanics-10              1435            784373 ns/op          256647 B/op       5009 allocs/op
BenchmarkWithBuiltinWgNoPanics-10           1599            753563 ns/op          248525 B/op       4008 allocs/op
BenchmarkWithConcWgRecovered-10              126           9408403 ns/op         2244327 B/op       7197 allocs/op
BenchmarkWithBuiltinWgRecovered-10          1464            776439 ns/op          248718 B/op       4011 allocs/op
PASS
ok      github.com/aliocode/conc-example        6.286s

```
Overall, the `BenchmarkWithBuiltinWgNoPanics` and `BenchmarkWithBuiltinWgRecovered` functions performed better in terms of time and memory usage. 

The performance diff in percentage for the `BenchmarkWithBuiltinWgNoPanics` function compared to `BenchmarkWithConcWgNoPanics` is 
- 24.64% for time
- 3.2% for memory
- 21.74% for allocations

The performance diff in percentage for the `BenchmarkWithBuiltinWgRecovered` function compared to `BenchmarkWithConcWgRecovered` is 
- 14.45% for time
- 9.81% for memory
- 0.04% for allocations
