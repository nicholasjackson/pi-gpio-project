# Building a GPIO based Raspberry Pi project controllable from Google Assistant
In this project we are going to build a small project for the Raspberry Pi which allows us to control some LEDs using Google Assistant, the 

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
	rpi.SO_51.Out(gpio.High)


	// Continuous loop blocking exit
  for {

  }
```

The line `rpi.SO_51.Out(gpio.High)` will set `Pin 15` on the Pi to High, the following lines ensure that our application will not exit until a key press is detected.  We also need to add two more package imports to our code:

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

To cycle the pin high and low we are going to create an endless loop which continuously flips the state and then waits for a set interval.  First let's create a new `struct` which will encapsulate this behavior.  We are going to add this after our `func main`.

The first thing we need to do is to define a new struct with the following fields:

```go
type PinCycle struct {
	Pin     gpio.PinIO
	Running bool
}
```

The struct has two fields `Pin` which is of type `gpio.PinIO` and `Running` of type `bool`, the `Pin` field will hold a reference to our Raspberry Pi pin and `Running` will hold the state if the pin is cycling between the high and low state.

We can then add a method to this struct which will allow us to flip the pin state:

```go
func (f *PinCycle) Cycle() {
	go func() {
		f.Running = true
		state := gpio.High

		for f.Running {
			f.Pin.Out(state)

			time.Sleep(500 * time.Millisecond)

			// flip the state
			if state == gpio.High {
				state = gpio.Low
			} else {
				state = gpio.High
			}
		}
	}()
}
```

Let's examine what is going on in this block.

First is that we are using a new keyword `go`, the `go` keyword will execute a statement and will not wait for it to return, since we are using a continuous loop inside of our method if we did not use the `go` keyword then the `Cycle` method would block forever and nothing beneath it in our program would execute.  When using the `go` keyword we could execute a single statement, for example `go logger.Println("Hello World")`, however when we need to execute multiple statement we can wrap these into an inline function with `go func() {//...}()`.

In the next line we are setting the running state to true with the line `f.Running = true`, `Running` is a field on our `PinCycle` struct, we use the `f.` accessor as this is what we defined in the method signature as a variable to attach the method to the struct, `func(f *PinCycle) Cycle()`.

We are then setting the initial pin state to `gpio.High`, we will use this variable to cycle the pin between high and low states and this will cause our LED to flash.

Next we start a loop `for f.Running { // ...}`, this loop will run as long as the value of `f.Running` is set to `true`, this is also the code which would block our program from executing future statements had we not wrapped it into a `go func() {}()` block.

We then set our pins output state to our current state, for the first iteration of the loop this will be `gpio.High`, before we change the state again we are going to `Sleep` for half a second.  To sleep we are going to use the `Sleep(time.Duration)` function from the `time` Go standard package.  This function has a single parameter which is a duration, we can create a duration simply by multiplying an integer by a time.Duration type, `500 *time.Millisecond)`.

After we have slept we need to cycle the pin state before repeating our loop, to do this we are using an `if` statement, an `if` statement allows us to do one thing if the condition in the statement is true and something else when it is false.

```go
if state == gpio.High {
	state = gpio.Low
} else {
	state = gpio.High
}
```

We are checking if the state is currently `gpio.High` and if so we set it to `gpio.Low`, if it is `gpio.Low` we then set things to `gpio.High`.

When we call the `Cycle()` method on our `PinCycle` struct the pin state will change every half a second and when there is a LED attached to the pin this will cause it to flash.

Since we are now using the `time` package add this to your list of imported packages at the top of your file.  We can then delete the line `rpi.P1_15.Out(gpio.High)` create an instance of our `PinCycle` struct and call the `Cycle` method.

```go
	p15 := PinCycle{
		Pin:     rpi.SO_51,// GPIO 15
		Running: false,
	}
	p15.Cycle()
```

When we are creating our `PinCycle` object we are passing it a reference to the Pin we would like to cycle and also setting the initial state for `Running` to be `false`.

Before we continue, let's also add a `Stop` method, we will need this later:

```go
func (f *PinCycle) Stop() {
	f.Running = false
}
```

All we are doing here is setting the internal Field to false, this will make the loop in the `Cycle` method exit.

If you run your app gain `go run main.go` you should now see pin 15 flashing every 0.5 seconds. 

```bash
######### Raspbery Pi Model Zero GPIO simulator ###########

 +| +| -|14|15|18| -|23|24| -|25|08|07|01| -|12| -|16|20|21
         ## ## ##    ## ##    ## ## ## ##    ##    ## ##
   ## ## ##    ##    ##    ## ## ##    ## ## ## ## ## ##
 +|02|03|04| -|17|27|22| +|10|09|11| -|00|05|06|13|19|16| -


Log (log is also written to ./out.log:
10:38:30.910722 Hello World
```

Lets add some more pins to our application, we are going to use pins `14`,`15`,`18`,`23`, `24`, and `25`, which correspond to `SO_51`, `SO_53`, `SO_64`, `SO_77`, `SO_81`, `SO_83`.  Define the new `PinCycle` structs and call the `Cycle` method for your new pins.

You should have defined something like the following in your code:

```go
	p14 := PinCycle{
		Pin:     rpi.SO_51,
		Running: false,
	}

	p15 := PinCycle{
		Pin:     rpi.SO_53,
		Running: false,
	}

	p18 := PinCycle{
		Pin:     rpi.SO_63,
		Running: false,
	}

	p23 := PinCycle{
		Pin:     rpi.SO_77,
		Running: false,
	}

	p24 := PinCycle{
		Pin:     rpi.SO_81,
		Running: false,
	}

	p25 := PinCycle{
		Pin:     rpi.SO_83,
		Running: false,
	}

	p14.Cycle()
	p15.Cycle()
	p18.Cycle()
	p23.Cycle()
	p24.Cycle()
  p25.Cycle()
  ```

  When you now run `go run main.go` you should see 6 of the pins flashing.

  This is somewhat boring though, what would be nice is if they flashed at different rates, lets make a little change to our code to change the flash duration.

  We are going to use the random standard package `rand` to create a random duration between 300 milliseconds and 1000 milliseconds, to do this we use the function `rand.Intn(max int)` which returns a random number up to the specified maximum.  To get a number in a range we generate a number which is our max with the minimum subtracted then add the minimum.

  Remember to add the package `rand` to your imports.

  Change the line in your `Cycle` method, `time.Sleep(500 * time.Millisecond)` and replace it with the following two lines.

  ```go
sleepDuration := rand.Intn(1000-300) + 300
time.Sleep(time.Duration(sleepDuration) * time.Millisecond)
```

Now when run `go run main.go` you should see 6 of the pins flashing at different rates.

Ok we are nearly there, there is one last step before we can control our application with Google Assistant, rather than just start the application with `go run` we need to be able to start and stop it based on a HTTP request.  To do this we are going to implement a simple HTTP server.  The Go standard library again has fantastic capability for this using the `http` package, we can register routes such as "/" which will trigger a function and it only requires a single command to start the server.

First we need to map a route which will trigger a function, to do this we use the function `http.HandleFunc(path string, handler func(http.ResponseWriter, *http.Request))`.

This function has two parameters the first is the path which will trigger the function provided in the second parameter.  To determine if  the LEDs are to be switched on or off, we are going to read a query string variable.  Add the following code block which will replace your code where you have `p14.Cycle()`, etc.

```go
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("mode") == "on" {
			p14.Cycle()
			p15.Cycle()
			p18.Cycle()
			p23.Cycle()
			p24.Cycle()
			p25.Cycle()
		} else {
			p14.Stop()
			p15.Stop()
			p18.Stop()
			p23.Stop()
			p24.Stop()
			p25.Stop()
		}
	})
```

If the webserver is called with the query string `?mode=on` then we will active the LEDs if it is called with `?mode=off` then we deactivate them.

All we need to do now is to start the server with the function `http.HttpListenAndServe(":9000",nil)`, this will start a webserver which is accessible on port `9000`.  Add the following line to your `func main` just before the `for` statement.

```go
	http.ListenAndServe(":9000", nil)
```

We can now test our application start the application using `go run main.go` and then in another terminal window make a call to the webserver using curl.

```go
# Start LEDs cycling
curl "http://localhost:9000/?mode=on"

# Stop LEDs cycling
curl "http://localhost:9000/?mode=off"
```

That is all we need to do, our application is almost complete, in the next steps we will hook things up to the Google Assistant.