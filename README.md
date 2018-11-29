# Building a GPIO based Raspberry Pi project

Go allows you to compile an application for multiple languages from a single source.  For the first part of this tutorial we will be working in a Linux environment. The first step is to create the `main.go` application entry point, let's create that file in the current folder.

```go
package main

func main() {

}
```

Once we have that file we can run it using the command `go run main.go`, run that now in your terminal.

```bash
$ go run main.go
```

When you run `go run main.go` the Go run time will compile a temporary version of your application and then execute it, we will not see any output from this run as the `main` func has no content.

An executable Go application always needs a `main` package, the package name is always the first line in the file and normally corresponds to the folder name.  In Go you can break up your application into packages by creating `.go` files inside of folders.

Let's add some content to our `func main`, we are going to add a logger which will log output to the command line, to do this we will use the `log` package which is part of Go's standard library.  Add the following lines to your `main.go`.

```go
  func main() {
    logger := log.New(os.Stdout, "", log.Lmicroseconds)
    logger.Println("Hello World")
  }
```

The line `logger := log.New(os.Stdout, "", log.Lmicroseconds)` will create a new logger and assign it to the variable `logger` in Go there is a difference between `:=` and `=`, `:=` is a convenience operator it will create a new variable and assign it a value in one line.  The operator `=` would only assign a value to an existing variable.  We could have written the same line as separate variable creation and assignments like so:

```go
var logger *log.Logger
logger = log.New(os.Stdout, "", log.Lmicroseconds)
```

What we are doing in this example is using the `var` keyword which defines a variable, we give it the name `logger` and define the type which is `*log.Logger` this is the object `Logger` from the `log` package and the `*` marks this as a reference which is allocated to the global memory area the `heap` rather than the local store the `stack`.

When we call `log.New(...)` what we are actually doing is calling a function `New` which is present in the `log` package, this function has the signature `func New(out io.Writer, prefix string, flag int) *Logger`.

Let's take a look at the parameters which are passed to this function.

The first is `out` this is defined as an io.Writer which is a Go interface, in Go you have `stucts` and `interfaces`, `interfaces` define behavior where as `structs` implement them.  An interface itself has no functionality it only defines behavior and even though this is a parameter in the function we would actually pass a reference to a `struct` which contains the methods defined in the interface and therefore implements its behavior.

By writing functions and methods which depend on `interfaces` rather than `structs` we are not being prescriptive on the concrete object which the function or method expects.

Looking at the definition for `io.Writer` we see it is:

```go
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

What we are defining with this `type` statement is a behavior a `struct` which is compatible with `Writer` must implement the method `Write` with the same signature as defined in the interface declaration.

In our code we are passing `os.Stdout`, the `Stdout` variable in the `os` package, holds a reference to a `struct` which implements the `Write` method, therefor it implements the behaviour of the `Writer` interface and can be passed as parameter to the `New` function.  This may feel a little confusing at first but it does not get much more complicated.  Before we continue any further Let's take a look at something which I have mentioned before `functions` and `methods`.

Go distinguishes between `functions` and `methods` however there is not a great deal of complexity, a `function` is attached to a package, in our code `func main()` is a function, a method is a function which is attached to a struct.  For example see the following code block which defines an `interface` with a `SayHello` method and a `struct` which implements this.

```go
type Greeter interface {
  SayHello()
}

type Person struct {}

func (p *Person) SayHello() {

}
```

To add a `method` for a struct you just write a normal `function` signature but you add the additional syntax `(p *Person)` after the `func` declaration and before the method name.  Where some languages use an encapsulated style such as:

```java
class Person {
  SayHello() {
    // ... do somethings
  }
}
```

Go uses the `([ref] *[StructType])` syntax to attach methods to `structs` which is not encapsulated inside the `struct` declaration.  `struct`s can have have variables which are defined by encapsulating them in the `struct` definition.

```go
type Person struct {
  Greeting string
}

