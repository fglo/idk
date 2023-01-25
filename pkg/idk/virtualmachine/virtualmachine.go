package virtualmachine

import (
	"fmt"

	"github.com/fglo/idk/pkg/idk/chunk"
	"github.com/fglo/idk/pkg/idk/opcodes"
)

type function struct {
	name string
	args int
	code []byte
}

type callFrame struct {
	functionName string
	args         []int
	returnAddr   int
}

type Stacks struct {
	intStack    []int
	floatStack  []float64
	boolStack   []bool
	charStack   []rune
	stringStack []string
}

type VirtualMachine struct {
	chunk *chunk.Chunk
	ip    int

	intStack    stack[int]
	floatStack  stack[float64]
	boolStack   stack[bool]
	charStack   stack[rune]
	stringStack stack[string]

	memory        []int
	symbolTable   *symbolTable
	callStack     []callFrame
	functionTable map[string]*function
	loopCounters  []int
	loopLimits    []int
}

func NewVirtualMachine(chunk *chunk.Chunk) *VirtualMachine {
	i := &VirtualMachine{
		chunk: chunk,

		intStack:    newStack[int](),
		floatStack:  newStack[float64](),
		boolStack:   newStack[bool](),
		charStack:   newStack[rune](),
		stringStack: newStack[string](),

		memory:       make([]int, 0),
		symbolTable:  newSymbolTable(),
		callStack:    make([]callFrame, 0),
		loopCounters: make([]int, 0),
		loopLimits:   make([]int, 0),
	}

	return i
}

func (vm *VirtualMachine) Run() {
	bytecode := vm.chunk.Bytecode
	codeLength := len(bytecode)

	constantPool := vm.chunk.ConstantPool

	for vm.ip < codeLength {
		switch bytecode[vm.ip] {
		case opcodes.IPUSH:
			vm.ip++
			addr := int(bytecode[vm.ip])
			value := constantPool.RetrieveInt(addr)
			vm.intStack.push(value)
		case opcodes.IADD:
			a := vm.intStack.pop()
			b := vm.intStack.pop()
			vm.intStack.push(b + a)
		case opcodes.ISUB:
			a := vm.intStack.pop()
			b := vm.intStack.pop()
			vm.intStack.push(b - a)
		case opcodes.IMUL:
			a := vm.intStack.pop()
			b := vm.intStack.pop()
			vm.intStack.push(b * a)
		case opcodes.IDIV:
			a := vm.intStack.pop()
			b := vm.intStack.pop()
			vm.intStack.push(b / a)
		case opcodes.INEG:
			val := vm.intStack.pop()
			vm.intStack.push(-val)
		case opcodes.IPRINT:
			value := vm.intStack.pop()
			fmt.Println(value)
		case opcodes.IVAR_BIND:
			vm.ip++
			varnameAddr := int(bytecode[vm.ip])
			varname := constantPool.RetrieveString(varnameAddr)
			value := vm.intStack.pop()
			vm.symbolTable.bindInt(varname, value)
		case opcodes.IVAR_LOOKUP:
			vm.ip++
			varnameAddr := int(bytecode[vm.ip])
			varname := constantPool.RetrieveString(varnameAddr)
			value := vm.symbolTable.lookupInt(varname)
			vm.intStack.push(value)
		case opcodes.IFUNC_CREATE:
			vm.ip++
			nameLength := int(bytecode[vm.ip])
			vm.ip++
			name := string(bytecode[vm.ip : vm.ip+nameLength])
			vm.ip += nameLength
			numArgs := int(bytecode[vm.ip])
			vm.ip++
			funcCode := bytecode[vm.ip:]
			vm.functionTable[name] = &function{name: name, args: numArgs, code: funcCode}
		case opcodes.IFUNC_CALL:
			vm.ip++
			nameLength := int(bytecode[vm.ip])
			vm.ip++
			name := string(bytecode[vm.ip : vm.ip+nameLength])
			vm.ip += nameLength
			numArgs := int(bytecode[vm.ip])
			vm.ip++
			args := make([]int, numArgs)
			for j := 0; j < numArgs; j++ {
				args[j] = vm.intStack[len(vm.intStack)-1-j]
			}
			vm.intStack = vm.intStack[:len(vm.intStack)-numArgs]
			vm.callStack = append(vm.callStack, callFrame{
				functionName: name,
				args:         args,
				returnAddr:   vm.ip,
			})
			vm.ip = 0
			bytecode = vm.functionTable[name].code
		case opcodes.IFUNC_RETURN:
			returnValue := vm.intStack.pop()
			vm.ip = vm.callStack[len(vm.callStack)-1].returnAddr
			vm.callStack = vm.callStack[:len(vm.callStack)-1]
			vm.intStack.push(returnValue)
			if len(vm.callStack) > 0 {
				bytecode = vm.functionTable[vm.callStack[len(vm.callStack)-1].functionName].code
			}
		case opcodes.IF:
			a := vm.intStack.pop()
			b := vm.intStack.pop()
			if a < b {
				vm.ip++
				jump := int(bytecode[vm.ip])
				vm.ip += jump
			} else {
				vm.ip++
			}
		case opcodes.ELSE:
			vm.ip++
			jump := int(bytecode[vm.ip])
			vm.ip += jump
		case opcodes.ENDIF:
			// do nothing
		case opcodes.FOR:
			vm.ip++
			loopCounter := int(bytecode[vm.ip])
			vm.ip++
			loopLimit := int(bytecode[vm.ip])
			vm.loopCounters = append(vm.loopCounters, loopCounter)
			vm.loopLimits = append(vm.loopLimits, loopLimit)
		case opcodes.NEXT:
			loopCounter := vm.loopCounters[len(vm.loopCounters)-1]
			loopLimit := vm.loopLimits[len(vm.loopLimits)-1]
			loopCounter++
			if loopCounter < loopLimit {
				vm.loopCounters[len(vm.loopCounters)-1] = loopCounter
				vm.ip++
				jump := int(bytecode[vm.ip])
				vm.ip -= jump
			} else {
				vm.loopCounters = vm.loopCounters[:len(vm.loopCounters)-1]
				vm.loopLimits = vm.loopLimits[:len(vm.loopLimits)-1]
			}
		case opcodes.BREAK:
			vm.loopCounters = vm.loopCounters[:len(vm.loopCounters)-1]
			vm.loopLimits = vm.loopLimits[:len(vm.loopLimits)-1]
		case opcodes.HALT:
			return
		}

		vm.ip++
	}
}
