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
<img src="../assets/devfest-logo.png" height="100" />
<img src="../assets/logo-circuit-board.svg" height="100" />
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

## Soluzioni Pre-Generics

- Fare una funzione che prende `any` e mettere tanti switch

- Utilizzare `go generate` [...]

- Copia incollare tante volte la funzione per ogni tipo

---

## Soluzione Post-Generics

#### Type Parameters

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
...
var a, b float32 = 3.14, 2.71 
Min[float32](a, b)
```

---

#### Type Inference

```go
var a, b int = 0, 1
Min(a, b)
...
var a, b float32 = 3.14, 2.71 
Min(a, b)
```

---

<style scoped>
code { font-size: 150% }
</style>

```
[T Vincolo1, R interface{ Method(), ... }, ...]
```

---

<style scoped>section { justify-content: space-between; }</style>

## Type Sets

<img src="../assets/method-sets.png" />

&nbsp;

---

<style scoped>section { justify-content: space-between; }</style>

## Type Sets

<img src="../assets/type-sets.png" />

&nbsp;


---

<style scoped>section { justify-content: space-between; }</style>

## Type Sets

<img src="../assets/type-sets-2.png" />

&nbsp;

---

#### Type Sets (Sintassi)

```go
[T interface{}] ~> [T any]

[T interface{ int | float32 }] ~> [T int | float32]
```

---

#### Type Sets

```go
func SumTwoIntegers[T int](x, y int) T {
    if x < y { return x }
    return y
}
```

```go
type Liter int
```

```go
var a, b int = 1, 2
SumTwoIntegers(a, b) // Ok

var a, b Liter = 1, 2
SumTwoIntegers(a, b) // Errore
```

---

#### Type Sets

```go
func SumTwoIntegers[T ~int](x, y int) T {
    if x < y { return x }
    return y
}
```

```go
type Liter int
```

```go
var a, b int = 1, 2
SumTwoIntegers(a, b) // Ok

var a, b Liter = 1, 2
SumTwoIntegers(a, b) // Ok
```


---

#### Type Sets

```go
package constraints

...

type Float interface {
    ~float32 | ~float64
}

...
```

---

#### Type Sets

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

---

<!-- _class: chapter -->

# Tipi Generici

---

<style scoped>
code { font-size: 120% }
</style>

```go
type Stack[T any] []T
```

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

Per ora ci tocca utilizzare questa funzione di _utility_

```go
func Zero[T any]() T {
    var zero T
    return zero
}
```

&nbsp;

<https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#the-zero-value>

---

<!-- _class: chapter -->

# Pattern: Tipi Contenitore

---

### Tipi generici nativi

- `[n]T` 
    
    Array di `n` elementi per il tipo `T`

- `[]T` 
    
    Slice per il tipo `T`

- `map[K]V` 
    
    Mappe con chiavi `K` e valori `V`

- `chan T` 
    
    Canali per elementi di tipo `T`

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

# Anti-Pattern 1
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
func DecodeAndValidateJSON(r *http.Request, target Validator) error {
    err := json.NewDecoder(r.Body).Decode(target)
    if err != nil {
        return err
    }

    if err := target.Validate(); err != nil {
        return err
    }

    return nil
}

...

var foo FooRequest
if err := DecodeAndValidateJSON(r, &foo); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
```

In realtà anche in questo caso non serviva introdurre necessariamente delle generics

---

Quindi nella maggior parte dei casi se ci ritroviamo a scrivere una funzione generica con un **parametro vincolato ad un'interfaccia** forse dobbiamo porci qualche domanda

---

<!-- _class: chapter -->

# Anti-Pattern 2
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

#### Go 1.18 Implementation of Generics via Dictionaries and Gcshape Stenciling

-  _A **gcshape** (or gcshape grouping) is a collection of types that all **share the same instantiation of a generic function/method**_.

- _Two concrete types are in the same gcshape grouping if and only if they have the **same underlying type** or they are **both pointer types**._

- _To avoid creating a different function instantiation for each generic call with distinct type arguments (which would be pure stenciling), we **pass a dictionary along with every call**_.