func (p *Person) SayHello() {
  logger.Println(p.Greeting)
}
```

An example of using this struct could be like the following example:

```go
p := Person{} // create a new person, note the {}
p.Greeting = "Hello World" // assign the greeting to the Greeting field on the Person struct
p.SayHello() // call the SayHello method on the struct
```

When we are creating a new struct we always use the curly brackets `{}`, in fact rather than setting the `Greeting` field on a separate line we could have written the code block as following:

```go
p := Person{
  Greeting: "Hello World",
} // create a new person, and set Greeting
p.SayHello() // call the SayHello method on the struct
```

`p` is now a value type or a `stack` allocated instance of `Person` the `stack` is memory local to our function and is incredibly fast and efficient, when the function exits all the `structs` and variables on the `stack` are automatically cleaned up.  Often however, we need to pass `structs` between functions or methods.  If we attempted to pass a `value` type it would be passed as a `copy` and any changes to state would only exist in the local copy.  Instead we often pass a `reference` this is what is happening with the block `New()` for our logger.  Creating a reference in Go is really easy, we can either do this at create time by adding an `&` to the beginning of our `struct` name. Or we can take a reference at a later point again using the `&` operator.

```go
p := Person{} // p is a value type of Person
p1 := Person{} // p1 is now a reference to a Person struct
p2 := &p // p2 is a reference to the variable p
```

Unlike C++ we do not need to use a different operator to call methods on values or references, in fact you can mostly ignore what they are, we only really need to understand if something is a variable or a reference when we are passing or receiving objects.

```go
a := "abc"
b := &a
*b = "123"

// what does a now equal?
// 123 because b is a reference to a
```

Last bit of theory before we dive into creating our application, note that to set the value `123` we had to use the syntax `*b` this is because b is not a `value` of type `string` but a `reference` to `a` which is of `type` `string`, in order to set a value for `b` (or a) we must dereference it which we do using the `*` operator.  If you had tried to set `b = "123"` the compiler would have complained with the error `cannot use "123" (type string) as type *string in assignment`.

To learn more about the basics of Go I recommend you take a look at the excellent tutorial from the Go team: [https://tour.golang.org/welcome/1](https://tour.golang.org/welcome/1).

Ok back to our application, our source code should now look like the following:

```go
package

func main() {
	logger := log.New(os.Stdout, "", log.Lmicroseconds)
	logger.Println("Hello World")
}
```

In this example we are using two packages `log` and `os`, in order to use a package in Go you need to add it to an `import` statement, many IDE's will do this for you when you save your file.  Let's add the import manually now.  Add the following lines to your source code.

```go
package main

import (
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.Lmicroseconds)
	logger.Println("Hello World")
}
```

Imports need to be added to every file where you are using a package if we had another source file `foo.go` and we used the `log` package we would also need to add `import` there too.  Let's now run our program and see what happens `go run main.go`.

```bash
$ go run main.go
16:46:15.134393 Hello World
```

This time we should see the message "Hello World", with the basics out of the way lets start writing our go programme for our Raspberry Pi.

To interact with the `GPIO` interface on the Pi we are going to use the `gpio` package from `perif.io`, [https://periph.io/](https://periph.io/). `perif.io` provides low level capabilities to interact with a number of different devices including the Raspberry Pi.  The first thing we need to do is to initialize `perif.io` we do this with the following code.

```go
	// Load all drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
  }
```

Add this code block beneath your logger.Println line, you will also need to import the `perif.io` host package, since this is not in the standard library we have to use a `URI` based syntax, often this corresponds to a `github` repository and path.

```go
package main

import (
	"log"
	"os"

	"periph.io/x/periph/host"

)

func main() {
	logger := log.New(os.Stdout, "", log.Lmicroseconds)
	logger.Println("Hello World")

	// Load all drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
}
```

Since we do are not coding onto a Raspberry Pi at the moment we can use a simulator which is a direct replacement for the Pi but will show us graphically the output of the pins.  Let's add some code to set one of the pins to `High` this would send voltage to the pin and if a LED was attached then the LED would illuminate.

Add the following code to your example:

```go
	rpi.P1_15.Out(gpio.High)

  // register for key presses
	c := make(chan os.Signal, 1)
	signal.Notify(c)

	// Block until a signal is received.
	<-c
```

The line `rpi.P1_15.Out(gpio.High)` will set `Pin 15` on the Pi to High, the following lines ensure that our application will not exit until a key press is detected.  We also need to add two more package imports to our code:

```go
  // GPIO simulator for the Raspberry Pi, compatible with periph.io rpi package
  "github.com/nicholasjackson/periph-gpio-simulator/host/rpi"
  // GPIO functionality for periph.io
	"periph.io/x/periph/conn/gpio"
```

When you run your application again you should now see the following output:

```bash
######### Raspbery Pi Model Zero GPIO simulator ###########

 +| +| -|14|15|18| -|23|24| -|25|08|07|01| -|12| -|16|20|21
         ## ## ##    ## ##    ## ## ## ##    ##    ## ##
   ## ## ##    ##    ##    ## ## ##    ## ## ## ## ## ##
 +|02|03|04| -|17|27|22| +|10|09|11| -|00|05|06|13|19|16| -


Log (log is also written to ./out.log:
17:05:52.682209 Hello World                                      
```

The Pin 15 should be highlighted in Red showing it's High state, all the other pins will be white.  Let's modify things a little further to cycle this pin on and off.

