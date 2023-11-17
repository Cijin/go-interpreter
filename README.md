# go-interpreter
Writing an interpreter in go for the monkey programming language.
Following along this [book](https://interpreterbook.com/) from Thorsten bell.

I'll be documenting the process on [seagin.me](https://www.seagin.me/2023/october)

October 2023: Lexer & Parser complete

## Running the REPL

The monkey programming language comes with a REPL. If you wish to run it yourself:

### Prerequisites

* [Go](https://go.dev/dl/) version `1.6` or above


### Installing
1. Clone the repo
2. `cd go-interpreter`
3. Run `go build -o monkey` (yes, that's it)
4. Now you can run the monkey repl in your terminal by running: `./monkey`

If everything went well you should see this in your terminal:

```
Welcome to monkey v0.0.0
Press ctrl-d to exit.
>> 5 + 5
10
>> 5 + true;
type mismatch: INTEGER + BOOLEAN
>> 

```

### v0.0.1

The interpreter has been expteded to be able to evaluate `function` literals and `Call` expressions. 
So you interpret inputs like this:

```
Welcome to monkey v0.0.1
Press ctrl-d to exit.
>> let newAdder = fn(x) { fn(y) { x + y } };
>> let addTwo = newAdder(2)
>> addTwo(3)
5
>> 
```
