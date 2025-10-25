package common

import (
	"encoding/json"
	"math"
)

type RoundedInt int

func (ri *RoundedInt) UnmarshalJSON(b []byte) error {
	var f float64
	if err := json.Unmarshal(b, &f); err != nil {
		return err
	}
	*ri = RoundedInt(math.Ceil(f))
	return nil
}
