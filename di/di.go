package di

import "github.com/sarulabs/di"

func Create() *di.Builder {
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
