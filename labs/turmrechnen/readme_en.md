# Tower Math

## Requirements

Your task is to write a command-line application in Go for [*Tower Math*](http://www.floriangeier.at/schule/kopf/kopf.php).

A user provides two parameters to your application via the command line:

* The starting value (e.g., *9*); the starting value must be > 1.
* The "height" (e.g., *5*); the height must be > 2.

You should output the result in the following format:

```txt
   9 * 2 = 18
  18 * 3 = 54
  54 * 4 = 216
 216 * 5 = 1080
1080 / 2 = 540
 540 / 3 = 180
 180 / 4 = 45
  45 / 5 = 9
```

If the user provides incorrect parameters on the command line or if parameters are missing, output an appropriate error message on the screen.

You can ignore *overflows*.

## Tips

* Accessing command-line parameters:
  * [Manual](https://gobyexample.com/command-line-arguments)
  * [Using *Flags*](https://gobyexample.com/command-line-flags)
* You can pad outputs to the left with spaces by specifying format options with `Printf` (and variants of this function). For example, `fmt.Printf("|%6d|\n", 42)` outputs four spaces followed by *42*, to reach six total characters.
* You can convert strings to integers using [`strconv.Atoi`](https://golang.org/pkg/strconv/).

## Levels

Depending on prior knowledge, some people may find this task easier or more difficult. Here are suggestions for how to solve the example step by step. Everyone can work through as many levels as they find appropriate for their programming experience.

### Level 1 - Calculation Logic

* Assume starting value and height, and omit command-line parameters.
* Output only the intermediate results. For example:

```txt
9
18
54
216
1080
540
180
45
9
```

### Level 2 - Command-Line Parameters

* Add the ability to pass parameters via the command line.
* Remember to validate the parameters and output appropriate error messages.

### Level 3 - Improved Output

* Improve the output so that the result looks like the one required at the beginning of this description.
* If possible, align the output value to the right.

### Level 4 - Unit Testing

* Structure your code in a way that makes it easy to test.
* Write at least one meaningful unit test ([Quick guide](https://gobyexample.com/testing)).

