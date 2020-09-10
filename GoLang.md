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

