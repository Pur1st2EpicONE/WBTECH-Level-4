package handler

import (
	"fmt"
	"io"
	"runtime"

	"runtime/debug"

	"github.com/gin-gonic/gin"
)

var gcPercent = func() int {
	p := debug.SetGCPercent(-1)
	debug.SetGCPercent(p)
	return p
}()

// Metrics exposes runtime and GC metrics in Prometheus text format.
func Metrics(c *gin.Context) {
	c.Header("Content-Type", `text/plain; version=0.0.4; charset=utf-8`)
	writePrometheusMetrics(c.Writer)
}

// writePrometheusMetrics collects and writes Go runtime metrics.
func writePrometheusMetrics(w io.Writer) {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	write := func(format string, args ...any) {
		_, _ = fmt.Fprintf(w, format, args...)
	}

	write("go_goroutines %d\n", runtime.NumGoroutine())

	write("go_memstats_alloc_bytes %d\n", m.Alloc)
	write("go_memstats_total_alloc_bytes %d\n", m.TotalAlloc)
	write("go_memstats_sys_bytes %d\n", m.Sys)

	write("go_memstats_heap_alloc_bytes %d\n", m.HeapAlloc)
	write("go_memstats_heap_inuse_bytes %d\n", m.HeapInuse)
	write("go_memstats_heap_idle_bytes %d\n", m.HeapIdle)
	write("go_memstats_heap_objects %d\n", m.HeapObjects)

	write("go_memstats_mallocs_total %d\n", m.Mallocs)
	write("go_memstats_frees_total %d\n", m.Frees)

	write("go_memstats_stack_inuse_bytes %d\n", m.StackInuse)

	write("go_memstats_num_gc %d\n", m.NumGC)
	write("go_memstats_num_forced_gc %d\n", m.NumForcedGC)

	write("go_memstats_last_gc_time_seconds %.3f\n", float64(m.LastGC)/1e9)
	write("go_memstats_pause_total_seconds %.9f\n", float64(m.PauseTotalNs)/1e9)
	write("go_memstats_gc_cpu_fraction %f\n", m.GCCPUFraction)

	write("go_gc_percent %d\n", gcPercent)

}
