package handlers

import (
	"testing"
)

func TestFilterStrings(t *testing.T) {
	inputs := []string{
		"This is a test, and replaced should be this: [[\"Jl. SMA Aek Kota Batu\",\"id\"],[\"Sumatera Utara\",\"de\"]], yes that is what should be replaced.",
		"[[\"а/д Вятка\",\"ru\"]]",
	}
	outputs := []string{
		"This is a test, and replaced should be this: [[\"\",\"\"],[\"\",\"\"]], yes that is what should be replaced.",
		"[[\"\",\"\"]]",
	}

	for i := range inputs {
		out := string(filterStrings([]byte(inputs[i])))
		if out != outputs[i] {
			t.Fatal("Expected\n", outputs[i], "\nbut got\n", out)
		}
	}
}
