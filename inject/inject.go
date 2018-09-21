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

	builder.Add(getCore()...)
	builder.Add(getRepositories()...)
	builder.Add(getControllers()...)

	return builder
}
