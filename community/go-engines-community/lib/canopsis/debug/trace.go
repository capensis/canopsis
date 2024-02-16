package debug

import (
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"

	"github.com/rs/zerolog"
)

const (
	envEnable    = "CPS_DEBUG_PPROF_ENABLE"
	envWebEnable = "CPS_DEBUG_WEB_PPROF_ENABLE"
	envCPU       = "CPS_DEBUG_PPROF_CPU"
	envMemory    = "CPS_DEBUG_PPROF_MEMORY"
	envTrace     = "CPS_DEBUG_TRACE"
)

// Trace is a struct containing informations about the traces and profiles that
// have been started, so that they can be stopped correctly by the Stop method.
type Trace struct {
	cpuStarted    bool
	memoryStarted bool
	memoryWriter  io.WriteCloser
	traceStarted  bool
	traceWriter   io.WriteCloser
	Logger        zerolog.Logger
}

// profilingEnabled returns true if the CPU_DEBUG_PPROF_ENABLE environment
// variable is equal to "1", in which case the profiling should be enabled.
func profilingEnabled() bool {
	return os.Getenv(envEnable) == "1"
}

// profilingWebEnabled returns true if the CPS_DEBUG_WEB_PPROF_ENABLE environment
// variable is equal to "1", in which case the profiling server should be enabled.
func profilingWebEnabled() bool {
	return os.Getenv(envWebEnable) == "1"
}

// getPath returns the absolute path of the file set in the provided
// environment variable.
func getPath(environmentVariable string) string {
	fpath := os.Getenv(environmentVariable)
	if fpath == "" {
		return ""
	}

	fpath, _ = filepath.Abs(fpath)

	return fpath
}

// startCPU enables CPU profiling for the current process, and writes the CPU
// profile to the path set in the CPS_DEBUG_PPROF_CPU environment variable.
func startCPU(logger zerolog.Logger) error {
	fpath := getPath(envCPU)
	if fpath == "" {
		return fmt.Errorf("the environment variable %s is empty", envCPU)
	}

	fh, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed to create CPU profile: %w", err)
	}

	err = pprof.StartCPUProfile(fh)
	if err != nil {
		return fmt.Errorf("failed to start CPU profile: %w", err)
	}

	logger.Info().Msgf("CPU profiling enabled on file: %s", fpath)

	return nil
}

// startMemory enables memory profiling for the current process, and writes the
// heap profile to the path set in the CPS_DEBUG_PPROF_MEMORY environment
// variable.
func startMemory(logger zerolog.Logger) (io.WriteCloser, error) {
	fpath := getPath(envMemory)
	if fpath == "" {
		return nil, fmt.Errorf("the environment variable %s is empty", envMemory)
	}

	fh, err := os.Create(fpath)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory profile: %w", err)
	}

	logger.Info().Msgf("Memory profiling enabled on file: %s", fpath)

	return fh, nil
}

// startTrace enabled tracing for the current process, and writes the trace to
// the path set in the CPS_DEBUG_TRACE environment variable.
func startTrace(logger zerolog.Logger) (io.WriteCloser, error) {
	fpath := getPath(envTrace)
	if fpath == "" {
		return nil, fmt.Errorf("the environment variable %s is empty",
			envMemory)
	}

	fh, err := os.Create(fpath)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace file: %w", err)
	}

	err = trace.Start(fh)
	if err != nil {
		return nil, fmt.Errorf("failed to start trace: %w", err)
	}

	logger.Info().Msgf("Tracing enabled on file: %s", fpath)

	return fh, nil
}

// Start starts CPU and memory profiling (depending on the values of the
// CPS_DEBUG_* environment variables).
// It returns a Trace, that should be stopped with the Stop method so that the
// profiles are written.
func Start(logger zerolog.Logger) Trace {
	t := Trace{Logger: logger}

	if profilingEnabled() {
		logger.Info().Msg("Profiling ENABLED")

		cpuerr := startCPU(logger)
		writer, memerr := startMemory(logger)

		if cpuerr != nil {
			logger.Err(cpuerr).Msg("Error")
		} else {
			t.cpuStarted = true
		}

		if memerr != nil {
			logger.Err(memerr).Msg("Error")
		} else {
			t.memoryStarted = true
			t.memoryWriter = writer
		}
	} else if profilingWebEnabled() {
		runtime.SetBlockProfileRate(1)
		runtime.SetMutexProfileFraction(1)
		logger.Info().Msg("Profiling web ENABLED")
		go func() {
			err := http.ListenAndServe("localhost:6060", nil) //nolint:gosec
			if err != nil {
				logger.Err(err).Msg("fail to start pprof server")
			}
		}()
	} else {
		logger.Info().Msg("Profiling DISABLED")
	}

	tracePath := getPath(envTrace)
	if tracePath != "" {
		logger.Info().Msg("Tracing ENABLED")

		writer, err := startTrace(logger)
		if err != nil {
			logger.Err(err).Msg("Error")
		} else {
			t.traceStarted = true
			t.traceWriter = writer
		}
	} else {
		logger.Info().Msg("Tracing DISABLED")
	}

	return t
}

// Stop stops the traces and profiles, and writes them.
func (t Trace) Stop() {
	if t.cpuStarted {
		pprof.StopCPUProfile()
	}

	if t.memoryStarted {
		runtime.GC()
		if err := pprof.WriteHeapProfile(t.memoryWriter); err != nil {
			t.Logger.Err(err).Msg("mem error: cannot write heap profile")
		}
		t.memoryWriter.Close()
	}

	if t.traceStarted {
		trace.Stop()
		t.traceWriter.Close()
	}
}
