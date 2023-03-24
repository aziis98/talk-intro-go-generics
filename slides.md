---
marp: true
theme: uncover
size: 4:3
---

<style>
:root {
    font-family: 'Open Sans', sans-serif;
    font-size: 175%;
    letter-spacing: 1px;

    color: #222;
}

section.chapter {
    background: #00acd7;
    color: #ecfbff;
}

code {
    font-size: 100%;
    line-height: 1.4;

    border-radius: 4px;
}

code.language-go {
    font-family: 'Go Mono', monospace;    
}

section.chapter code {
    background: #00809f;
    color: #ecfbff;   
}

@import "https://unpkg.com/@highlightjs/cdn-assets@11.7.0/styles/github.min.css";
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

&nbsp;

Fonte: https://go.dev/blog/intro-generics

---

## Il Problema

---

<style scoped>
code { font-size: 150% }
</style>

```go
func Min(x, y int) int {
    if x < y {
        return x
    }
    
    return y
}
```

---

```go
func MinInt8(x, y int8) int8 {
    if x < y {
        return x
    }
    
    return y
}

func MinInt16(x, y int16) int8 {
    if x < y {
        return x
    }
    
    return y
}

func MinFloat32(x, y float32) float32 {
    if x < y {
        return x
    }
    
    return y
}
```

---

<style scoped>
code { font-size: 150% }
</style>

```go
...
if x < y {
    return x
}

return y
...
```

---

## La Soluzione

---

#### Type Parameters & Type Sets

```go
import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](x, y T) T {
    if x < y { 
        return x
    }
    return y
}
```

#### Type Inference

```go
var a, b int = 0, 1
Min(a, b)
```

```go
var a, b float32 = 3.14, 2.71 
Min(a, b)
```

---

<style scoped>
code { font-size: 150% }
</style>

```go
func Min[T constraints.Ordered](x, y T) T {
    if x < y { 
        return x
    }

    return y
}
```

---

```go
package constraints

...

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

...
```

<!-- Le interfacce possono introdurre questo type-set per limitare i tipi a cui vengono applicate, l'unico tipo che non possiamo utilizzare nei type-sets sono le interfacce con metodi. -->

---

```go
type Liter float64

type Meter float64

type Kilogram float64
```

---

<style scoped>
code { font-size: 150% }
</style>

## Tipi Generici

```go
type Stack[T interface{}] []T
```

<!-- In realtà non serve usare "interface{}" tutte le volte -->

---

<style scoped>
code { font-size: 150% }
</style>

## Tipi Generici

```go
type Stack[T any] []T
```

---

```go
func (s *Stack[T]) Push(value T) {
    *s = append(*s, value)
}

func (s Stack[T]) Peek() T {
    return s[len(s)-1]
}

func (s Stack[T]) Len() int {
    return len(s)
}
```

---

```go
func (s *Stack[T]) Pop() (T, bool) {
    items := *s

    if len(items) == 0 {
        var zero T
        return zero, false
    }

    newStack, poppedValue := items[:len(items)-1], items[len(items)-1]
    *s = newStack

    return poppedValue, true
}
```

---

```go
func Zero[T any]() T {
    var zero T
    return zero
}
```

---

<!-- _class: chapter -->

# Pattern (1)
Quando usare le generics?

---

### Tipi "Contenitore"

- `[n]T` 
    
    Array di `n` elementi per il tipo `T`

- `[]T` 
    
    Slice per il tipo `T`

- `map[K]V` 
    
    Mappe con chiavi `K` e valori `V`

- `chan T` 
    
    Canali per elementi di tipo `T`

---

In Go sono sempre esistite queste strutture dati "generiche"

Solo che prima delle generics non era possibile definire algoritmi generali per questi tipi di container, ora invece possiamo ed infatti alcune di questi sono "già in prova"

---

<style scoped>
section {
    font-size: 140%;
    line-height: 1.75;
}
</style>

## `golang.org/x/exp/slices`

- `func Index[E comparable](s []E, v E) int`
- `func Equal[E comparable](s1, s2 []E) bool`
- `func Sort[E constraints.Ordered](x []E)`
- `func SortFunc[E any](x []E, less func(a, b E) bool)`
- e molte altre...

