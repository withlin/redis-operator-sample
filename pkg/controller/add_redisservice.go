package controller

import (
	"github.com/ym/redis-operator/pkg/controller/redisservice"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, redisservice.Add)
}
