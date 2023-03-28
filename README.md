
# Introduzione alle Generics in Go - DevFest GDG

Repo con tutti gli esempi e le slides della presentazione.

**Descrizione.** In questo talk introdurremo le generics del Go 1.18 e vedremo alcuni _pattern_ ed _anti-pattern_ del loro utilizzo. 

- [Scarica il PDF con le slides](https://github.com/aziis98/talk-intro-go-generics/raw/build/slides.pdf)

&nbsp;

<div align="center">
<img src="/assets/devfest-logo.png" height="100" />
&nbsp; &nbsp;
<img src="/assets/logo-circuit-board.svg" height="100" />
</div>

&nbsp;

## Setup

These slides are made using _Marp_

```bash
$ npm install
```

## Usage

To preview and build the slides use

```bash
# Show slides preview
$ npm run preview

# Build slides
$ npm run build:html
$ npm run build:pdf
```

## Go

There is a Makefile with various utilities for running, build and decompiling the Go examples. 

```bash
# Show usage
$ make

# Run/build/decomp examples
$ make run-<subproject> 
$ make compile-<subproject> 
$ make compile-noinline-<subproject> 
$ make decomp-<subproject> 
$ make decomp-noinline-<subproject> 
```

---

# Extra: Per chi non sa cosa sono le interfacce leggere prima questo

```go
type Circle struct {
    Radius float64
}
type Rectangle struct {
    Width, Height float64
}

type Shape interface {
    Area() float64
    Perimeter() float64
}

func (c Circle) Area() float64 {
    return c.Radius * c.Radius * math.Pi
}
func (c Circle) Perimeter() float64 {
    return 2 * c.Radius * math.Pi
}
func (c Circle) Curvature() float64 {
    return 1 / c.Radius
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}
func (r Rectangle) Perimeter() float64 {
    return 2 * r.Radius * math.Pi
}
func (r Rectangle) CornerCount() int {
    if (r.Width == 0 && r.Height == 0) { return 0 }
    if (r.Width == 0 || r.Height == 0) { return 2 }
    return 4
}

func AreaOverPerimeter(s Shape) float64 {
    return s.Area() / s.Perimeter()
}

// c1 := Circle{ Radius: 5.0 }
c1 := Circle{ 5.0 }

// r1 := Rectangle{ Width: 2.0, Height: 3.0 }
r1 := Rectangle{ 2.0, 3.0 }

AreaOverPerimeter(c1)
AreaOverPerimeter(r1)
```