---

<style scoped>
section {
    font-size: 140%;
    line-height: 1.75;
}
</style>

## `golang.org/x/exp/maps`

- `func Keys[M ~map[K]V, K comparable, V any](m M) []K`
- `func Values[M ~map[K]V, K comparable, V any](m M) []V`
- e molte altre...

---

<style scoped>
section {
    font-size: 140%;
    line-height: 1.75;
}
</style>

## Strutture Dati Generiche

Esempio notevole: <https://github.com/zyedidia/generic> (1K:star: su GitHub)
- `mapset.Set[T comparable]`, set basato su un dizionario.
- `multimap.MultiMap[K, V]`, dizionario con anche più di un valore per chiave.
- `stack.Stack[T]`, slice ma con un'interfaccia più simpatica rispetto al modo idiomatico del Go.
- `cache.Cache[K comparable, V any]`, dizionario basato su `map[K]V` con una taglia massima e rimuove gli elementi usando la strategia LRU.
- `bimap.Bimap[K, V comparable]`, dizionario bi-direzionale.
- `hashmap.Map[K, V any]`, implementazione alternativa di `map[K]V` con supporto per _copy-on-write_.
- e molte altre...

---

<!-- _class: chapter -->

# Anti-Pattern (1)
Generics vs Interfacce

---

## Momento Quiz

```go
func WriteOneByte(w io.Writer, data byte) {
    w.Write([]byte{data})
}

...

d := &bytes.Buffer{}
WriteOneByte(d, 42)
```

```go
func WriteOneByte[T io.Writer](w T, data byte) {
    w.Write([]byte{data})
}

...

d := &bytes.Buffer{}
WriteOneByte(d, 42)
```

---

```
BenchmarkInterface
BenchmarkInterface-4            135735110            9.017 ns/op

BenchmarkGeneric
BenchmarkGeneric-4              50947912            22.26 ns/op
```

---

```go
//go:noinline
func WriteOneByte(w io.Writer, data byte) {
    w.Write([]byte{data})
}

...

d := &bytes.Buffer{}
WriteOneByte(d, 42)
```

---

```
BenchmarkInterface
BenchmarkInterface-4            135735110            9.017 ns/op

BenchmarkInterfaceNoInline
BenchmarkInterfaceNoInline-4    46183813            23.64 ns/op

BenchmarkGeneric
BenchmarkGeneric-4              50947912            22.26 ns/op
```

---

```go
d := &bytes.Buffer{} /* (*bytes.Buffer) */

WriteOneByte(d /* (io.Writer) */, 42)
```

↓

```go
d := &bytes.Buffer{} /* (*bytes.Buffer) */

(io.Writer).Write(d /* (io.Writer) */, []byte{ 42 })
```

↓

```go
d := &bytes.Buffer{} /* (*bytes.Buffer) */

(*bytes.Buffer).Write(d /* (*bytes.Buffer) */, []byte{ 42 })
```

---

<!-- _class: chapter -->

# Anti-Pattern (2)
Utility HTTP

---

```go
type Validator interface {
    Validate() error
}

func DecodeAndValidateJSON[T Validator](r *http.Request) (T, error) {
    var value T
    err := json.NewDecoder(r.Body).Decode(&value)
    if err != nil {
        var zero T
        return zero, err
    }

    if err := value.Validate(); err != nil {
        var zero T
        return zero, err
    }

    return value, nil
}
```

---

```go
type FooRequest struct {
    A int    `json:"a"`
    B string `json:"b"`
}

func (foo FooRequest) Validate() error {
    if foo.A < 0 {
        return fmt.Errorf(`parameter "a" cannot be lesser than zero`)
    }
    if !strings.HasPrefix(foo.B, "baz-") {
        return fmt.Errorf(`parameter "b" has wrong prefix`)
    }

    return nil
}

...

foo, err := DecodeAndValidateJSON[FooRequest](r)
if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
```

---

```go
func DecodeAndValidateJSON(r *http.Request, target *Validator) error {
    err := json.NewDecoder(r.Body).Decode(target)
    if err != nil {
        return err
    }

    if err := (*target).Validate(); err != nil {
        return err
    }

    return nil
}

...

var foo Validator = FooRequest{}
if err := DecodeAndValidateJSON(r, &foo); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
```

