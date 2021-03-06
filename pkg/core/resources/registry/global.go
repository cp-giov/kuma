package registry

import "github.com/kumahq/kuma/pkg/core/resources/model"

var global = NewTypeRegistry()

func Global() TypeRegistry {
	return global
}

func RegisterType(res model.ResourceTypeDescriptor) {
	if err := global.RegisterType(res); err != nil {
		panic(err)
	}
}
