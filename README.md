# IDK Programming Language

IDK is a weakly typed language with immutable variables. It's interpreter is currently implemented in Python.

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

## Syntax

Currently all tokens in IDK must be separated by the space symbol: ' '.

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
h := 2 + 1 *2
i := 3 - 4 / 2
```

### Printing

You can print variables, integers, character and expressions using the `print` keyword:
```
print a
```

## Running programs

You can run an IDK program by calling the idk.py script and providing a path to a file with an IDK program:
```bash
python idk.py test.idk
```