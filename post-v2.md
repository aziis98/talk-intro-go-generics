
# Introduzione alle Generics in Go

Dalla versione 1.18 del Go è stata aggiunta la possibilità di definire funzioni e strutture parametrizzate da tipi con i cosiddetti _type parameters_ o anche dette semplicemente _generics_. Lo scopo principale è che ci permettono di scrivere codice indipendente dai tipi specifici che utilizzano.

Più precisamente le tre novità relative alle _generics_ sono

- Sia funzioni che tipi possono essere parametrizzati rispetto a dei tipi (_type parameters_)

- In un modo ristretto le interfacce possono essere utilizzare per definire "insiemi di tipi" (_type sets_)

- Un minimo di _type inference_ che ci permette di omettere i _type parameters_ quando si riescono a dedurre dal contesto.

## Il problema

Uno degli esempi più lampanti della necessità di aggiungere le _generics_ al Go è che ad esempio manca la funzione `Min` per interi nella libreria standard del linguaggio e bisogna scriversi ogni volta un'implementazione speciale di `Min(x, y)` per il tipo numerico che vogliamo utilizzare (al momento c'è solo `math.Min(float64, float64) float64` che però necessita di conversioni se la vogliamo usare per interi o anche solo `float32`)

```go
func Min(x, y int) int {
    if x < y {
        return x
    }
    
    return y
}
```
e quindi siamo costretti o a ricopiare tante volte la stessa funzione specializzandola a mano per i vari tipi per cui la vogliamo usare

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

