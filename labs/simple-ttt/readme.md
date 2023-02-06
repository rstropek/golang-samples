# Exercise: Simple TicTacToe

## Requirements

Implement a simple, console-based TicTacToe game.

The console app displays the board and asks the user where she wants to set (format *A1*, *B2*, etc.). After each set, the board is displayed again and the app asks for another location. This process has to continue until we have a winner.

The following sample output illustrates the concept:

```txt
┏━━┯━━┯━━┓
┃  |  |  ┃
┠──┼──┼──┨
┃  |  |  ┃
┠──┼──┼──┨
┃  |  |  ┃
┗━━┷━━┷━━┛

Enter location (A1-C3): A1
┏━━┯━━┯━━┓
┃XX|  |  ┃
┠──┼──┼──┨
┃  |  |  ┃
┠──┼──┼──┨
┃  |  |  ┃
┗━━┷━━┷━━┛

Enter location (A1-C3): B2
┏━━┯━━┯━━┓
┃XX|  |  ┃
┠──┼──┼──┨
┃  |OO|  ┃
┠──┼──┼──┨
┃  |  |  ┃
┗━━┷━━┷━━┛

Enter location (A1-C3): A2
┏━━┯━━┯━━┓
┃XX|XX|  ┃
┠──┼──┼──┨
┃  |OO|  ┃
┠──┼──┼──┨
┃  |  |  ┃
┗━━┷━━┷━━┛

Enter location (A1-C3): A3
┏━━┯━━┯━━┓
┃XX|XX|OO┃
┠──┼──┼──┨
┃  |OO|  ┃
┠──┼──┼──┨
┃  |  |  ┃
┗━━┷━━┷━━┛

Enter location (A1-C3): C2
┏━━┯━━┯━━┓
┃XX|XX|OO┃
┠──┼──┼──┨
┃  |OO|  ┃
┠──┼──┼──┨
┃  |XX|  ┃
┗━━┷━━┷━━┛

Enter location (A1-C3): C1
┏━━┯━━┯━━┓
┃XX|XX|OO┃
┠──┼──┼──┨
┃  |OO|  ┃
┠──┼──┼──┨
┃OO|XX|  ┃
┗━━┷━━┷━━┛

Winner: O
```

## Tips

Here are the characters to draw the board:

```go
top := []rune("┏━┯┓")
middle := []rune("┠─┼┨")
bottom := []rune("┗━┷┛")

// This is how you print a rune:
fmt.Printf("%c", top[0])
```

You can read from *stdin* as follows:

```go
var userName string
fmt.Println("Enter your name:")
fmt.Scanln(&userName)
fmt.Println("Your name is", userName)```
```

## Advanced Version

Try to separate UI logic (input from *stdin*, output to *stdout*) nicely from calculation logic.
