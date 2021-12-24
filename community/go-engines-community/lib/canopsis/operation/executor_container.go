package operation

import "fmt"

// ExecutorContainer interface is used to implement an storage of executors.
type ExecutorContainer interface {
	// Get returns executor.
	Get(operationType string) (Executor, bool)
	// Set stores executor.
	Set(operationType string, executor Executor)
	// Has returns if container contains executor.
	Has(operationType string) bool
}

// NewExecutorContainer creates new container.
func NewExecutorContainer() ExecutorContainer {
	return &mapExecutorContainer{
		executors: make(map[string]Executor),
	}
}

type mapExecutorContainer struct {
	executors map[string]Executor
}

func (c *mapExecutorContainer) Get(operationType string) (Executor, bool) {
	e, ok := c.executors[operationType]

	return e, ok
}

func (c *mapExecutorContainer) Set(operationType string, executor Executor) {
	if c.Has(operationType) {
		panic(fmt.Errorf("operation executor %q already exists", operationType))
	}

	c.executors[operationType] = executor
}

func (c *mapExecutorContainer) Has(operationType string) bool {
	_, ok := c.executors[operationType]

	return ok
}
