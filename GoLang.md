# Go-Lang

## Interfaces

```go
type InterfaceName interface {
    // convention : ....er interface // Worker/Marshaler
    Method1()
}
type Type1 struct {
    val1 string
}
type Type2 struct {
    val2 string
}
// aka Polymorphism
func(t Type1) Method1() {
    fmt.Println("Method1 type1", t)
}
func(t *Type2) Method1() {
    fmt.Println("Method1 type2", t)
}
func main() {
    t1 := Type1{"value"}
    t2 := Type2{"value2"}
    var ii InterfaceName = t1
    ii.Method1();
    ii = &t1
    ii.Method1()
    ii = &t2
    ii.Method1()
    ii = t2
    ii.Method1() // Not possible -> compilation error as pointer the concrete value of interface is not addressable = Now there is no way to get &t2 from ii, hence this doesn't work
    // empty interface
    var i interface{}
    // Type Assertion
    i = 56
    int_from_i = i.(int) // 56
    string_from_i := i.(string) // panic
    string_from_i, ok := i.(string) // "", false
    // Type switch
    // val contains the value of the interface variable and not the type
    switch val := interface_variable.(type) {
        case Type1:
    case InterfaceName:
    case Type2:
    default:
    }
}
```

- If a type implements all methods of an interface it is said to be implementing that interface and an instance of that type can be assigned to an instance of the interface.
- All types implement empty interface, since no methods for empty interface(can be used for arbitrary arguments to methods.)
- Type assertions -> Extract the given type from the interface
- Type Switch -> switch case for interface type
- Pointer types implementing interfaces
- A type can implement multiple interfaces.
- Inheritance can be achieved using embedding interfaces(same as types)

## Inheritance(sort of)

```go
type parent1 struct {
    field1 string
}
type parent2 struct {
    field1 string
}
type parent3 struct {
    field3 string
}
type ii interface() {
    Method1() string
}
func (p parent1) Method string {
    return p.field1()
}
type child struct {
     // Also known as Embedding
    ii
    parent1
    parent2
    parent3
    []parent1 //  Breaks -> Slice cannot be embedded anonymously
    pars []parent1 // Correct
    field2 string
}
c := child(parent1{"intefacestring1"}, parent1{"parstring1"}, parent2{"parstring2"}, parent3{"parstring3"}, "childstring")
c.field1 // Incorrect -> ambiguous field
c.parent1.field1 // Correct
c.field3 // Works
c.ii.field1  or c.field1 // No field field1 for c or ii
c.Method1() // Works -> can be used for private vs public classifiers?
```

## HTTP and things

```go
import "net/http"
type myhandler struct {
}
func (h *myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "42")
}
handler := new(myhandler)
http.Handle("/endpoint", handler) // handler type needs to implement http.Handler interface (ServeHTTP method)
http.ListenAndServe(8080, nil) // DefaultServeMux handler -> any other URL 404
//TODO: Routing

// Making external requests
res, err := http.Get(url) // http.Response, Error
defer res.Body.Close()
data, err := ioutil.ReadAll(res.Body) // []byte, Error
numBytes, err := io.Copy(outAsIOWriter, res.Body) // if large response -> buffered can be used to download large files (Streams data 32 kbytes at a time)
res, err := http.Post(url, contenttype, bodyAsIOReader)
// http.Get/Post/Head/PostForm wrapper around (http.Client.Get/*) -> where client is DefaultClient(0 value for http.Client)
```

