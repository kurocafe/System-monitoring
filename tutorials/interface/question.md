# Q8. Interface Basics

Define an interface:
```go
type Shape interface {
    Area() float64
}
```

Implement `Circle` and `Rectangle` so this code works:
```go
shapes := []Shape{
    Circle{Radius: 5},
    Rectangle{W: 3, H: 4},
}

for _, s := range shapes {
    fmt.Println(s.Area())
}

```