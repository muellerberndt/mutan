package mutan

import (
	"fmt"
	"math/big"
	"os"
)

// Op codes
var OpCodes = map[string]byte{
	// 0x0 range - arithmetic ops
	"stop": 0x00,
	"add":  0x01,
	"mul":  0x02,
	"sub":  0x03,
	"div":  0x04,
	"sdiv": 0x05,
	"mod":  0x06,
	"smod": 0x07,
	"exp":  0x08,
	"neg":  0x09,
	"lt":   0x0a,
	"gt":   0x0b,
	"eq":   0x0c,
	"not":  0x0d,

	// 0x10 range - bit ops
	"and":  0x10,
	"or":   0x11,
	"xor":  0x12,
	"byte": 0x13,

	// 0x20 range - crypto
	"sha3": 0x20,

	// 0x30 range - closure state
	"address":      0x30,
	"balance":      0x31,
	"origin":       0x32,
	"caller":       0x33,
	"callvalue":    0x34,
	"calldataload": 0x35,
	"calldatasize": 0x36,
	"gasprice":     0x38,

	// 0x40 range - block operations
	"prevhash":   0x40,
	"coinbase":   0x41,
	"timestamp":  0x42,
	"number":     0x43,
	"difficulty": 0x44,
	"gaslimit":   0x45,

	// 0x50 range - 'storage' and execution
	"pop":     0x51,
	"dup":     0x52,
	"swap":    0x53,
	"mload":   0x54,
	"mstore":  0x55,
	"mstore8": 0x56,
	"sload":   0x57,
	"sstore":  0x58,
	"jump":    0x59,
	"jumpi":   0x5a,
	"pc":      0x5b,
	"msize":   0x5c,

	// 0x60 range - 'push'
	"push1":  0x60,
	"push2":  0x61,
	"push3":  0x62,
	"push4":  0x63,
	"push5":  0x64,
	"push6":  0x65,
	"push7":  0x66,
	"push8":  0x67,
	"push9":  0x68,
	"push10": 0x69,
	"push11": 0x6a,
	"push12": 0x6b,
	"push13": 0x6c,
	"push14": 0x6d,
	"push15": 0x6e,
	"push16": 0x6f,
	"push17": 0x70,
	"push18": 0x71,
	"push19": 0x72,
	"push20": 0x73,
	"push21": 0x74,
	"push22": 0x75,
	"push23": 0x76,
	"push24": 0x77,
	"push25": 0x78,
	"push26": 0x79,
	"push27": 0x7a,
	"push28": 0x7b,
	"push29": 0x7c,
	"push30": 0x7d,
	"push31": 0x7e,
	"push32": 0x7f,

	// 0xf0 range - closures
	"create": 0xf0,
	"call":   0xf1,
	"return": 0xf2,

	// 0x70 range - other
	"log":     0x70,
	"suicide": 0x7f,
}

// Big to bytes
//
// Returns the bytes of a big integer with the size specified by
// **base**
// Attempts to pad the byte array with zeros.
func bigToBytes(num *big.Int, base int) []byte {
	ret := make([]byte, base/8)

	return append(ret[:len(ret)-len(num.Bytes())], num.Bytes()...)
}

// Is op code
//
// Check whether the given string matches anything in
// the OpCode list
func IsOpCode(s string) bool {
	for key, _ := range OpCodes {
		if key == s {
			return true
		}
	}
	return false
}

// Compile instruction
//
// Attempts to compile and parse the given instruction in "s"
// and returns the byte sequence
func CompileInstr(s interface{}) ([]byte, error) {
	switch s.(type) {
	case string:
		str := s.(string)
		isOp := IsOpCode(str)
		if isOp {
			return []byte{OpCodes[str]}, nil
		}

		// Check for pre formatted byte array
		// Jumps are preformatted
		if []byte(str)[0] == 0 {
			return []byte(str), nil
		}

		num := new(big.Int)
		_, success := num.SetString(str, 0)
		// Assume regular bytes during compilation
		if !success {
			num.SetBytes([]byte(str))
		}

		return num.Bytes(), nil
	case int:
		//num := bigToBytes(big.NewInt(int64(s.(int))), 256)
		return big.NewInt(int64(s.(int))).Bytes(), nil
	case []byte:
		return new(big.Int).SetBytes(s.([]byte)).Bytes(), nil
	}

	return nil, nil
}

// Assemble
//
// Assembles the given instructions and returns EVM byte code
func Assemble(instructions ...interface{}) (script []byte) {
	//script = make([]string, len(instructions))

	for _, val := range instructions {
		instr, err := CompileInstr(val)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(instr) == 0 {
			instr = []byte{0}
		}

		//script[i] = string(instr)
		script = append(script, instr...)
	}

	return
}

// Pre parse script
//
// Take data apart and attempt to find the "init" section and
// "main" section. `main { } init { }`
func PreParse(data string) (mainInput, initInput string) {
	mainInput = getCodeSectionFor("main", data)
	if mainInput == "" {
		mainInput = data
	}
	initInput = getCodeSectionFor("init", data)

	return
}

// Very, very dumb parser. Heed no attention :-)
func getCodeSectionFor(blockMatcher, input string) string {
	curCount := -1
	length := len(blockMatcher)
	matchfst := rune(blockMatcher[0])
	var currStr string

	for i, run := range input {
		// Find init
		if curCount == -1 && run == matchfst && input[i:i+length] == blockMatcher {
			curCount = 0
		} else if curCount > -1 {
			if run == '{' {
				curCount++
				if curCount == 1 {
					continue
				}
			} else if run == '}' {
				curCount--
				if curCount == 0 {
					// we are done
					curCount = -1
					break
				}
			}

			if curCount > 0 {
				currStr += string(run)
			}
		}
	}

	return currStr
}
