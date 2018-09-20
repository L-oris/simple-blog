package inject

import "github.com/sarulabs/di"

// Container stores all dependencies, allowing to easily inject them
var Container di.Container

func init() {
	Container = createBuilder().Build()
}

func createBuilder() *di.Builder {
	obj := di.Def{
		Name: "my-object",
		Build: func(ctn di.Container) (interface{}, error) {
			return &struct{ Name string }{Name: "Loris"}, nil
		},
	}

	builder, _ := di.NewBuilder()
	builder.Add(obj)

	return builder
}
