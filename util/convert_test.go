package util_test

import (
	"testing"

	"github.com/ntchjb/martian-querybody/util"
	"github.com/stretchr/testify/require"
)

func TestConvertAnyToString(t *testing.T) {
	var v interface{}

	v = float64(15)
	require.Equal(t, "15", util.ConvertAnyToString(v))

	v = "hello"
	require.Equal(t, "hello", util.ConvertAnyToString(v))

	v = true
	require.Equal(t, "true", util.ConvertAnyToString(v))

	v = float64(12.345)
	require.Equal(t, "12.345", util.ConvertAnyToString(v))
}
