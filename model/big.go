package model

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"strings"
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

func (z *BigInt) SetString(x string, base int) *BigInt {
	if res, err := (*big.Int)(z).SetString(x, base); err == true {
		return (*BigInt)(res)
	}
	return nil
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

func (x *BigInt) Cmp(y *BigInt) int {
	return (*big.Int)(x).Cmp((*big.Int)(y))
}

type BigFloat big.Float

func (b *BigFloat) Value() (driver.Value, error) {
	if b != nil {
		return (*big.Float)(b).String(), nil
	}
	return nil, nil
}

func (b *BigFloat) Scan(value interface{}) error {
	if value == nil {
		b = nil
	}

	switch t := value.(type) {
	case []uint8:
		ta := value.([]uint8)
		tt := string(ta)
		_, ok := (*big.Float)(b).SetString(tt)
		if !ok {
			return fmt.Errorf("failed to load value to []uint8: %v", value)
		}
	default:
		return fmt.Errorf("could not scan type %T into BigInt", t)
	}
	return nil
}

func (z *BigFloat) Set(x *BigFloat) *BigFloat {
	return (*BigFloat)((*big.Float)(z).Set((*big.Float)(x)))
}

func (z *BigFloat) SetString(x string) *BigFloat {
	if res, f := new(big.Float).SetPrec(uint(256)).SetString(x); f == true {
		return (*BigFloat)(res)
	}
	return nil
}

func (b *BigFloat) String() string {
	if buf, err := b.MarshalText(); err == nil {
		return string(buf)
	} else {
		return ""
	}
}

func (b *BigFloat) Text(format byte, prec int) string {
	return (*big.Float)(b).Text(format, prec)
}

func (x *BigFloat) MarshalText() (text []byte, err error) {
	return (*big.Float)(x).MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (z *BigFloat) UnmarshalText(text []byte) error {
	return (*big.Float)(z).UnmarshalText(text)
}

// The JSON marshalers are only here for API backward compatibility
// (programs that explicitly look for these two methods). JSON works
// fine with the TextMarshaler only.

// MarshalJSON implements the json.Marshaler interface.
func (x *BigFloat) MarshalJSON() ([]byte, error) {
	return x.MarshalText()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (z *BigFloat) UnmarshalJSON(text []byte) error {
	return z.UnmarshalText(text)
}

func (x *BigFloat) Cmp(y *BigFloat) int {
	return (*big.Float)(x).Cmp((*big.Float)(y))
}

func (z *BigFloat) ConvertToBigInt() *BigInt {
	s := strings.Split(z.String(), ".")
	ss := s[0] + s[1]
	return new(BigInt).SetString(ss, 10)
}
