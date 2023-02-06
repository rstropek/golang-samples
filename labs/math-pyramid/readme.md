# Math Pyramid

## Introduction

A math pyramid is a frequently used tool for teaching math to elementary school children. The pyramid's base consists of a configurable number of random values between 0 and 9. Going up, each layer has one value less than the layer below. The values above the base are calculated by adding the left and right neighbours one level below.

Here is an example of such a math pyramid with a base consisting of five values:

```txt
                    ┌──────┐
                    │   70 │
                    └──────┘
               ┌──────┐   ┌──────┐
               │   29 │   │   41 │
               └──────┘   └──────┘
          ┌──────┐   ┌──────┐   ┌──────┐
          │   12 │   │   17 │   │   24 │
          └──────┘   └──────┘   └──────┘
     ┌──────┐   ┌──────┐   ┌──────┐   ┌──────┐
     │    8 │   │    4 │   │   13 │   │   11 │
     └──────┘   └──────┘   └──────┘   └──────┘
┌──────┐   ┌──────┐   ┌──────┐   ┌──────┐   ┌──────┐
│    8 │   │    0 │   │    4 │   │    9 │   │    2 │
└──────┘   └──────┘   └──────┘   └──────┘   └──────┘
```

Your job is to write a program that generates such math pyramids. Math teachers will use your program to generate quizzes for kids. Before handing out the quizzes, they will remove numbers from the pyramid and the kids have to fill in the blanks.

## Requirements

### Level 0: Basic Requirements

* Write a command-line application that generates a math pyramid as described above and print it on *stdout*.

* The application can receive the width of the base as an optional [command line argument](https://gobyexample.com/command-line-arguments).
  * If no argument is given, a width of *5* is used as the default width.
  * If an argument is given, but it is not a number or not between 2 and 10 (including), the program should not generate the pyramid and print a proper error message on *stderr*.

The first version of the program should not care about nice formatting. Printing the raw values on the screen is sufficient. Here is an example output:

```txt
70
29 41
12 17 24
 8  4 13 11
 8  0  4  9  2
```

### Level 1: Packages

Extract different aspects of your application into separate packages.

* Use the [`flag`](https://gobyexample.com/command-line-flags) package for getting the command line arguments.

* Create a package that contains the calculation logic for the math pyramid.

### Level 2: Unit Tests

Add [unit tests](https://gobyexample.com/testing-and-benchmarking) to your modules.

### Level 3: Output Formatting

Add nice output formatting as shown at the beginning of this document. Here are some constants that should make formatting easier:

```go
const top string = "┌──────┐"
const bottom string  = "└──────┘"
const separatorLength = 3
const separatorChar byte = ' ' // Note that this is a byte, not a rune
const bar = '│'
```

* Create a package that contains helper functions and constants for output formatting.
  * Don't forget to add unit tests.
