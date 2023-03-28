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

&nbsp;

<div style="display: flex; align-items: center; justify-content: center; gap: 2rem;">
<img src="devfest-logo.png" height="100" />
<img src="logo-circuit-board.svg" height="100" />
</div>

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

func MinInt16(x, y int16) int16 {
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

```go
var a, b int = 0, 1
Min[int](a, b)
```

```go
var a, b float32 = 3.14, 2.71 
Min[float32](a, b)
```

---

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
type Liter float64

type Meter float64

type Kilogram float64
```

```go
func Min[T float64](x, y T) T {
    if x < y { 
        return x
    }
    return y
}
```

```go
var a, b Liter = 1, 2
Min(a, b) // Errore
```

---

```go
type Liter float64

type Meter float64

type Kilogram float64
```

```go
func Min[T ~float64](x, y T) T {
    if x < y { 
        return x
    }
    return y
}
```

```go
var a, b Liter = 1, 2
Min(a, b) // Ok
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
WriteOneByte[*bytes.Buffer](d, 42)
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

# Confronto con altri linguaggi
Vediamo come funzionano le generics in Go confrontandole con altri linguaggi

---

## C

```c
// versione semplificata da linux/minmax.h
#define min(x, y) ({                \ // block expr (GCC extension)
    typeof(x) _min1 = (x);          \ // eval x once
    typeof(y) _min2 = (y);          \ // eval y once
    (void) (&_min1 == &_min2);      \ // check same type 
    _min1 < _min2 ? _min1 : _min2;  \ // do comparison
}) 
```

---

## C++

```cpp
template<typename T>
T min(T const& a, T const& b)
{
    return (a < b) ? a : b;
}
```

---

## Rust

```rust
pub fn min<T: PartialOrd>(a: T, b: T) -> T {
    if a < b {
        a
    } else {
        b
    }
}
```

---

## Go _Gcshape Stenciling_

-  _A **gcshape** (or gcshape grouping) is a collection of types that can all **share the same instantiation of a generic function/method** in our implementation when specified as one of the type arguments_.

- _Two concrete types are in the same gcshape grouping if and only if they have the **same underlying type** or they are **both pointer types**._

- _In order to avoid creating a different function instantiation for each invocation of a generic function/method with distinct type arguments (which would be pure stenciling), we **pass a dictionary along with every call** to a generic function/method_.

<!-- :link: [Go 1.18 implementation of generics via dictionaries and gcshape stenciling](https://github.com/golang/proposal/blob/master/design/generics-implementation-dictionaries-go1.18.md) -->

---

<!-- _class: chapter -->

# Anti-Pattern (2)
Utility HTTP

---

```go
// library code
type Validator interface {
    Validate() error
}

func DecodeAndValidateJSON[T Validator](r *http.Request) (T, error) {
    var value T
    if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
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
// client code
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
```

```go
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
```

---

```go
type PromiseFunc[T any] func(resolve func(T), reject func(error))

func Run[T any](f PromiseFunc[T]) *Promise[T] {
    done := make(chan struct{})
    p := Promise{ done: done }

    go f(
        func(value T) { p.value = value; done <- struct{} },
        func(err error) { p.err = err; done <- struct{} }
    )

    return &p
}
```

---

```go
type Waiter { Wait() error }

func (p Promise[T]) Wait() error {
    <-p.done
    return p.err
}

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

```go
Validate()
```

---

<!-- _class: chapter -->

# 1 + 1 = 2
_Proof checking_ in Go 

---

## Premesse

```go
type Bool interface{ isBool() }

type Term interface{ isTerm() }

type Term2Term interface{ isTerm2Term() }

// trick to encode higher-kinded types
type V[H Term2Term, T Term] Term
```

---

## Assiomi dei Naturali 

```go
type Zero Term
type Succ Term2Term

// Alcuni alias utili
type One = V[Succ, Zero]
type Two = V[Succ, V[Succ, Zero]]
type Three = V[Succ, V[Succ, V[Succ, Zero]]]
```

---

## Uguaglianza

```go
type Eq[A, B any] Bool

// Eq_Refl ~> forall x : x = x
func Eq_Reflexive[T any]() Eq[T, T] { 
    panic("axiom")
}

// Eq_Symmetric ~> forall a, b: a = b => b = a
func Eq_Symmetric[A, B any](_ Eq[A, B]) Eq[B, A] { 
    panic("axiom")
}

// Eq_Transitive ~> forall a, b, c: a = b e b = c => a = c
func Eq_Transitive[A, B, C any](_ Eq[A, B], _ Eq[B, C]) Eq[A, C] { 
    panic("axiom")
}
```

---

## "Funzionalità dell'uguale"

Per ogni funzione `F`, ovvero tipo vincolato all'interfaccia `Term2Term` vorremmo dire che

```
               F
Eq[ A , B ] ------> Eq[ F[A] , F[B] ] 

```

---

## "Funzionalità dell'uguale"

Data una funzione ed una dimostrazione che due cose sono uguali allora possiamo applicare la funzione ed ottenere altre cose uguali

```go
// Function_Eq ~> forall f function, forall a, b term:
//     a = b  => f(a) = f(b)
func Function_Eq[F Term2Term, A, B Term](_ Eq[A, B]) Eq[V[F, A], V[F, B]] {
    panic("axiom")
}
```

---

## Assiomi dell'addizione

```go
type Plus[L, R Term] Term

// "n + 0 = n"

// Plus_Zero ~> forall n, m: n + succ(m) = succ(n + m)
func Plus_Zero[N Term]() Eq[Plus[N, Zero], N] { 
    panic("axiom")
}

// "n + (m + 1) = (n + m) + 1"

// Plus_Sum ~> forall a, m: n + succ(m) = succ(n + m)
func Plus_Sum[N, M Term]() Eq[
    Plus[N, V[Succ, M]], 
    V[Succ, Plus[N, M]],
] { panic("axiom") }
```

---

## 1 + 1 = 2

```go
func Theorem_OnePlusOneEqTwo() Eq[Plus[One, One], Two] {
    var en1 Eq[ Plus[One, Zero], One ] = Plus_Zero[One]()

    var en2 Eq[
        V[Succ, Plus[One, Zero]],
        Two
    ] = Function_Eq[Succ](en1)

    var en3 Eq[
        Plus[One, One],
        V[Succ, Plus[One, Zero]],
    ] = Plus_Sum[One, Zero]()

    return Eq_Transitive(en3, en2)
}
```

---

## 1 + 1 = 2

```go
func Theorem_OnePlusOneEqTwo() Eq[Plus[One, One], Two] {
    return Eq_Transitive(
        Plus_Sum[One, Zero](), 
        Function_Eq[Succ](
            Plus_Zero[One](),
        ),
    )
}
```

---

<!-- _class: chapter -->

# Conclusione

---

<style scoped>
section {
    text-align: left;
}
</style>

### Regole generali

Per scrivere _codice generico_ in Go

- Se l'implementazione dell'operazione che vogliamo supportare non dipende del tipo usato allora conviene usare dei **type-parameter**

- Se invece dipende dal tipo usato allora è meglio usare delle **interfacce**

- Se invece dipende sia dal tipo e deve anche funzionare per tipi che non supportano metodi (ad esempio per i tipi primitivi) allora conviene usare **reflection** 

---

# Fine :C

_Domande_

---

<style scoped>
li {
    font-size: 80%;
}
</style>

## Bibliografia

- <https://go.dev/blog/intro-generics>

- <https://go.dev/blog/when-generics>

- <https://github.com/golang/proposal/blob/master/design/generics-implementation-dictionaries-go1.18.md>
