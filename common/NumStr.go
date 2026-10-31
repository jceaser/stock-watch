//***********************************************************************************************100
/*
The struct NumStr should contain either a Number or a string with supporting functions to test which
value is contained.
*/

package common

import (
	"fmt"
)

type NumStr struct {
	NumberValue *float64
	StringValue *string
}

func (n NumStr) HasNumber() bool {
  	return n.NumberValue != nil
}

func (n NumStr) HasString() bool {
  	return n.StringValue != nil
}

func (n NumStr) GetNumber() float64 {
  	if n.NumberValue == nil {
	  	return 0.0
  	}
  	return *n.NumberValue
}

func (n NumStr) GetString() string {
  	if n.StringValue == nil {
		return ""
  	}
  	return *n.StringValue
}

func (n NumStr) String() string {
	if n.HasNumber() {
		return fmt.Sprintf("%f", *n.NumberValue)
	}
	return *n.StringValue
}

func NewNumStrNumber(val float64) NumStr {
  	return NumStr{NumberValue: &val}
}

func NewNumStrString(val string) NumStr {
  	return NumStr{StringValue: &val}
}
