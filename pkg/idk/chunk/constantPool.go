package chunk

type ConstantPool struct {
	intPool    []int
	floatPool  []float64
	boolPool   []bool
	charPool   []rune
	stringPool []string
}

func NewConstantPool() *ConstantPool {
	return &ConstantPool{
		intPool:    make([]int, 0),
		floatPool:  make([]float64, 0),
		boolPool:   make([]bool, 0),
		charPool:   make([]rune, 0),
		stringPool: make([]string, 0),
	}
}

func (cp *ConstantPool) InsertInt(val int) int {
	cp.intPool = append(cp.intPool, val)
	return len(cp.intPool) - 1
}

func (cp *ConstantPool) RetrieveInt(address int) int {
	return cp.intPool[address]
}

func (cp *ConstantPool) InsertFloat(val float64) int {
	cp.floatPool = append(cp.floatPool, val)
	return len(cp.floatPool) - 1
}

func (cp *ConstantPool) RetrieveFloat(address int) float64 {
	return cp.floatPool[address]
}

func (cp *ConstantPool) InsertBool(val bool) int {
	cp.boolPool = append(cp.boolPool, val)
	return len(cp.boolPool) - 1
}

func (cp *ConstantPool) RetrieveBool(address int) bool {
	return cp.boolPool[address]
}

func (cp *ConstantPool) InsertChar(val rune) int {
	cp.charPool = append(cp.charPool, val)
	return len(cp.charPool) - 1
}

func (cp *ConstantPool) RetrieveChar(address int) rune {
	return cp.charPool[address]
}

func (cp *ConstantPool) InsertString(val string) int {
	cp.stringPool = append(cp.stringPool, val)
	return len(cp.stringPool) - 1
}

func (cp *ConstantPool) RetrieveString(address int) string {
	return cp.stringPool[address]
}
