# IDK Programming Language

IDK is a statically typed language with immutable variables. Its interpreter is currently implemented in Python.

This is very much a work in progress, so ANYTHING can change any moment.

## The Name

IDK means literally I Don't Know. The name is a placeholder. If I decide to work on it further I will think of a better name.

## Features

- int, char and bool variable types
- addition
- subtraction
- multiplication
- division
- printing
- conditional statements

#TODO: 
- logical (`not`, `and`, `or`, `xor`) keywors/operators
- local variables (better scopes)
- comments
- loops
- strings
- procedures
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
c := true
d := false
```

Only global variables are currently supported (you can't declare variables inside an `if` statement).

### Calculations

Calculations are done using normal set of operators: `+`, `-`, `*`, `/`:
```
a := 1 + 2
b := 3 - a
c := 4 * b
d := d / 5
```

You can also chain operators:
```
e := 3 + 4 - 5
f := 2 + 1 * 2
g := 3 - 4 / 2
```

IDK keeps chars and bools as integers (chars as ASCII code values and bool using convention that true is 1 and false is 0), so it's possible to do calculations on them (and doing so convert them to integers):
```
a_code := 'a' + 0
true_val := true + 0
false_val := false + 0
```

### Comparison

IDK supports following comparison operators: `=`, `>` and `<`:
```
eq := 1 = 1
gt := 2 > 1
lt := 1 < 2
```

### Logical operators

Following logical operators are supported: `xor`, `or`, `and`, `not`:
```
negated := not true
anded := true and false
ored := true or false
xored := true xor false
```

### Printing

You can print variables, integers, character and expressions using the `print` keyword:
```
a := 'a'
a_code := 'a' + 0

print a
print a_code
print '0'
print 'a'
print 'a' + 0
print 'a' = 'a'
print true
print false
```

### Conditional statements

You can write `if`, `if-else` and `if-else-if` statements:
```
a := 67
b := 't'
c := 'f'

if a = 67
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
if 3 > 2
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

You can assign variables and do oneline calculations. if you want to exit just type 'exit' and click enter:
```
$ print 1
1 
$ x := 2
$ print x + 3
5 
$ exit
```