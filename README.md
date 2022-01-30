# IDK Programming Language

IDK is a dynamically typed language with immutable variables. Its interpreter is currently implemented in Python.

This is very much a work in progress, so ANYTHING can change any moment.

## The Name

IDK means literally I Don't Know. The name is a placeholder. If I decide to work on it further I will think of a better name.

## Features

- integer variables
- character variables
- addition
- subtraction
- multiplication
- division
- printing
- conditional statements

#TODO: 
- loops
- strings
- prcedures
- functions
- tokenization not only with spaces but also with operators

## Syntax

Currently all tokens in IDK must be separated by a space symbol: ' '.

Correct code:
```
a := 1 + 2
```

Incorrect code:
```
a:=1+2
```
### Assignment

You can assign value to a new variable using the `:=` operator:
```
a := 67
b := 'b'
```

Only global variables are currently supported (you can't declare variables inside an `if` statement).

### Calculations

Calculations are done using normal set of operators: `+`, `-`, `*`, `/`:
```
c := 1 + 2
d := 3 - a
e := 4 * b
f := d / 5
```

You can also chain operators:
```
g := 3 + 4 - 5
h := 2 + 1 * 2
i := 3 - 4 / 2
```

### Printing

You can print variables, integers, character and expressions using the `print` keyword:
```
print a
```

### Conditional statements

You can write simple `if`, `if-else` and `if-else-if` statements:
```
if a < 68
    print b
end

if a > 67
    print c
else
    print d
end

if a > 67
    print e
else if a > 68
    print f
else
    print g
end
```

## Running programs

You can run an IDK program by calling the idk.py script and providing a path to a file with an IDK program:
```bash
python idk.py test.idk
```
