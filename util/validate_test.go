package util

import (
	"testing"

	"github.com/asaskevich/govalidator"
)

type (
	customValidate struct {
		Age int `json:"age,omitempty" valid:"xMyValidate(0|10)"`
	}
	validateStruct struct {
		Age  int `json:"age,omitempty" valid:"xIntRange(0|100)"`
		Type int `json:"type,omitempty" valid:"xIntIn(1|5|10)"`
	}
)

func TestValidate(t *testing.T) {
	t.Run("default custom valid", func(t *testing.T) {
		buf := []byte(`{
			"age": 10,
			"type": 1
		}`)
		s := &validateStruct{}
		err := Validate(s, buf)
		if err != nil {
			t.Fatalf("default custom valid fail, %v", err)
		}
	})

	t.Run("validate fail", func(t *testing.T) {
		p := &params{}
		buf := []byte(`{"account":"abd"}`)
		err := Validate(p, buf)
		if err == nil {
			t.Fatalf("validate should be fail")
		}
	})
	t.Run("validate fail with not json buffer", func(t *testing.T) {
		p := &params{}
		buf := []byte(`{"account":"vicanso}`)
		err := Validate(p, buf)
		he := err.(*HTTPError)
		if he.Category != ErrCategoryJSON {
			t.Fatalf("validate should be json fail")
		}
	})

	t.Run("validate success", func(t *testing.T) {
		p := &params{}
		account := "vicanso"
		buf := []byte(`{"account":"vicanso"}`)
		err := Validate(p, buf)
		if err != nil || p.Account != account {
			t.Fatalf("validate fail, %v", err)
		}
		tmp := &params{}
		err = Validate(tmp, p)
		if err != nil || tmp.Account != account {
			t.Fatalf("validate fail, %v", err)
		}
	})

	t.Run("custom validate", func(t *testing.T) {

		AddRegexValidate("xMyValidate", "^xMyValidate\\((\\d+)\\|(\\d+)\\)$", func(value string, params ...string) bool {
			return govalidator.InRangeInt(value, params[0], params[1])
		})
		s := &customValidate{}
		err := Validate(s, []byte(`{
			"age": 10
		}`))
		if err != nil {
			t.Fatalf("add regexp validate fail, %v", err)
		}
		err = Validate(s, []byte(`{
			"age": 11
		}`))
		if err == nil {
			t.Fatalf("the age over the limit should return error")
		}
	})
}
