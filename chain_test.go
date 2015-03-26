package gostrich

import (
	"testing"
)

func TestBuildOnNewChainCreatesEmptyArray(t *testing.T) {
	result := New().Build()
	if (result == nil || len(result) != 0) {
		t.Fail()
	}
}

func TestNewCreatesNewChain(t *testing.T) {
	chain := New()
	otherChain := New()
	if (&chain == &otherChain || &chain.links == &otherChain.links || chain.links != nil || otherChain.links != nil) {
		t.Fail()
	}
}

func TestComposeCreatesNewChain(t *testing.T) {
	chain := New()
	composed := chain.Compose(1)
	if (&chain == &composed || &chain.links == &composed.links || chain.links != nil || composed.links == nil) {
		t.Fail()
	}
}

func TestThenCreatesNewChain(t *testing.T) {
	chain := New()
	then := chain.Then(1)
	if (&chain == &then || &chain.links == &then.links || chain.links != nil || then.links == nil) {
		t.Fail()
	}
}

func TestMergeComposeCreatesNewChain(t *testing.T) {
	chain := New()
	merged := chain.MergeCompose(New())
	if (&chain == &merged || &chain.links == &merged.links || chain.links != nil || merged.links != nil) {
		t.Fail()
	}
}

func TestMergeThenCreatesNewChain(t *testing.T) {
	chain := New()
	merged := chain.MergeThen(New())
	if (&chain == &merged || &chain.links == &merged.links || chain.links != nil || merged.links != nil) {
		t.Fail()
	}
}

func TestComposeCreatesChainWithLinksInTheGivenOrder(t *testing.T) {
	chain := New().Compose(1, 2)
	if (chain.links[0] != 1 || chain.links[1] != 2) {
		t.Fail()
	}
}


func TestThenCreatesChainWithLinksInTheReverseOrder(t *testing.T) {
	chain := New().Then(1, 2)
	if (chain.links[0] != 2 || chain.links[1] != 1) {
		t.Fail()
	}
}

func TestMergeComposeCreatesChainWithLinksFromChainsInTheGivenOrder(t *testing.T) {
	chain1 := New().Compose(1, 2)
	chain2 := New().Compose(3, 4)
	chain := New().Compose(0).MergeCompose(chain1, chain2)
	if (chain.links[0] != 0 || chain.links[1] != 1 || chain.links[2] != 2 || chain.links[3] != 3 || chain.links[4] != 4) {
		t.Fail()
	}
}

func TestMergeThenCreatesChainWithLinksFromChainsInTheReverseOrder(t *testing.T) {
	chain1 := New().Compose(1, 2)
	chain2 := New().Compose(3, 4)
	chain := New().Compose(0).MergeThen(chain1, chain2)
	if (chain.links[0] != 4 || chain.links[1] != 3 || chain.links[2] != 2 || chain.links[3] != 1 || chain.links[4] != 0) {
		t.Fail()
	}
}

func TestBuildWithVariadicFunctionAndNoVarArgsTakesAsManyArgsAsPossible(t *testing.T) {
	varFunc := func(a int, ns ...int) int {
		result := 0
		for _, n := range ns {
			result += a * n
		}
		return result
	}
	simpleFunc := func() int {
		return 3
	}
	result := New().Compose(varFunc).Compose(10, 1, 2, simpleFunc, "dog", 9).Build()
	if (len(result) != 3 || result[0] != 60 || result[1] != "dog" || result[2] != 9) {
		t.Fail()
	}
}

func TestBuildWithVariadicFunctionAndVarArgsSetTakesLimitedNumberOfArgs(t *testing.T) {
	varFunc := func(a int, ns ...int) int {
		result := 0
		for _, n := range ns {
			result += a * n
		}
		return result
	}
	simpleFunc := func() int {
		return 3
	}
	result := New().Compose(varFunc, VarArgs(2)).Compose(10, 1, 2, simpleFunc, "dog", 9).Build()
	if (len(result) != 4 || result[0] != 30 || result[1] != 3 || result[2] != "dog" || result[3] != 9) {
		t.Fail()
	}
}

func TestBuildShouldExecuteInTheLeftToRightManner(t *testing.T) {
	f1 := func(a int, b int, c int, d int) (int, int, int, int) {
		return a, b, c, d
	} 
	f2 := func(a int, b int) int {
		return a + b
	}
	f3 := func(a int, b int, c int) (int, int) {
		return a * b, b * c
	}
	result := New().Compose(f1, f2, 10, 5, f3, 2, 4, 8, 100).Build()
	if (len(result) != 4 || result[0] != 15 || result[1] != 8 || result[2] != 32 || result[3] != 100) {
		t.Fail()
	}
}

func TestComposeAcceptsNoArguments(t *testing.T) {
	chain := New()
	composed := chain.Compose()
	if (&chain == &composed || &chain.links == &composed.links || chain.links != nil || composed.links != nil) {
		t.Fail()
	}
}

func TestThenAcceptsNoArguments(t *testing.T) {
	chain := New()
	then := chain.Compose()
	if (&chain == &then || &chain.links == &then.links || chain.links != nil || then.links != nil) {
		t.Fail()
	}
}

func TestMergeComposeAcceptsNoArgumentsAndReturnChainCopy(t *testing.T) {
	chain := New()
	composed := chain.MergeCompose()
	if (&chain == &composed || &chain.links == &composed.links || chain.links != nil || composed.links != nil) {
		t.Fail()
	}
}

func TestMergeThenAcceptsNoArgumentsAndReturnsChainCopy(t *testing.T) {
	chain := New()
	then := chain.MergeThen()
	if (&chain == &then || &chain.links == &then.links || chain.links != nil || then.links != nil) {
		t.Fail()
	}
}

func TestComposeAcceptsNil(t *testing.T) {
	chain := New()
	composed := chain.Compose(nil)
	if (&chain == &composed || &chain.links == &composed.links || chain.links != nil || composed.links == nil ||
		composed.links[0] != nil) {
		t.Fail()
	}
}

func TestThenAcceptsNil(t *testing.T) {
	chain := New()
	then := chain.Compose(nil)
	if (&chain == &then || &chain.links == &then.links || chain.links != nil || then.links == nil ||
		then.links[0] != nil) {
		t.Fail()
	}
}

func TestIncompleteBuildShouldPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		} else if s, ok := r.(string); !ok || s != "Gost: Build with incomplete chain" {
			t.Fail()
		}
	}()
	f := func(a int){}
	New().Compose(f).Build()
}
