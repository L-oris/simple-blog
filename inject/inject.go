package inject

import (
	"github.com/sarulabs/di"
)

// CreateContainer creates the container where all dependencies are bound
func CreateContainer() di.Container {
	builder, _ := di.NewBuilder()

	builder.Add(core()...)
	builder.Add(repositories()...)
	builder.Add(controllers()...)
	builder.Add(routers()...)

	return builder.Build()
}
