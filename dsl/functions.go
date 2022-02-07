package dsl

import (
	"math/rand"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

const RandomID = "random"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var Random = function.New(&function.Spec{
	VarParam: nil,
	Params: []function.Parameter{
		{Type: cty.Number},
	},
	Type: func(args []cty.Value) (cty.Type, error) {
		return cty.String, nil
	},
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		n, _ := args[0].AsBigFloat().Int64()
		b := make([]rune, n)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		return cty.StringVal(string(b)), nil

	},
})