o ad esempio scrivere una funzione che prende input `interface{}` (ora c'è un alias ad `any`) ed utilizzare uno switch sul tipo per decidere cosa fare.

Alternativamente ci sono anche tecniche per risolvere questo problema che utilizzano `go generate`.

## Type Parameters

Dalla versione 1.18 come già detto sono stati introdotti i _type parameters_ che possiamo ad esempio applicare ad una funzione come segue per introdurre una _generics_ alla funzione.

```go
import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](x, y T) T {
    if x < y { 
        return x
    }
    return y
}
```

Questa funzione generica può essere utilizzata come segue

```go
var a, b int = 0, 1
Min[int](a, b)
...
var a, b float32 = 3.14, 2.71 
Min[float32](a, b)
```

o anche omettendo la chi

```go
var a, b int = 0, 1
Min(a, b)
...
var a, b float32 = 3.14, 2.71 
Min(a, b)
```

Più precisamente la sintassi per introdurre un _type parameter_ è la seguente, tra quadre indichiamo il nome del tipo che vogliamo introdurre seguito da un'interfaccia che indica il vincolo che il parametro deve rispettare.

```
[T Vincolo1, R interface{ Method(), ... }, ...]
```

Vediamo meglio come definire questi vincoli. 

## Type Sets

Fino ad ora quando in Go si parla di interfacce si prende in considerazione il _method set_ per quell'interfaccia e dato un tipo esso implementerà l'interfaccia se ha tutti i metodi del _method set_.

<img src="../assets/method-sets.png" />

Un modo duale di vedere la cosa è di pensare al _type set_ generato da un'interfaccia ovvero l'insieme di tutti i tipi che rispettano un'interfaccia come nel diagramma che segue. 

<img src="../assets/type-sets.png" />

Seguendo questa linea di pensiero è stata estesa la sintassi delle interfacce per ammettere una dichiarazione esplicita del _type set_ sotto forma di un'unione di tipi che l'interfaccia accetta (se omessa allora accetterà tutti i tipi).

<img src="../assets/type-sets-2.png" />

Sono state inoltre aggiunte delle semplificazioni nella sintassi per cui quando scriviamo il vincolo di un tipo e vogliamo usare un'interfaccia senza metodi possiamo scrivere direttamente l'unione di tipi.

- Se vogliamo scrivere `[T interface{ int | float32 }]` può anche essere scritto come `[T int | float32]`

- Inoltre è stato aggiunto (finalmente) un alias per il tipo `interface{}` detto `any`

### Tipi con la tilde

Consideriamo ad esempio il seguente frammento di codice.

```go
func SumTwoIntegers[T int](x, y int) T {
    if x < y { return x }
    return y
}
```

Però se abbiamo un tipo "sinonimo" di `int` ma non suo alias allora non potremo usarlo nella chiamata generica perché di per sé `Liter` e `int` non sono compatibili.

```go
type Liter int
...

var a, b int = 1, 2
SumTwoIntegers(a, b) // Ok

var a, b Liter = 1, 2
SumTwoIntegers(a, b) // Errore
```

È stata però inserita la seguente sintassi con la `~` prima di un tipo per intendere anche tutti i suoi _type alias_.

```go
func SumTwoIntegers[T ~int](x, y int) T {
    if x < y { return x }
    return y
}
```

ed infatti poi il seguente frammento compila

```go
type Liter int
...

var a, b int = 1, 2
SumTwoIntegers(a, b) // Ok

var a, b Liter = 1, 2
SumTwoIntegers(a, b) // Ok
```

Utilizzando la stessa tecnica possiamo considerare un caso più utile definito direttamente la modulo `constraints`.

```go
package constraints

...

type Float interface {
    ~float32 | ~float64
}

...
```

e sulla stessa riga possiamo finalmente vedere com'è definita l'interfaccia `constraints.Ordered` vista in precedenza. 

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

## Tipi Generici

Possiamo ad esempio definire uno _stack_ generico come segue

```go
type Stack[T any] []T
```

Questo ci permette di vedere come possono essere definiti metodi per tipi generici 

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

Il metodo _pop_ è un po' più interessante, decidiamo che ritornerà il valore tolto dallo stack seguito da `true` se lo stack non era vuoto e "`0.(T)`" e `false` se lo stack era vuoto.

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

Per ora per ottenere un valore rappresentante il valore di default per un tipo serve introdurre una nuova variabile e poi ritornarla. 

Alternativamente possiamo definire questa funzione di _utility_

```go
func zero[T any]() T {
    var zero T
    return zero
}
```

e diciamo che questo è un primo _pattern_ che spesso può essere utile quando lavoriamo con le generics.

In realtà [anche il resto del mondo](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#the-zero-value) si è accorto di questo trick e già si sta pensando a delle soluzioni come utilizzare `nil` o `_` per indicare un valore di default per un tipo.

## Pattern: Tipi Contenitore

Vediamo ora qualche caso utile in cui utilizzare le _generics_.

### Tipi generici nativi

Fin dall'inizio il Go ha avuto alcune strutture generiche _backed in_

- `[n]T` 
    
    Array di `n` elementi per il tipo `T`

- `[]T` 
    
    Slice per il tipo `T`

- `map[K]V` 
    
    Mappe con chiavi `K` e valori `V`

- `chan T` 
    
    Canali per elementi di tipo `T`

solo che non essendoci le generics non era possibile definire algoritmi generici che le utilizzassero. 

Ora finalmente è possibile ed infatti già ci sono moduli sperimentali della libreria standard che introducono una manciata di funzioni utili per lavorare con queste strutture

- `golang.org/x/exp/slices`

    - `func Index[E comparable](s []E, v E) int`
    
    - `func Equal[E comparable](s1, s2 []E) bool`
    
    - `func Sort[E constraints.Ordered](x []E)`
    
    - `func SortFunc[E any](x []E, less func(a, b E) bool)`
    
    - e molte altre...

- `golang.org/x/exp/maps`

    - `func Keys[M ~map[K]V, K comparable, V any](m M) []K`
    
    - `func Values[M ~map[K]V, K comparable, V any](m M) []V`
    
    - e molte altre...

### Strutture Dati Generiche

Stanno anche nascendo alcuni moduli esterni con varie strutture dati generiche come ad esempio <https://github.com/zyedidia/generic> (~1K stelle su GitHub) che fornisce le seguenti strutture

- `mapset.Set[T comparable]`, set basato su un dizionario.

- `multimap.MultiMap[K, V]`, dizionario con anche più di un valore per chiave.

- `stack.Stack[T]`, slice ma con un'interfaccia più simpatica rispetto al modo idiomatico del Go.

- `cache.Cache[K comparable, V any]`, dizionario basato su `map[K]V` con una taglia massima e rimuove gli elementi usando la strategia LRU.

- `bimap.Bimap[K, V comparable]`, dizionario bi-direzionale.

- `hashmap.Map[K, V any]`, implementazione alternativa di `map[K]V` con supporto per _copy-on-write_.

- e molte altre...

## Anti-Pattern (1)

Vediamo ora un esempio (a posteriori discutibile) di come creare una utility http.

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

# Pattern: "PhantomData"
Vediamo un analogo di `PhantomData<T>` dal Rust per rendere _type-safe_ l'interfaccia di una libreria

---

Proviamo ad usare questa tecnica per rendere _type-safe_ l'interfaccia con `*sql.DB`

```go
type DatabaseRef[T any] string
```

```go
package tables 

// tables metadata
var Users = database.Table[User]{ ... }
var Products = database.Table[Product]{ ... }
```


```go
userRef1 := DatabaseRef[User]("j.smith@example.org") 
...

// Ok
user1, err := database.Read(dbConn, tables.Users, userRef1) 
// Errore
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

...

func Read[T WithPK](d DB, t Table[T], ref Ref[T]) (*T, error) 
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
func AwaitAll[T any](ps ...*Promise[T]) error {
    ...
}
```

```go
type Waiter interface { Wait() error }

func (p Promise[T]) Wait() error {
    <-p.done
    return p.err
}

func AwaitAll(ws ...Waiter) error {
    ...
}
```

---

```go
func ResolveInto[T any](p *Promise[T], target *T) *Promise[T] {
    ...
}
```

```go
AwaitAll(
    ResolveInto(httpRequest1, &result1), // :: *Promise[int]
    ResolveInto(httpRequest2, &result2), // :: *Promise[struct{ ... }] 
    ResolveInto(httpRequest3, &result3), // :: *Promise[any]
    timer1,                              // :: *Promise[struct{}]
)
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

