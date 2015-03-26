// Package gostrich provides a simple function chaining mechanism.
// It allows chaining functions and the arguments they require.
// It does not ensure the correctness of a chain until the build is requested.
package gostrich

import (
	"reflect"
)

// Chain holds a current list of links to chain.
type Chain struct {
	links []interface{}
}

// VarArgs is a named type used to define the number of arguments
// that should be passed to the variadic parameter of the preceeding function.
type VarArgs int

// New creates a new empty chain where links are set to nil.
func New() *Chain {
	return &Chain{nil}
}

// Compose is a method which appends links passed to the current chain.
// ie.
// Chain([f, g]).Compose(h, k) -> Chain([f, g, h, k]); f(g(h(k)))
func (c *Chain) Compose(links ...interface{}) *Chain {
	return &Chain{append(c.links, links...)}
}

// Then is a method which prepends links passed to the current chain.
// It uses the reverse order.
// ie.
// Chain([f, g]).Then(h, k) -> Chain([k, h, f, g]); k(h(f(g)))
func (c *Chain) Then(links ...interface{}) *Chain {
	var reverse []interface{}
	if n := len(links); n > 0 {
		reverse = make([]interface{}, n)
		for i, link := range links {
			reverse[n - 1 - i] = link
		}
	}
	return &Chain{append(reverse, c.links...)}
}

// MergeCompose is a method which appends links from the chains passed
// to the current chain.
// ie.
// Chain([f, g]).MergeCompose(Chain([h, k]), Chain([l, m]))
// -> Chain([f, g, h, k, l, m]); f(g(h(k(l(m)))))
func (c *Chain) MergeCompose(chains ...*Chain) *Chain {
	result := &Chain{c.links}
	for _, chain := range chains {
		result = result.Compose(chain.links...)
	}
	return result
}

// MergeThen is a method which prepends links from the chains passed
// to the current chain in the revese order.
// ie.
// Chain([f, g]).MergeThen(Chain([h, k]), Chain([l, m]))
// -> Chain([k, h, m, l, f, g]); k(h(m(l(f(g)))))
func (c *Chain) MergeThen(chains ...*Chain) *Chain {
	result := &Chain{c.links}
	n := len(chains)
	reverse := make([]*Chain, n)
	for i, chain := range chains {
		reverse[n - 1 - i] = chain
	}
	for _, chain := range chains {
		result = result.Then(chain.links...)
	}
	return result
}

// Build is a method which executes the constructed chain and
// returns an array of interface{} containing the return values.
//
// eg.
// plus: f(int, int) -> int
// timesTwo: f(int) -> int
// n: int = 2
// m: int = 3
// chain: Chain([f, timesTwo, n, m])
// chain.Build() returns [f(timesTwo(n), m)] = [7]
//
// It traverses the chain in the reverse order, puts values on the stack,
// executes functions with the correct number of arguments taken from the stack,
// puts return values on the stack in the reverse order.
//
// For functions with a variadic parameter it:
// -if top(stack) is VarArgs: takes VarArgs number of arguments for the variadic parameter
// -if top(stack) is not VarArgs: takes as many arguments as possible, depending on the type
//
// It panics when there are not enough values on the stack for the current function call.
//
// Be careful about the references as links inside the chain.
func (c *Chain) Build() []interface{} {
	var stack []reflect.Value
	for i := len(c.links) - 1; i >= 0; i-- {
		value := reflect.ValueOf(c.links[i])
		switch value.Kind() {
		case reflect.Func:
			n, m := len(stack), value.Type().NumIn()
			if (value.Type().IsVariadic() && n != 0) {
				if variadicNo, ok := stack[n-1].Interface().(VarArgs); ok {
					m, n = m + int(variadicNo) - 1, n - 1
					stack = stack[:n]
				} else {
					variadicType := value.Type().In(m - 1).Elem()
					for n - m >= 0 && stack[n-m].Type() == variadicType {
						m++
					}
					m--
				}
			}
			if (n - m < 0) {
				panic("Gost: Build with incomplete chain")
			}
			args := make([]reflect.Value, m)
			for j, arg := range stack[n-m:] {
				args[m - 1 - j] = arg
			}
			rets := value.Call(args)
			p := len(rets)
			reverseRets := make([]reflect.Value, p)
			for j, ret := range rets {
				reverseRets[p - 1 - j] = ret
			}
			stack = append(stack[:n-m], reverseRets...)
		default:
			stack = append(stack, value)
		}
	}
	n := len(stack)
	result := make([]interface{}, n)
	for i, value := range(stack) {
		result[n - 1 - i] = value.Interface()
	}
	return result
}
