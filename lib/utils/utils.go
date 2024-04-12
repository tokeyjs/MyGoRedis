package utils

import "strconv"

// ToCmdLine convert strings to [][]byte
func ToCmdLine(cmd ...string) [][]byte {
	args := make([][]byte, len(cmd))
	for i, s := range cmd {
		args[i] = []byte(s)
	}
	return args
}

// ToCmdLine2 convert commandName and []byte-type argument to CmdLine
func ToCmdLine2(commandName string, args ...[]byte) [][]byte {
	result := make([][]byte, len(args)+1)
	result[0] = []byte(commandName)
	for i, s := range args {
		result[i+1] = s
	}
	return result
}

// BytesEquals check whether the given bytes is equal
func BytesEquals(a []byte, b []byte) bool {
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	size := len(a)
	for i := 0; i < size; i++ {
		av := a[i]
		bv := b[i]
		if av != bv {
			return false
		}
	}
	return true
}

// float64 to string
func Float64ToString(val float64) string {
	return strconv.FormatFloat(val, 'g', -1, 64)
}

// float64 to string
func Float64ToByte(val float64) []byte {
	return []byte(strconv.FormatFloat(val, 'g', -1, 64))
}

func StringToFloat64(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}
