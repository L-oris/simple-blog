package inject

import (
	"github.com/sarulabs/di"
)

// Container stores all dependencies, allowing to easily inject them
var Container di.Container

func init() {
	Container = createBuilder().Build()
}

func createBuilder() *di.Builder {
	builder, _ := di.NewBuilder()

	builder.Add(core()...)
	builder.Add(repositories()...)
	builder.Add(controllers()...)
	builder.Add(routers()...)

	return builder
}
