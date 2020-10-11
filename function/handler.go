// function template

// Package function name is fixed.
package function

// Input struct is for ...<need to doc>
type Input struct {
	Name      string `json:"name"`
	NameState string `json:"name_state"`
}

// Handler func is ...<need to doc>
func Handler(i Input) interface{} {
	return i.Name
}
