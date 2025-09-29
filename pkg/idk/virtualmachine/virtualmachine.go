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

		memory:        make([]int, 0),
		symbolTable:   newSymbolTable(),
		callStack:     make([]callFrame, 0),
		functionTable: make(map[string]*function),
		loopCounters:  make([]int, 0),
		loopLimits:    make([]int, 0),
	}

	return i
}

func (vm *VirtualMachine) PrintStacks() {
	fmt.Println("CALL stack:", vm.callStack)
	fmt.Println("INT stack:", vm.intStack)
	fmt.Println("FLOAT stack:", vm.floatStack)
	fmt.Println("BOOL stack:", vm.boolStack)
	fmt.Println("CHAR stack:", vm.charStack)
	fmt.Println("STRING stack:", vm.stringStack)
}

func (vm *VirtualMachine) Run() {
	defer func() {
		if r := recover(); r != nil {
			vm.PrintStacks()
			panic(fmt.Sprintf("runtime error at ip %d", vm.ip))
		}
	}()

	bytecode := vm.chunk.Bytecode
	codeLength := len(bytecode)

	constantPool := vm.chunk.ConstantPool

	// TODO: handle calling unexistant functino
	// TODO: handle calling function on variable which isn't a function
	// TODO: using function as a variable

	for vm.ip < codeLength {
		switch bytecode[vm.ip] {
		// INT
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
		case opcodes.IGT:
			b := vm.intStack.pop()
			a := vm.intStack.pop()
			vm.boolStack.push(a > b)
		case opcodes.ILT:
			b := vm.intStack.pop()
			a := vm.intStack.pop()
			vm.boolStack.push(a < b)
		case opcodes.IGE:
			b := vm.intStack.pop()
			a := vm.intStack.pop()
			vm.boolStack.push(a >= b)
		case opcodes.ILE:
			b := vm.intStack.pop()
			a := vm.intStack.pop()
			vm.boolStack.push(a <= b)
		case opcodes.IEQ:
			b := vm.intStack.pop()
			a := vm.intStack.pop()
			vm.boolStack.push(a == b)
		case opcodes.INEQ:
			b := vm.intStack.pop()
			a := vm.intStack.pop()
			vm.boolStack.push(a != b)
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
			value, err := vm.symbolTable.lookupInt(varname)
			if err != nil {
				panic(err) // TODO: maybe better error handling
			}
			vm.intStack.push(value)
		case opcodes.IFUNC_CREATE:
			vm.ip++
			funcNameAddr := int(bytecode[vm.ip])
			funcName := constantPool.RetrieveString(funcNameAddr)
			vm.ip++
			numArgs := int(bytecode[vm.ip])
			vm.ip++
			funcLen := int(bytecode[vm.ip])
			vm.ip++
			funcCode := bytecode[vm.ip : vm.ip+funcLen+(numArgs*2)]
			vm.functionTable[funcName] = &function{name: funcName, args: numArgs, code: funcCode}
			vm.ip += funcLen + (numArgs * 2) - 1
		case opcodes.IFUNC_CALL:
			vm.ip++
			funcNameAddr := int(bytecode[vm.ip])
			funcName := constantPool.RetrieveString(funcNameAddr)
			vm.ip++
			numArgs := int(bytecode[vm.ip])
			vm.ip++
			args := make([]int, numArgs)
			for j := range numArgs {
				args[j] = vm.intStack[len(vm.intStack)-1-j]
			}
			vm.callStack = append(vm.callStack, callFrame{
				functionName: funcName,
				args:         args,
				returnAddr:   vm.ip,
			})
			vm.ip = -1
			bytecode = vm.functionTable[funcName].code
			codeLength = len(bytecode)
		case opcodes.IFUNC_RETURN:
			returnValue := vm.intStack.pop()
			vm.ip = vm.callStack[len(vm.callStack)-1].returnAddr - 1
			vm.callStack = vm.callStack[:len(vm.callStack)-1]
			vm.intStack.push(returnValue)
			if len(vm.callStack) > 0 {
				bytecode = vm.functionTable[vm.callStack[len(vm.callStack)-1].functionName].code
				codeLength = len(bytecode)
			} else {
				bytecode = vm.chunk.Bytecode
				codeLength = len(bytecode)
			}
		// FLOAT
		case opcodes.FPUSH:
			vm.ip++
			addr := int(bytecode[vm.ip])
			value := constantPool.RetrieveFloat(addr)
			vm.floatStack.push(value)
		case opcodes.FADD:
			a := vm.floatStack.pop()
			b := vm.floatStack.pop()
			vm.floatStack.push(b + a)
		case opcodes.FSUB:
			a := vm.floatStack.pop()
			b := vm.floatStack.pop()
			vm.floatStack.push(b - a)
		case opcodes.FMUL:
			a := vm.floatStack.pop()
			b := vm.floatStack.pop()
			vm.floatStack.push(b * a)
		case opcodes.FDIV:
			a := vm.floatStack.pop()
			b := vm.floatStack.pop()
			vm.floatStack.push(b / a)
		case opcodes.FNEG:
			val := vm.floatStack.pop()
			vm.floatStack.push(-val)
		case opcodes.FGT:
			b := vm.floatStack.pop()
			a := vm.floatStack.pop()
			vm.boolStack.push(a > b)
		case opcodes.FLT:
			b := vm.floatStack.pop()
			a := vm.floatStack.pop()
			vm.boolStack.push(a < b)
		case opcodes.FGE:
			b := vm.floatStack.pop()
			a := vm.floatStack.pop()
			vm.boolStack.push(a >= b)
		case opcodes.FLE:
			b := vm.floatStack.pop()
			a := vm.floatStack.pop()
			vm.boolStack.push(a <= b)
		case opcodes.FEQ:
			b := vm.floatStack.pop()
			a := vm.floatStack.pop()
			vm.boolStack.push(a == b)
		case opcodes.FNEQ:
			b := vm.floatStack.pop()
			a := vm.floatStack.pop()
			vm.boolStack.push(a != b)
		case opcodes.FPRINT:
			value := vm.floatStack.pop()
			fmt.Println(value)
		case opcodes.FVAR_BIND:
			vm.ip++
			varnameAddr := int(bytecode[vm.ip])
			varname := constantPool.RetrieveString(varnameAddr)
			value := vm.floatStack.pop()
			vm.symbolTable.bindFloat(varname, value)
		case opcodes.FVAR_LOOKUP:
			vm.ip++
			varnameAddr := int(bytecode[vm.ip])
			varname := constantPool.RetrieveString(varnameAddr)
			value, err := vm.symbolTable.lookupFloat(varname)
			if err != nil {
				panic(err) // TODO: maybe better error handling
			}
			vm.floatStack.push(value)
		// BOOL
		case opcodes.BPUSH:
			vm.ip++
			value := false
			if bytecode[vm.ip] == 1 {
				value = true
			}
			vm.boolStack.push(value)
		case opcodes.BNEG:
			val := vm.boolStack.pop()
			vm.boolStack.push(!val)
		case opcodes.BEQ:
			a := vm.boolStack.pop()
			b := vm.boolStack.pop()
			vm.boolStack.push(a == b)
		case opcodes.BNEQ:
			a := vm.boolStack.pop()
			b := vm.boolStack.pop()
			vm.boolStack.push(a != b)
		case opcodes.BAND:
			a := vm.boolStack.pop()
			b := vm.boolStack.pop()
			vm.boolStack.push(a && b)
		case opcodes.BOR:
			a := vm.boolStack.pop()
			b := vm.boolStack.pop()
			vm.boolStack.push(a || b)
		case opcodes.BXOR:
			a := vm.boolStack.pop()
			b := vm.boolStack.pop()
			vm.boolStack.push((a || b) && !(a && b))
		case opcodes.BPRINT:
			value := vm.boolStack.pop()
			fmt.Println(value)
		case opcodes.BVAR_BIND:
			vm.ip++
			varnameAddr := int(bytecode[vm.ip])
			varname := constantPool.RetrieveString(varnameAddr)
			value := vm.boolStack.pop()
			vm.symbolTable.bindBool(varname, value)
		case opcodes.BVAR_LOOKUP:
			vm.ip++
			varnameAddr := int(bytecode[vm.ip])
			varname := constantPool.RetrieveString(varnameAddr)
			value, err := vm.symbolTable.lookupBool(varname)
			if err != nil {
				panic(err) // TODO: maybe better error handling
			}
			vm.boolStack.push(value)
		// CHAR
		case opcodes.CPUSH:
			vm.ip++
			addr := int(bytecode[vm.ip])
			value := constantPool.RetrieveChar(addr)
			vm.charStack.push(value)
		case opcodes.CEQ:
			a := vm.charStack.pop()
			b := vm.charStack.pop()
			vm.boolStack.push(a == b)
		case opcodes.CNEQ:
			a := vm.charStack.pop()
			b := vm.charStack.pop()
			vm.boolStack.push(a != b)
		case opcodes.CPRINT:
			value := vm.charStack.pop()
			fmt.Println(value)
		case opcodes.CVAR_BIND:
			vm.ip++
			varnameAddr := int(bytecode[vm.ip])
			varname := constantPool.RetrieveString(varnameAddr)
			value := vm.charStack.pop()
			vm.symbolTable.bindChar(varname, value)
		case opcodes.CVAR_LOOKUP:
			vm.ip++
			varnameAddr := int(bytecode[vm.ip])
			varname := constantPool.RetrieveString(varnameAddr)
			value, err := vm.symbolTable.lookupChar(varname)
			if err != nil {
				panic(err) // TODO: maybe better error handling
			}
			vm.charStack.push(value)
		// STRING
		case opcodes.SPUSH:
			vm.ip++
			addr := int(bytecode[vm.ip])
			value := constantPool.RetrieveString(addr)
			vm.stringStack.push(value)
		case opcodes.SEQ:
			a := vm.stringStack.pop()
			b := vm.stringStack.pop()
			vm.boolStack.push(a == b)
		case opcodes.SNEQ:
			a := vm.stringStack.pop()
			b := vm.stringStack.pop()
			vm.boolStack.push(a != b)
		case opcodes.SPRINT:
			value := vm.stringStack.pop()
			fmt.Println(value)
		case opcodes.SVAR_BIND:
			vm.ip++
			varnameAddr := int(bytecode[vm.ip])
			varname := constantPool.RetrieveString(varnameAddr)
			value := vm.stringStack.pop()
			vm.symbolTable.bindString(varname, value)
		case opcodes.SVAR_LOOKUP:
			vm.ip++
			varnameAddr := int(bytecode[vm.ip])
			varname := constantPool.RetrieveString(varnameAddr)
			value, err := vm.symbolTable.lookupString(varname)
			if err != nil {
				panic(err) // TODO: maybe better error handling
			}
			vm.stringStack.push(value)
		// STATEMENTS
		case opcodes.IF:
			if !vm.boolStack.pop() {
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

	fmt.Println()
	fmt.Println("---------------------")
	fmt.Println()

	vm.PrintStacks()
}