:link: [generics-implementation-dictionaries-go1.18.md](https://github.com/golang/proposal/blob/master/design/generics-implementation-dictionaries-go1.18.md)

<!-- :link: [Go 1.18 implementation of generics via dictionaries and gcshape stenciling](https://github.com/golang/proposal/blob/master/design/generics-implementation-dictionaries-go1.18.md) -->

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
type DatabaseRef[T any] string
```

```go
package tables 

// tables metadata
var Users = Table[User]{ ... }
var Products = Table[Product]{ ... }
```


```go
userRef1 := DatabaseRef[User]("j.smith@example.org") 

...
// Ok
user1, err := database.Read(dbConn, tables.Users, userRef1) 

// Error
user2, err := database.Read(dbConn, tables.Products, userRef1)
```

---

```go
package database

type WithPK interface {
    PrimaryKey() *string
}

type Ref[T WithPK] string

type Table[T WithPK] struct {
    Name     string
    PkColumn string
    Columns  func(*T) []any
}
```

---

```go
package database

func Create[T WithPK](d DB, t Table[T], row T) (Ref[T], error) 

func Insert[T WithPK](d DB, t Table[T], row T) (Ref[T], error)

func Read[T WithPK](d DB, t Table[T], ref Ref[T]) (*T, error) 

func Update[T WithPK](d DB, t Table[T], row T) error 

func Delete[T WithPK](d DB, t Table[T], id Ref[T]) error 
```

---

```go
func Read[T WithPK](d DB, t Table[T], ref Ref[T]) (*T, error) {
    result := d.QueryRow(
        fmt.Sprintf(
            `SELECT * FROM %s WHERE %s = ?`, 
            t.Name, t.PkColumn,
        ), 
        string(ref),
    )

    var value T
    if err := result.Scan(t.Columns(&value)...); err != nil {
        return nil, err
    }

    return &value, nil
}
```

---

```go
package model

type User struct {
    Username  string
    FullName  string
    Age       int
}

func (u *User) PrimaryKey() *string {
    return &u.Username
}
```

```go
package tables

var Users = Table[User]{
    Name: "users",
    PkColumn: "username",
    Columns: func(u *User) []any {
        return []any{ &u.Username, &u.FullName, &u.Age }
    }
}
```

---

```go
user1 := &model.User{ "j.smith@example.org", "John Smith", 36 }

userRef1, _ := database.Insert(db, tables.Users, user1)

...

user1, _ := database.Read(db, tables.Users, userRef1) 
```

---

<!-- _class: chapter -->

# Altro esempio caotico
Vediamo come implementare le _promise_ in Go con le generics

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
func Resolve[T](value T) *Promise[T] {
    return &Promise{ value: value }
}

func Reject[T](err error) *Promise[T] {
    return &Promise{ error: err }
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
type Waiter interface { Wait() error }

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
func ResolveInto[T any](p *Promise[T], target *T) *Promise[T] {
    return Run[T](func(resolve func(T), reject func(error)) {
        value, err := p.Await()
        if err != nil {
            reject(err)
            return
        }

        *target = value
        resolve(value)
    })
}
```

```go
err := AwaitAll(
    ResolveInto(httpRequest1, &result1), // :: *Promise[int]
    ResolveInto(httpRequest2, &result2), // :: *Promise[struct{ ... }] 
    ResolveInto(httpRequest3, &result3), // :: *Promise[any]
    timer1,                              // :: *Promise[struct{}]
)
...
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

// Trick per codificare higher-kinded types
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

// Eq_Refl ovvero l'assioma 
//   forall x : x = x
func Eq_Reflexive[T any]() Eq[T, T] { 
    panic("axiom")
}

// Eq_Symmetric ovvero l'assioma 
//   forall a, b: a = b => b = a
func Eq_Symmetric[A, B any](_ Eq[A, B]) Eq[B, A] { 
    panic("axiom")
}

// Eq_Transitive ovvero l'assioma 
//   forall a, b, c: a = b e b = c => a = c
func Eq_Transitive[A, B, C any](_ Eq[A, B], _ Eq[B, C]) Eq[A, C] { 
    panic("axiom")
}
```

---

## Uguaglianza e Sostituzione

Per ogni funzione `F`, ovvero tipo vincolato all'interfaccia `Term2Term` vorremmo dire che

```
               F
Eq[ A , B ] ------> Eq[ F[A] , F[B] ] 

```

---

## Uguaglianza e Sostituzione

Data una funzione ed una dimostrazione che due cose sono uguali allora possiamo applicare la funzione ed ottenere altre cose uguali

```go
// Function_Eq ovvero l'assioma
//   forall f function, forall a, b term: a = b  => f(a) = f(b)
func Function_Eq[F Term2Term, A, B Term](_ Eq[A, B]) Eq[V[F, A], V[F, B]] {
    panic("axiom")
}
```

---

## Assiomi dell'addizione

```go
type Plus[L, R Term] Term

// "n + 0 = n"

// Plus_Zero ovvero l'assioma 
//   forall n, m: n + succ(m) = succ(n + m)
func Plus_Zero[N Term]() Eq[Plus[N, Zero], N] { 
    panic("axiom")
}

// "n + (m + 1) = (n + m) + 1"

// Plus_Sum ovvero l'assioma 
//   forall a, m: n + succ(m) = succ(n + m)
func Plus_Sum[N, M Term]() Eq[
    Plus[N, V[Succ, M]], 
    V[Succ, Plus[N, M]],
] { panic("axiom") }
```

---

## 1 + 1 = 2

```go
func Theorem_OnePlusOneEqTwo() Eq[Plus[One, One], Two] {
    // 1 + 0 = 1
    var en1 Eq[ Plus[One, Zero], One ] = Plus_Zero[One]()

    // (1 + 0) + 1 = 2
    var en2 Eq[ V[Succ, Plus[One, Zero]], Two ] = Function_Eq[Succ](en1)

    // 1 + 1 = (1 + 0) + 1 
    var en3 Eq[ Plus[One, One], V[Succ, Plus[One, Zero]] ] = Plus_Sum[One, Zero]()
 
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

--- 


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

Impl. ⇝ _Sostituzione testuale post-tokenizzazione_

---

## C++

```cpp
template<typename T>
T min(T const& a, T const& b)
{
    return (a < b) ? a : b;
}
```

Impl. ⇝ _se funziona allora ok_

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

Impl. ⇝ _Monomorfizzazione_