In realtà anche in questo caso non serviva introdurre necessariamente delle generics

---

<style scoped>
code { font-size: 150% }
</style>

L'unico problema è che siamo obbligati a fare questo cast che non è molto estetico

```go
var foo Validator = FooRequest{}
```

---

<!-- _class: chapter -->

# Pattern (2)
Vediamo un analogo di `PhantomData<T>` dal Rust per rendere _type-safe_ l'interfaccia di una libreria

---

In Rust è obbligatorio utilizzare tutte le generics usate. 

```rust
// Ok
struct Foo<T> { a: String, value: T }

// Errore
struct Foo<T> { a: String }
```

In certi casi però vogliamo introdurle solo per rendere _type-safe_ un'interfaccia o per lavorare con le _lifetime_.

```go
// Ok
use std::marker::PhantomData;

struct Foo<T> { a: String, foo_type: PhantomData<T> }
```

---

Proviamo ad usare questa tecnica per rendere _type-safe_ l'interfaccia con `*sql.DB`

```go
package database

type IdTable interface {
    PrimaryKey() *string
    Columns()    []any
}

type Ref[T IdTable] string

type Table[T IdTable] struct {
    Name     string
    PkColumn string
}
```

---

```go
// a random db library
// type DB = *sql.DB

func Create[T IdTable](d DB, t Table[T], row T) (Ref[T], error) 

func Insert[T IdTable](d DB, t Table[T], row T) (Ref[T], error)

func Read[T IdTable](d DB, t Table[T], ref Ref[T]) (*T, error) 

func Update[T IdTable](d DB, t Table[T], row T) error 

func Delete[T IdTable](d DB, t Table[T], id string) error 
```

---

```go
func Read[T IdTable](d DB, t Table[T], ref Ref[T]) (*T, error) {
    result := d.QueryRow(
        fmt.Sprintf(
            `SELECT * FROM %s WHERE %s = ?`, 
            t.Name, t.PkColumn,
        ), 
        string(ref),
    )

    var value T
    if err := result.Scan(value.Columns()...); err != nil {
        return nil, err
    }

    return &value, nil
}
```

---

```go
type User struct {
    Username  string
    FullName  string
    Age       int
}

func (u *User) PrimaryKey() *string {
    return &u.Username
}

func (u User) Columns() []any {
    return []any{ &u.Username, &u.FullName, &u.Age }
}

var UsersTable = Table[User]{
    Name: "users",
    PkColumn: "username",
}
```

---

```go
db := ...

user1 := &User{ "aziis98", "Antonio De Lucreziis", 24 }
ref1, _ := database.Insert(db, UsersTable, user1)

...

user1, _ := database.Read(db, UsersTable, ref1) 
```

---

<!-- _class: chapter -->

# Altro esempio caotico
Vediamo come implementare le promise in Go con le generics

---

```go
type Promise[T any] struct {
    value T
    err   error
    done  <-chan struct{}
}

func (p Promise[T]) Await() (T, error) {
    <-p.done
    return p.value, p.err
}

type Waiter { Wait() error }

func (p Promise[T]) Wait() error {
    <-p.done
    return p.err
}
```

---

```go
type PromiseFunc[T any] func(resolve func(T), reject func(error))

func Start[T any](f PromiseFunc[T]) *Promise[T] {
    done := make(chan struct{})
    p := Promise{ done: done }

    f(
        func(value T) { p.value = value; done <- struct{} },
        func(err error) { p.err = err; done <- struct{} }
    )

    return &p
}
```

---

```go
func AwaitAll(ws ...Waiter) error {
    var wg sync.WaitGroup
    wg.Add(len(ws))
    
    errChan := make(chan error, len(ws))
    for _, w := range ws {
        go func(w Waiter) {
            defer wg.Done()
            err := w.Wait()
            if err != nil {
                errChan <- err
            }
        }(w)
    }
    ...
```

---

```go
    ...
    done := make(chan struct{})
    go func() {
        defer close(done)
        wg.Wait()
    }()

    select {
    case err := <-errChan:
        return err
    case <-done:
        return nil
    }
}
```

---

# Fine :C
