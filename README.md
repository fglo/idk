﻿# IDK Programming Language

IDK is a statically typed language. Its interpreter is currently implemented in Go.

This is very much a work in progress, so ANYTHING can change any moment.

**Note:** I started with almost zero knowledge of language development and interpreter writing. This project is about trying my own ideas and testing my intuition. Maybe the next one will be backed by an actuall knowledge about languages, interpreters and compilers. It started with experimenting with some ideas in Python, now it's based on a book "writing an INTERPRETER in go" by Thorsten Ball

## The Name

IDK means literally I Don't Know. The name is a placeholder. If I decide to work on it further I will think of a better name.

## Features

- int, char and bool variable types
- declassign operator: `:=`
- declare operator: `:`
- assignment operator: `=`
- arithmetic operators: `+`, `-`, `*`, `/`
- comparison operators: `==`, `>`, `<`, `>=`, `<=`
- logical operators: `not`, `and`, `or`, `xor`
- printing
- conditional statements: `if`, `if-else`, `if-else-if`
- comments
- for loops (like while loops)

#TODO (must-have):
- actual for loops
- arrays
- range operators: `..`, `..=`
- array `in` operator
- parentheses in operations
- loops: while and for with variable
- nested loops
- strings (arrays?)
- procedures
- functions
- structs
- tokenization not only with spaces but also with operators
- better variables, statically typed

#TODO (maybe):
- mutable/immutable variable modifiers
- ternary operators or oneline if expressions: `i := 1 < 2 ? true : false` or `i := if 1 < 2 then true else false`
- python-like comprehensions (generators)
- lambdas/function literals
- c#-like extension methods 

## Syntax

### Comments

You can write comments in code using the `//`:
```
// comment test
to_print := 1 //assigning variable
print to_print // printing
```

### Variable Declaration and Assignment

You can declare a variable using the `:` operator:
```
a:int
b:char
c:bool
d:string
```

You can assign value to a declared variable using the `=` operator:
```
a = 67
b = 'b'
c = true
d = "false"
```

You can assign value in the same line as you declared it:
```
a:int = 67
b:char = 'b'
c:bool = true
d:string = "false"
```

You can declare variable and assign value to it using the `:=` operator (it's type will be infered):
```
a := 67
b := 'b'
c := true
d := "false"
```

### Arithmetic operators

Calculations are done using normal set of operators: `+`, `-`, `*`, `/`:
```
a := 1 + 2
b := 3 - a
c := 4 * b
d := d / 5
```

### Comparison operators

IDK supports following comparison operators: `==`, `>`, `<`, `>=` and `<=`:
```
eq := 1 == 1
gt := 2 > 1
gte := 2 >= 1
lt := 1 < 2
lte := 1 <= 2
```

### Logical operators

Following logical operators are supported: `xor`, `or`, `and`, `!`:
```
negated := !true
anded := true and false
ored := true or false
xored := true xor false
```

### Range operators // TODO: range operators

IDK supports exclusive range operator `..` and inclusive range operator `..=`:
```
exclusive := 1..3  // creates following array: [1, 2]
inclusive := 1..=3 // creates following array: [1, 2, 3]
```

### Arrays // TODO: arrays

Currently IDK supports only arrays of digits created by the range operators. The range operator is exclusive.
```
x := 1..3          // creates following array: [1, 2]
```

IDK supports only one operation on arrays : `in`.  
```
print 3 in 1..=3   // prints true

x := 1..3
print 3 in x       // prints false
```

### Printing

You can print variables, integers, character and expressions using the `print` keyword: // TODO: add `print` keyword
```
a := 'a'
a_code := a + 0

print(a) // prints a
print(a, 'a') prints a a

print a            // prints a
print a_code       // prints 97
print 'a'          // prints a
print 'a' + 0      // prints 97
print 'a' = 'a'    // prints true
print 'a' = 2      // prints false
print true         // prints true
print false        // prints false
print 1..6         // prints 1, 2, 3, 4, 5
print 6..1         // prints 6, 5, 4, 3, 2
```

### Conditional statements

You can write `if`, `if-else` and `if-else-if` statements:
```
a := 67
b := 't'
c := 'f'

if a == 67
    print b
end

if a < 68
    print b
end

if a > 67
    print c
else
    print b
end

if a > 67
    print c
else if a > 68
    print c
else
    print b
end

if true
    print 't'
end

if false
    print 'f'
end

```

You can also nest `if` statements:
```
if 3 >= 2
    if 3 > 4
        print '0'
    else
        print '1'
    end
    print '2'
end
```

Logical operators are also supported:
```
if 1 < 2 and 2 > 1
    print 1
end
if 1 < 2 or 2 < 1
    print 2
end
if 1 < 2 xor 2 < 1
    print 3
end
if not 1 < 2 xor 2 > 1
    print 4
end
```

### Loops

#### For loop

IDK currently supports only while loops (using the `for` keyword):

```
i := 0
for i < 3
    i = i + 1
    print 'x'
end
```

// TODO: what's below

IDK currently supports only simple for loop using range operators `..`:
```
for 1..3
    print 'x'
end
```

Each for loop implicitly declares variable `_it` which contains current iterator value:
```
for 1..5
    print _it
end
```

// TODO: update the rest of the README

Currently loops cannot be nested.

## Running programs

You can run an IDK program by calling the idk.py script and providing a path to a file with an IDK program:
```bash
python idk.py test.idk
```

## Running interactive interpreter

IDK interpreter provides simple interactive mode. Only oneline statements are currently supported.

You can run it using following command:
```bash
python idk.py -it
```

if everything went as it should, you should see following text and prompt:
```
Welcome to IDK interactive!
$ 
```

If you want to exit just type `exit` and click enter:
```
$ print 1
1 
$ exit
```

Example:
```
$ x := 2
$ print x + 3
5
$ if true
> i := 1
> print i
> end
1 
$ for 1..=3
> print _it
> end
1
2
3
$ exit
```