- See: [gorilla-mux](https://github.com/gorilla/mux)

- Good [Read](https://medium.com/rungo/creating-a-simple-hello-world-http-server-in-go-31c7fd70466e)

### Templates

```go
import "html/template"
Example Template
"""edit.html
<<h1>Editing {{.Title}}</h1>

<form action="/save/{{.Title}}" method="POST">
<div><textarea name="body" rows="20" cols="80">{{printf "%s" .Body}}</textarea></div>
<div><input type="submit" value="Save"></div>
</form>
"""
data = &Data(Title: "This is a title", Body: []byte("This is the body"))
t, _ = template.ParseFiles("edit.html")
t.Execute(w, data)
//Read All Templates at load
template.Must(template.ParseFiles("view.html", "edit.html"))
templates.ExecuteTemplate(w, "edit.html", data)
```

- [Doc Page](https://golang.org/doc/articles/wiki/final.go)

## Defer/Panic/Recover

```go
// Arguments to deferred functions are evaluated at the time of defer statement evaluation
i := 0
defer fmt.Prinln(i) //Output: 0
i++
return

//LIFO
for i := 1; i <= 3; i++ {
    defer fmt.Println(i) //Output: 321
}
func f() int {
    def func() {return 5}()
    return 1
} // returns 1: Return value of deferrred functions is discarded

// Deferred functions may read and assign to the returning function's named return values.
func f() (i int) {
    defer func() {i++} ()
    return 1
}
// panic rescuing
func f() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("panic handled", r, string(debug.Stack()))
        }
    }()
    panic("panic attack!!!")
}
```

- After panic is called, normal execution stops and already deferred functions are executed.

## JSON

```go
import "encoding/json"
type response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"` //JSON tags See package reflect for StructTags
    Withouttag string
}
mapD := &response2{
    Page:   1,
    Fruits: []string{"apple", "peach", "pear"}
    Withouttag: "unTagged"
}
mapB, _ := json.Marshal(mapD)// keys are always strings (int converted to strings)
//mapB, _ := json.MarshalIndent(mapD, "", "  ") pretty-print
fmt.Printf("%T\n%v\n%#v\n", mapB, mapB, mapB) // []uint8, [123, 34], []byte{values}
fmt.Println(string(mapB))// {"page": 1, "fruits": [values], "Withouttag": "unTagged"}
var decodedJSON response2
json.Unmarshal(mapB, &decodedJSON)
fmt.Println(reflect.DeepEqual(decodedJSON, *mapD)) // true

// Arbitrary JSON
b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
var f interface{}
err := json.Unmarshal(b, &f) // map[string]interface{}{}
m := f.(map[string]interface{})
for k, v := range m {switch vv := v.(type){}} // type switch on json default types
// {"Name": "Wednesday, "Age", 6 float64, "Parents"[]interface{}{"Gomez",...}}

Encoding/Decoding Streams

func NewDecoder(r io.Reader) *Decoder
func NewEncoder(w io.Writer) *Encoder
var m map[string]interface{}
json.NewDecoder(os.Stdin/http.(Request/Response).Body/files).Decode(&m)

```

See [go-blog](https://blog.golang.org/json)

## Exporting identifiers

```go
package main

import (
	"fmt"
	"play.ground/animals"
)

func main() {
	dog := animals.Dog{
		BarkStrength: 10,
	}
	dog.Age = 1
	dog.SetName("dogname")
	fmt.Printf("Counter: %#v\n%#v", dog, dog.GetName())
}
-- go.mod --
module play.ground
-- animals/animals.go --
package animals

// animal represents information about all animals.
type animal struct {
	name string // not exported to outer packages but can be embedded to types in the same package
	Age  int // exported to outer packages if type can be exported
}

// Dog represents information about dogs.
type Dog struct {
	animal // this type and its fields will not be exported to outer packages but embedded to Dog
	BarkStrength int
}

func (d *Dog) SetName(name string) {
	d.name = name
}
func (d *Dog) GetName() string {
	return d.name
}

```

## Error Handling

```go
// Errors are values
{val, err := Work()} * 100
//Instead
type errWorker {
    w Worker
    err Error
}
for errWorker.w.Work() {
}
if errWorker.err != nil {
}
// Similar to what bufio.Writer.Write + bufio.Writer.Flush and bufio.Scanner.Scan() + Scanner.Err()

RandomConstantError := errors.New("random error")
if err != RandomConstantError {
}
type MyError struct {
    d Data
    Err error // wrapping errors in errors Eg: os.PathError
}
func (e *MyError) Error() string {return errorAsString}
var myerror error = &MyError{RandomConstantError}
if err, ok := myerror.(*MyError); ok {
    if err.Err == RandomConstantError {
    }
    if err2, ok2 := err.Err.(*MyError2); ok {
    }
    fmt.Println("this is my error")
}
//Go 1.13 stuff
func (e *MyError) Unwrap() error { return e.Err}

if errors.Is(myerror, RandomConstantError) {
    //unwraps until inside error matches
}
var errorinside *MyError
if errors.As(myerror, &errorinside){
    //unwraps until type matches and assigns errorinside
    errorinside.Error() == myerror.Err.Error() // true
}

wrapped := fmt.Errorf("wrapped error: %w", errorinside)
```

- [Intro](https://blog.golang.org/errors-are-values) Blog Post
- See [fmt.Errorf](https://golang.org/src/fmt/errors.go?s=624:674#L7)
- https://blog.golang.org/go1.13-errors

## Concurrency

### Context:

`context` package that makes it easy to pass request-scoped values, cancellation signals, and deadlines across API boundaries to all the goroutines involved in handling a request.

```go
type Context interface {
    Done() <- chan struct{} //receive only
    Err() error
    Deadline() (deadle time.Time, ok bool)
    Value(key interface{}) interface{}
}
func Background/WithCancel/WithDeadline(context, args ...argType) ... //deriving contexts from existing ones.
type CancelFunc func()

```

- Done channel -> cancellation signal to all functions running "on behalf of the context".
- `WithCancel/WithTimeout` contexts can be cancelled before the parent context.
- Should be used in all incomingoutgoing requests -> gives easy control over cancellations/timeouts. See [http.Client.Do](https://golang.org/pkg/net/http/#Client.Do)
- [Blog](https://blog.golang.org/context)

### Select Statement

```go
select {
case xc <- x:
    //sent x on xc
case y := <-yc:
    //received y from yc
default:
    fmt.Println("No communication could happen")
case <-quit:
    teardown()
    quit <- "Safe to exit"
    return
}
//timeout slow requests
select {
case val := <- c:
case <-time.After(1 * time.Second):
    // timed out
}
//Multiplex using select
func Multiplex(in1, in2 <- chan string) {
    c := make(chan string)
    go func() {
        for {
            select {
            case s:= in1:
                c <- s
            case s:= <-in2:
                c <- s
            }
        }
    }
    return c
}
func Multiplex(inputs ... <- chan string) <- chan string{
    cases := make([]reflect.SelectCase, len(inputs))
    ch := make(chan string)
    for i := range inputs {
        cases[i] = reflect.SelectCase({
            Dir: reflect.SelectRecv,
            Chan: reflect.ValueOf(inputs[i])
        })
    }
    go func() {
        for {
            chosenidx, valAsreflectDotValue, ok := reflect.Select(cases)
            ch <- valAsreflectDotValue.String()
        }
    }()
    return ch
}

```

- Evaluates all and invokes the first one ready to proceed, pseudo-random if multiple are ready. `default` if none are ready.

### Channels

```go
ch := make(chan int)
<-ch // receive
ch <- 1 //send
//Generator returns a read-only channel
func Generator() <- chan string {
    c := make(chan string)
    go func() {
        c <- 1
    }()
    return c
}
c1 := Generator(); c2 := Generator()
for {
    <-c1; <-c2;
    // these are now lockstepped (even if generation on one is quicker both will be executed in order one by one)
}
// Multiplex channels to avoid lockstep for each channel
func Multiplex(inputs ... <- chan string) <- chan string {
    c := make(chan string)
    for i := range inputs {
        ch := inputs[i]
        go func() {for { c <- <-ch}}()
    }
    return c
}
ch := Multiplex(Generator(), Generator())
// Restore sequence
// Add wait channel in the message that is listened to in the generator and sent values to by the consumer
// Essentially, this stops further generation until the slowest generator is ready
type Message struct{
    msg string
    wait chan bool
}
Generator() <-chan Message{
    waitForIt := make(chan bool)
    go func() {
        c <- Message{1, waitForIt}
        doSomeTimeTakingStuff()
        <-waitForIt
    }
}
ch := Multiplex(Generator(), Generator())
for {
    m1 := <-ch;
    m2 := <-ch // * number of generators;
    m1.wait <- true
    m2.wait <- true
}
```

- Sending/receiving is blocking. Channels communicate and synchronize with the same operation. Exception: **Buffered channels**.
- TODO: Buffered channels
- Generators: functions that returns channels
- Multiplexing -> FanIn outputs of multiple channels using into one channel using goroutines to avoid locksteps. 
- [Restoring](https://github.com/golang/talks/blob/master/content/2012/concurrency/support/sequenceboring.go) sequence
- [Slides](https://talks.golang.org/2012/concurrency.slide#1) -> [Code](https://github.com/golang/talks/tree/master/content/2012/concurrency/support)

