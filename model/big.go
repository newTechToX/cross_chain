package model

import (
	"database/sql/driver"
	"fmt"
	"math/big"
)

type BigInt big.Int

func (b *BigInt) Value() (driver.Value, error) {
	if b != nil {
		return (*big.Int)(b).String(), nil
	}
	return nil, nil
}

func (b *BigInt) Scan(value interface{}) error {
	if value == nil {
		b = nil
	}

	switch t := value.(type) {
	case []uint8:
		_, ok := (*big.Int)(b).SetString(string(value.([]uint8)), 10)
		if !ok {
			return fmt.Errorf("failed to load value to []uint8: %v", value)
		}
	default:
		return fmt.Errorf("could not scan type %T into BigInt", t)
	}
	return nil
}

func (z *BigInt) Set(x *BigInt) *BigInt {
	return (*BigInt)((*big.Int)(z).Set((*big.Int)(x)))
}

func (b *BigInt) String() string {
	return (*big.Int)(b).String()
}

func (b *BigInt) Text(base int) string {
	return (*big.Int)(b).Text(base)
}

func (x *BigInt) MarshalText() (text []byte, err error) {
	return (*big.Int)(x).MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (z *BigInt) UnmarshalText(text []byte) error {
	return (*big.Int)(z).UnmarshalText(text)
}

// The JSON marshalers are only here for API backward compatibility
// (programs that explicitly look for these two methods). JSON works
// fine with the TextMarshaler only.

// MarshalJSON implements the json.Marshaler interface.
func (x *BigInt) MarshalJSON() ([]byte, error) {
	return x.MarshalText()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (z *BigInt) UnmarshalJSON(text []byte) error {
	return z.UnmarshalText(text)
}
