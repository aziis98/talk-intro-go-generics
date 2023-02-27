---
marp: true
theme: uncover
size: 4:3
---

<style>
:root {
    font-family: 'Inter', sans-serif;
    font-size: 175%;
}

section.chapter {
    background: #00acd7;
    color: #ecfbff;
}
</style>

<!-- _class: chapter -->

# Introduzione alle Generics in Go

---



## Chi sono?

Antonio De Lucreziis, studente di Matematica e macchinista del PHC

### Cos'è il PHC?

Il PHC è un gruppo di studenti di Matematica con interessi per, open source, Linux, self-hosting e soprattutto smanettare sia con hardware e software (veniteci pure a trovare!)

---

_The Go 1.18 release adds support for generics. Generics are the biggest change we’ve made to Go since the first open source release_

Fonte: https://go.dev/blog/intro-generics

---

Ad esempio una delle grandi mancanze della _stdlib_ del Go è stata l'assenza delle funzioni `Min(x, y)` e `Max(x, y)`. 

Non c'è mai stati modo di definirle in modo generico.


```go
func MinInt(x, y int) int {
    if x < y {
        return x
    }
    return y
}
```

```go
func MinFloat32(x, y float32) float32 {
    if x < y {
        return x
    }
    return y
}
```

---

Invece ora con le nuove generics possiamo scrivere

```go
import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](x, y T) T {
    if x < y { 
        return x
    }
    return y
}
```

---

Cos'è `constraints.Ordered`?

```go
type Ordered interface {
    Integer | Float | ~string
}

type Float interface {
    ~float32 | ~float64
}

type Integer interface {
    Signed | Unsigned
}

type Signed interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
```

---

# Prova

Lorem ipsum dolor sit amet consectetur adipisicing elit. Perferendis, voluptas. Doloribus ea consectetur fugiat quaerat eum magni eos earum placeat dolorum. Nesciunt nostrum tenetur magnam facere magni sapiente illo pariatur?



