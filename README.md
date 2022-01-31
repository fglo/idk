# IDK Programming Language

IDK is a dynamically typed language with immutable variables. Its interpreter is currently implemented in Python.

This is very much a work in progress, so ANYTHING can change any moment.

## The Name

IDK means literally I Don't Know. The name is a placeholder. If I decide to work on it further I will think of a better name.

## Features

- int, char and bool variable types
- assignment operator: `:=`
- arithmetic operators: `+`, `-`, `*`, `/`
- comparison operators: `=`, `>`, `<`, `>=`, `<=`
- logical operators: `not`, `and`, `or`, `xor`
- range operators: `..`, `..=`
- array `in` operator
- printing
- conditional statements: `if`, `if-else`, `if-else-if`
- comments
- for loops

#TODO (must-have):
- parentheses in operations
- loops: while and for with variable
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
- lambdas
- c#-like extension methods 

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

### Comments

You can write comments in code using the `//`:
```
// comment test
to_print := 1 //assigning variable
print to_print // printing
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

### Arithmetic operators

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

### Comparison operators

IDK supports following comparison operators: `=`, `>`, `<`, `>=` and `<=`:
```
eq := 1 = 1
gt := 2 > 1
gte := 2 >= 1
lt := 1 < 2
lte := 1 <= 2
```

### Logical operators

Following logical operators are supported: `xor`, `or`, `and`, `not`:
```
negated := not true
anded := true and false
ored := true or false
xored := true xor false
```

### Range operators

IDK supports exclusive range operator `..` and inclusive range operator `..=`:
```
exclusive := 1..3        // creates following array: [1, 2]
inclusive := 1..=3        // creates following array: [1, 2, 3]
```

### Arrays

Currently IDK supports only arrays of digits created by the range operators. The range operator is exclusive.
```
x := 1..3        // creates following array: [1, 2]
```

IDK supports only one operation on arrays : `in`.  
```
print 3 in 1..=3 // prints true

x := 1..3
print 3 in x    // prints false
```

### Printing

You can print variables, integers, character and expressions using the `print` keyword:
```
a := 'a'
a_code := 'a' + 0

print a         // prints a
print a_code    // prints 97
print 'a'       // prints a
print 'a' + 0   // prints 97
print 'a' = 'a' // prints true
print 'a' = 2   // prints false
print true      // prints true
print false     // prints false
print 1..6      // print 1, 2, 3, 4, 5
print 6..1      // print 6, 5, 4, 3, 2
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

IDK currently supports only simple for loop using range operators `..`:
```
for 1..3
    print 'x'
end
```

Each for loop implicitly declater variable `_it` which contains current iterator value:
```
for 1..5
    print _it
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