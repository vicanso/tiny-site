package config

import (
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestGetENV(t *testing.T) {
	env := GetENV()
	if env != "test" {
		t.Fatalf("get env fail")
	}
}

func TestGetValue(t *testing.T) {
	noneKey := "none"
	i := "i"
	iValue := 1

	s := "s"
	sValue := "a"

	d := "d"
	dValue := time.Second

	ss := "ss"
	ssValue := []string{
		"a",
	}

	t.Run("get int", func(t *testing.T) {
		viper.Set(i, iValue)
		if GetInt(i) != iValue {
			t.Fatalf("get int fail")
		}
		if GetIntDefault(i, 2) != 1 {
			t.Fatalf("get default int(exist) fail")
		}
		if GetIntDefault(noneKey, 2) != 2 {
			t.Fatalf("get default int(not exitst) fail")
		}
	})

	t.Run("get string", func(t *testing.T) {
		viper.Set(s, sValue)
		if GetString(s) != sValue {
			t.Fatalf("get string fail")
		}

		if GetStringDefault(s, "b") != sValue {
			t.Fatalf("get default string(exists) fail")
		}

		if GetStringDefault(noneKey, "b") != "b" {
			t.Fatalf("get default string(not exists) fail")
		}
	})

	t.Run("get duration", func(t *testing.T) {
		viper.Set(d, dValue)

		if GetDuration(d) != dValue {
			t.Fatalf("get duration fail")
		}

		if GetDurationDefault(d, time.Hour) != dValue {
			t.Fatalf("get default duration(exists) fail")
		}

		if GetDurationDefault(noneKey, time.Hour) != time.Hour {
			t.Fatalf("get default duration(not exists) fail")
		}
	})

	t.Run("get string slice", func(t *testing.T) {
		viper.Set(ss, ssValue)
		if strings.Join(GetStringSlice(ss), ",") != strings.Join(ssValue, "") {
			t.Fatalf("get string slice fail")
		}
	})

	t.Run("get track key", func(t *testing.T) {
		if GetTrackKey() != "jt" {
			t.Fatalf("get track key fail")
		}
		customTrack := "vicanso"
		viper.Set("track", customTrack)

		if GetTrackKey() != customTrack {
			t.Fatalf("get track key fail")
		}
	})

	t.Run("get session keys", func(t *testing.T) {
		keys := GetSessionKeys()
		if strings.Join(keys, ",") != "cuttlefish" {
			t.Fatalf("get session keys fail")
		}
	})

	t.Run("get session cookie", func(t *testing.T) {
		cookieName := GetSessionCookie()
		if cookieName != "sess" {
			t.Fatalf("get session cookie's name fail")
		}
	})

	t.Run("get cookie path", func(t *testing.T) {
		path := GetCookiePath()
		if path != "/" {
			t.Fatalf("get cookie's path fail")
		}
	})
}

func TestLoadFile(t *testing.T) {
	err := viperInit("../configs")
	if err != nil {
		t.Fatalf("init from file fail, %v", err)
	}
}
