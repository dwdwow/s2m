package s2m

import (
	"fmt"
	"reflect"
	"testing"
)

type testStruct0 struct {
	A int
	B string
}

type testStruct struct {
	A string         `s2m:"a,omitempty"`
	B bool           `s2m:"b"`
	C int64          `s2m:"c,omitempty"`
	D int32          `s2m:"d"`
	E int16          `s2m:"e"`
	F int8           `s2m:"f"`
	G int            `s2m:"g"`
	H uint64         `s2m:"h"`
	I uint32         `s2m:"i"`
	J uint16         `s2m:"j"`
	K uint8          `s2m:"k"`
	L uint           `s2m:"l"`
	M uintptr        `s2m:"m"`
	N float64        `s2m:"n"`
	O float32        `s2m:"o"`
	P []int          `s2m:"p"`
	Q map[string]int `s2m:"q"`
	R testStruct0    `s2m:"r"`
	S [3]string      `s2m:"s"`
	T chan int       `s2m:"t"`
	U func()         `s2m:"u"`

	PA *string         `s2m:"pa"`
	PB *bool           `s2m:"pb"`
	PC *int64          `s2m:"pc"`
	PD *int32          `s2m:"pd"`
	PE *int16          `s2m:"pe"`
	PF *int8           `s2m:"pf"`
	PG *int            `s2m:"pg"`
	PH *uint64         `s2m:"ph"`
	PI *uint32         `s2m:"pi"`
	PJ *uint16         `s2m:"pj"`
	PK *uint8          `s2m:"pk"`
	PL *uint           `s2m:"pl,omitempty"`
	PM *uintptr        `s2m:"pm"`
	PN *float64        `s2m:"pn"`
	PO *float32        `s2m:"po"`
	PP *[]int          `s2m:"pp"`
	PQ *map[string]int `s2m:"pq"`
	PR *testStruct0    `s2m:"pr"`
	PS *[3]string      `s2m:"ps"`
	PT *chan int       `s2m:"pt"`
	PU *func()         `s2m:"pu"`

	Invalid reflect.Value `s2m:"invalid"`

	a string `s2m:"a"`
	b int    `s2m:"b"`
}

func TestToWithErr(t *testing.T) {
	a := "aaaaa"
	ts := testStruct{
		A:  a,
		PA: &a,
		R:  testStruct0{111111, "rrrrrr"},
	}
	m, err := ToWithErr(ts)
	if err != nil {
		panic(err)
	}
	for k, v := range m {
		switch k {
		case "c", "pl":
			panic("c is zero value and tag is omitempty, but is obtained in map")
		}
		fmt.Println(k, v)
	}
}

func TestToStrMapWithErr(t *testing.T) {
	a := "aaaaa"
	ts := testStruct{
		A:  a,
		PA: &a,
		R:  testStruct0{111111, "rrrrrr"},
	}
	m, err := ToStrMapWithErr(ts)
	if err != nil {
		panic(err)
	}
	for k, v := range m {
		switch k {
		case "c", "pl":
			panic("c is zero value and tag is omitempty, but is obtained in map")
		}
		fmt.Println(k, v)
	}
}
