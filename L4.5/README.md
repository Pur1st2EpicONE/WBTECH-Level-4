## L4.5

![optimization banner](assets/banner.png)

<h3 align="center">A Go project focused on performance optimization of a simple HTTP service using benchmarks, pprof, and execution tracing.</h3>

<br>

## API

A minimal HTTP endpoint for adding two integers, used as a baseline workload for performance profiling and optimization.


### POST /add

**Request:**

```json
{
  "first": 2,
  "second": 2
}
```

**Response:**

```json
{
  "result": 4
}
```

<br>

## Version History

###  add V1 
```go
func add(c *gin.Context) {

	var request request

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := request.First + request.Second

	c.JSON(http.StatusOK, response{Result: result})

}
```

<br>

###  Benchmarks

```
goos: linux
goarch: amd64
pkg: L4.5
cpu: AMD Ryzen 3 7320U with Radeon Graphics         
Benchmark-8   	  144180	      8489 ns/op	    7503 B/op	      30 allocs/op
PASS
ok  	L4.5	1.241s
```


The benchmark results indicate ~30 allocations per request and ~7.5 KB of memory usage for a trivial arithmetic operation. This suggests that execution cost is dominated by HTTP and JSON processing overhead rather than computation. This already indicates that the cost of execution is not driven by computation, but by request/response processing overhead and JSON handling in the HTTP layer.

<br>

```
      flat  flat%   sum%        cum   cum%
     230ms 13.77% 13.77%      230ms 13.77%  runtime.futex
     100ms  5.99% 19.76%      210ms 12.57%  runtime.scanobject
      90ms  5.39% 25.15%       90ms  5.39%  runtime.memclrNoHeapPointers
      80ms  4.79% 29.94%       80ms  4.79%  runtime.tgkill
      50ms  2.99% 32.93%       50ms  2.99%  runtime.(*spanSet).reset
      40ms  2.40% 35.33%       60ms  3.59%  runtime.findObject
      40ms  2.40% 37.72%       60ms  3.59%  runtime.typePointers.next
      40ms  2.40% 40.12%       40ms  2.40%  runtime.typePointers.nextFast (inline)
      30ms  1.80% 41.92%       30ms  1.80%  gcWriteBarrier
      30ms  1.80% 43.71%       30ms  1.80%  internal/runtime/syscall.Syscall6
```

The CPU profile is dominated by Go runtime internals (GC, scheduler, memory scanning) rather than application logic. This confirms that the workload is not compute-bound, but instead heavily influenced by runtime overhead caused by frequent allocations and garbage collection activity.

<br>

```
      flat  flat%   sum%        cum   cum%
  633.91MB 56.07% 56.07%   633.91MB 56.07%  bufio.NewReaderSize
   88.04MB  7.79% 63.86%    88.04MB  7.79%  encoding/json.(*Decoder).refill
   60.02MB  5.31% 69.17%    60.02MB  5.31%  net/http.Header.Clone
   58.02MB  5.13% 74.30%    58.02MB  5.13%  net/textproto.MIMEHeader.Set
   51.52MB  4.56% 78.85%    51.52MB  4.56%  github.com/gin-gonic/gin/render.writeContentType
   49.01MB  4.34% 83.19%    49.01MB  4.34%  encoding/json.NewDecoder
   48.01MB  4.25% 87.44%    74.52MB  6.59%  net/http.readRequest
   43.01MB  3.80% 91.24%    43.01MB  3.80%  net/http.(*Request).WithContext
   25.50MB  2.26% 93.50%    25.50MB  2.26%  net/http/httptest.NewRecorder
   16.50MB  1.46% 94.96%    16.50MB  1.46%  net/url.parse
```

The memory profile highlights that most allocations originate from the HTTP and JSON processing pipeline (request parsing, headers handling, and decoder buffering), not from the business logic itself. In particular, a large portion of memory is consumed by request decoding and HTTP infrastructure overhead, indicating that the service is allocation-heavy due to framework and encoding layers rather than computation.

<br>

![trace screenshot](assets/trace.png)

The highlighted goroutine shows the execution path of a single benchmark iteration. The add handler is executed deep in the call stack, after testing, httptest, and Gin routing layers. A significant part of the trace is taken by mallocgc and string/JSON-related allocations triggered by request creation and response encoding. This confirms that execution is not CPU-bound in the handler itself, but driven by allocation-heavy HTTP and JSON infrastructure, with runtime work dominating the actual business logic.
