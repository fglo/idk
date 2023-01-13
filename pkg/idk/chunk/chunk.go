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

func (c *Chunk) AddBoolConstant(val bool) int {
	return c.ConstantPool.InsertBool(val)
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

func (c *Chunk) GetBoolConstant(addr int) bool {
	return c.ConstantPool.RetrieveBool(addr)
}

func (c *Chunk) GetCharConstant(addr int) rune {
	return c.ConstantPool.RetrieveChar(addr)
}

func (c *Chunk) GetStringConstant(addr int) string {
	return c.ConstantPool.RetrieveString(addr)
}

func (c *Chunk) Disassemble() string {
	var out bytes.Buffer

	out.WriteString("──────┬────────────────────────\n")
	out.WriteString(" IP   │ OPCODE           PARAM \n")
	out.WriteString("──────┼────────────────────────\n")

	for ip := 0; ip < len(c.Bytecode); ip++ {
		bcode := c.Bytecode[ip]
		code := opcodes.ToString(bcode)
		if bcode == opcodes.IPUSH || bcode == opcodes.IVAR_BIND || bcode == opcodes.IVAR_LOOKUP {
			if ip < len(c.Bytecode)-1 {
				ip++
				param := c.Bytecode[ip]
				out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-12s %v\n", ip, bcode, code, param))
			}
			out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-12s\n", ip, bcode, code))
		} else {
			out.WriteString(fmt.Sprintf(" %-4d │ %04d  %-12s\n", ip, bcode, code))
		}
	}

	out.WriteString("──────┴────────────────────────\n")

	return out.String()
}
