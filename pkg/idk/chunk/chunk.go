package chunk

import (
	"bytes"
	"fmt"

	"github.com/fglo/idk/pkg/idk/opcodes"
)

type Chunk struct {
	Bytecode     []byte
	ConstantPool *ConstantPool
}

func NewChunk() *Chunk {
	return &Chunk{
		Bytecode:     make([]byte, 0),
		ConstantPool: NewConstantPool(),
	}
}

func (c *Chunk) Write(byte byte) int {
	c.Bytecode = append(c.Bytecode, byte)
	return len(c.Bytecode) - 1
}

func (c *Chunk) WriteBytes(bytes []byte) int {
	c.Bytecode = append(c.Bytecode, bytes...)
	return len(c.Bytecode) - 1
}

func (c *Chunk) AddIntConstant(val int) int {
	return c.ConstantPool.InsertInt(val)
}

func (c *Chunk) AddFloatConstant(val float64) int {
	return c.ConstantPool.InsertFloat(val)
}

func (c *Chunk) AddCharConstant(val rune) int {
	return c.ConstantPool.InsertChar(val)
}

func (c *Chunk) AddStringConstant(val string) int {
	return c.ConstantPool.InsertString(val)
}

func (c *Chunk) GetIntConstant(addr int) int {
	return c.ConstantPool.RetrieveInt(addr)
}

func (c *Chunk) GetFloatConstant(addr int) float64 {
	return c.ConstantPool.RetrieveFloat(addr)
}

func (c *Chunk) GetCharConstant(addr int) rune {
	return c.ConstantPool.RetrieveChar(addr)
}

func (c *Chunk) GetStringConstant(addr int) string {
	return c.ConstantPool.RetrieveString(addr)
}

func (c *Chunk) Disassemble() string {
	var out bytes.Buffer

	out.WriteString("──────┬───────────────────────────────────\n")
	out.WriteString(" IP   │ OPCODE             ADDR     VALUE \n")
	out.WriteString("──────┼───────────────────────────────────\n")

	for ip := 0; ip < len(c.Bytecode); ip++ {
		currentIP := ip
		bcode := c.Bytecode[ip]
		code := opcodes.ToString(bcode)
		switch bcode {
		case opcodes.IPUSH:
			if ip < len(c.Bytecode)-1 {
				ip++
				param := c.Bytecode[ip]
				value := c.ConstantPool.RetrieveInt(int(param))
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s %-9v %d\n", currentIP, bcode, code, param, value))
			} else {
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s\n", currentIP, bcode, code))
			}
		case opcodes.FPUSH:
			if ip < len(c.Bytecode)-1 {
				ip++
				param := c.Bytecode[ip]
				value := c.ConstantPool.RetrieveFloat(int(param))
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s %-9v %f\n", currentIP, bcode, code, param, value))
			} else {
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s\n", currentIP, bcode, code))
			}
		case opcodes.BPUSH:
			if ip < len(c.Bytecode)-1 {
				ip++
				valueByte := c.Bytecode[ip]
				value := false
				if valueByte == 1 {
					value = true
				}
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s %-9s %t\n", currentIP, bcode, code, "", value))
			} else {
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s\n", currentIP, bcode, code))
			}
		case opcodes.CPUSH:
			if ip < len(c.Bytecode)-1 {
				ip++
				param := c.Bytecode[ip]
				value := c.ConstantPool.RetrieveChar(int(param))
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s %-9v %c\n", currentIP, bcode, code, param, value))
			} else {
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s\n", currentIP, bcode, code))
			}
		case opcodes.SPUSH:
			if ip < len(c.Bytecode)-1 {
				ip++
				param := c.Bytecode[ip]
				value := c.ConstantPool.RetrieveString(int(param))
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s %-9v %s\n", currentIP, bcode, code, param, value))
			} else {
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s\n", currentIP, bcode, code))
			}
		case opcodes.IVAR_BIND, opcodes.IVAR_LOOKUP,
			opcodes.FVAR_BIND, opcodes.FVAR_LOOKUP,
			opcodes.BVAR_BIND, opcodes.BVAR_LOOKUP,
			opcodes.CVAR_BIND, opcodes.CVAR_LOOKUP,
			opcodes.SVAR_BIND, opcodes.SVAR_LOOKUP:
			if ip < len(c.Bytecode)-1 {
				ip++
				param := c.Bytecode[ip]
				value := c.ConstantPool.RetrieveString(int(param))
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s %-9v %s\n", currentIP, bcode, code, param, value))
			} else {
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s\n", currentIP, bcode, code))
			}
		case opcodes.IFUNC_CREATE,
			opcodes.FFUNC_CREATE,
			opcodes.BFUNC_CREATE,
			opcodes.CFUNC_CREATE,
			opcodes.SFUNC_CREATE:
			if ip < len(c.Bytecode)-1 {
				ip++
				param := c.Bytecode[ip]
				value := c.ConstantPool.RetrieveString(int(param))
				ip++
				numArgs := c.Bytecode[ip]
				ip++
				funcLen := c.Bytecode[ip]
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s %-9v %-10s %d %d\n", currentIP, bcode, code, param, value, numArgs, funcLen)) // TODO: display number of arguments
			} else {
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s\n", currentIP, bcode, code))
			}
		case opcodes.IFUNC_CALL,
			opcodes.FFUNC_CALL,
			opcodes.BFUNC_CALL,
			opcodes.CFUNC_CALL,
			opcodes.SFUNC_CALL:
			if ip < len(c.Bytecode)-1 {
				ip++
				param := c.Bytecode[ip]
				value := c.ConstantPool.RetrieveString(int(param))
				ip++
				numArgs := c.Bytecode[ip]
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s %-9v %-10s %d\n", currentIP, bcode, code, param, value, numArgs)) // TODO: display number of arguments
			} else {
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s\n", currentIP, bcode, code))
			}
		case opcodes.IF:
			ip++
			out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s %-9d\n", currentIP, bcode, code, int(c.Bytecode[ip])))
		default:
			out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-15s\n", currentIP, bcode, code))
		}
	}

	out.WriteString("──────┴───────────────────────────────────\n")

	return out.String()
}
