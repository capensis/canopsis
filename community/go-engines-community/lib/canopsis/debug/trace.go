package debug

import (
	"fmt"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

const (
	envEnable    = "CPS_DEBUG_PPROF_ENABLE"
	envWebEnable = "CPS_DEBUG_WEB_PPROF_ENABLE"
	envCPU       = "CPS_DEBUG_PPROF_CPU"
	envMemory    = "CPS_DEBUG_PPROF_MEMORY"
	envTrace     = "CPS_DEBUG_TRACE"
	// custom host:port for pprof web server
	envWebPort     = "CPS_DEBUG_WEB_PPROF_PORT"
	webDefaultPort = "6060"
	envWebHost     = "CPS_DEBUG_WEB_PPROF_HOST"
	webDefaultHost = "localhost"
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

// profilingWebPort returns port number provided by CPS_DEBUG_WEB_PPROF_PORT environment
// variable. This will be default as webDefaultPort when not provided or invalid
func profilingWebPort() (string, error) {
	val := os.Getenv(envWebPort)

	if val == "" {
		return webDefaultPort, nil
	}

	_, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		return webDefaultPort, fmt.Errorf("the environment variable %s must be 16 bit unsigned integer but got %s; use default value %s instead", envWebPort, val, webDefaultPort)
	}

	return val, nil
}

// profilingWebHost returns host provided by CPS_DEBUG_WEB_PPROF_HOST environment
// variable.
func profilingWebHost() (string, error) {
	val := os.Getenv(envWebHost)

	switch val {
	case "":
		return webDefaultHost, fmt.Errorf("the environment variable %s is empty; use default value %s instead", envWebHost, webDefaultHost)
	case "*":
		return "", nil
	}

	return val, nil
}

// getPath returns the absolute path of the file set in the provided
// environment variable.
func getPath(environmentVariable string, logger zerolog.Logger) (string, error) {
	fpath := os.Getenv(environmentVariable)
	if fpath == "" {
		logger.Info().Msgf("the environment variable %s is empty, skipping", environmentVariable)
		// Not an error, just not enabled
		return "", nil
	}

	fpath, err := filepath.Abs(fpath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %s: %w", environmentVariable, err)
	}

	return fpath, nil
}

// startCPU enables CPU profiling for the current process, and writes the CPU
// profile to the path set in the CPS_DEBUG_PPROF_CPU environment variable.
func startCPU(logger zerolog.Logger) error {
	fpath, err := getPath(envCPU, logger)
	if err != nil {
		return err
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
	fpath, err := getPath(envMemory, logger)
	if err != nil {
		return nil, err
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
func startTrace(fpath string, logger zerolog.Logger) (io.WriteCloser, error) {
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
		port, err := profilingWebPort()
		if err != nil {
			logger.Warn().Msg(err.Error())
		}

		host, err := profilingWebHost()
		if err != nil {
			logger.Warn().Msg(err.Error())
		}

		profAddr := net.JoinHostPort(host, port)
		runtime.SetBlockProfileRate(1)
		runtime.SetMutexProfileFraction(1)
		logger.Info().Str("address", profAddr).Msg("Profiling web ENABLED")
		go func() {
			srv := &http.Server{
				Addr:         profAddr,
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
			}
			err := srv.ListenAndServe()
			if err != nil {
				logger.Err(err).Msg("fail to start pprof server")
			}
		}()
	} else {
		logger.Info().Msg("Profiling DISABLED")
	}

	tracePath, err := getPath(envTrace, logger)
	if err != nil {
		logger.Err(err).Msg("Error")
	} else if tracePath != "" {
		logger.Info().Msg("Tracing ENABLED")

		writer, err := startTrace(tracePath, logger)
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
	var errorCount int

	if t.cpuStarted {
		pprof.StopCPUProfile()
		t.Logger.Info().Msg("CPU profiling stopped")
	}

	if t.memoryStarted {
		runtime.GC()
		if err := pprof.WriteHeapProfile(t.memoryWriter); err != nil {
			t.Logger.Err(err).Msg("failed to write heap profile")
			errorCount++
		} else {
			t.Logger.Info().Msg("Memory profile written successfully")
		}

		if err := t.memoryWriter.Close(); err != nil {
			t.Logger.Err(err).Msg("failed to close memory profile file")
			errorCount++
		}
	}

	if t.traceStarted {
		trace.Stop()
		t.Logger.Info().Msg("Tracing stopped")

		if err := t.traceWriter.Close(); err != nil {
			t.Logger.Err(err).Msg("failed to close trace file")
			errorCount++
		}
	}

	if errorCount > 0 {
		t.Logger.Error().Int("error_count", errorCount).Msg("Stop completed with errors")
	} else {
		t.Logger.Info().Msg("All profiling and tracing stopped successfully")
	}
}
