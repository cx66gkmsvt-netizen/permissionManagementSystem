package trace

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type contextKey struct{}

var traceKey = contextKey{}

// Step represents a single execution step
type Step struct {
	Time    string `json:"time"`
	Step    string `json:"step"`
	Details string `json:"details,omitempty"`
}

// Tracer holds the steps for a request
type Tracer struct {
	Steps []Step
	mu    sync.RWMutex
}

// WithContext returns a new context with an initialized Tracer
func WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, traceKey, &Tracer{
		Steps: make([]Step, 0),
	})
}

// fromContext retrieves the Tracer from context
func fromContext(ctx context.Context) *Tracer {
	if v, ok := ctx.Value(traceKey).(*Tracer); ok {
		return v
	}
	return nil
}

// AddStep logs a step to the trace
func AddStep(ctx context.Context, step string, format string, args ...interface{}) {
	tracer := fromContext(ctx)
	if tracer == nil {
		return
	}

	details := fmt.Sprintf(format, args...)
	tracer.mu.Lock()
	defer tracer.mu.Unlock()

	tracer.Steps = append(tracer.Steps, Step{
		Time:    time.Now().Format("15:04:05.000"),
		Step:    step,
		Details: details,
	})
}

// GetTraceString returns the JSON representation of the trace
func GetTraceString(ctx context.Context) string {
	tracer := fromContext(ctx)
	if tracer == nil {
		return ""
	}

	tracer.mu.RLock()
	defer tracer.mu.RUnlock()

	if len(tracer.Steps) == 0 {
		return ""
	}

	bytes, _ := json.Marshal(tracer.Steps)
	return string(bytes)
}
