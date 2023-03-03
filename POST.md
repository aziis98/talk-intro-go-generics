# Introduzione alle Generics in Go

Dalla versione 1.18 del Go è stata aggiunta la possibilità di definire funzioni e strutture parametrizzate da tipi con i cosiddetti _type parameters_ o anche dette semplicemente _generics_. Lo scopo principale è che ci permettono di scrivere codice indipendente dai tipi specifici che utilizzano.

Più precisamente le tre novità relative alle _generics_ sono

- Sia funzioni che tipi possono essere parametrizzati rispetto a dei tipi (_type parameters_)

- In un modo ristretto le interfacce possono essere utilizzare per definire "insiemi di tipi" (_type sets_)

- Un minimo di _type inference_ che ci permette di omettere i _type parameters_ quando si riescono a dedurre dal contesto.

## Il problema

Uno degli esempi più lampanti della necessità di aggiungere le _generics_ al Go è che ad esempio manca la funzione `Min` per interi nella libreria standard del linguaggio e bisogna scriversi ogni volta un'implementazione speciale di `Min(x, y)` per il tipo numerico che vogliamo utilizzare (al momento c'è solo `math.Min(float64, float64) float64` che però necessita di conversioni se la vogliamo usare per interi o anche solo `float32`)

```go
func MinInt(x, y int) int {
    if x < y {
        return x
    }
    return y
}

func MinInt32(x, y int32) int32 {
    if x < y {
        return x
    }
    return y
}

func MinInt64(x, y int64) int64 {
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

Notiamo che l'implementazione è sempre la stessa ma cambia solo la segnatura della funzione. Dal Go 1.18 però possiamo scrivere

```go
import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](x, y T) T {
    if x < y {
        return x
    }
    return y
}
```

Qui la parte nuova da notare è la stringa `[T constraints.Ordered]` che indica che stiamo introducendo un parametro `T` vincolato ad essere ordinabile. 

Questa funzione può essere usata ad esempio con `Min[int64](2, 5)` oppure `Min[float32](2.71, 3.14)`. In particolare dopo aver passato i _type parameters_ possiamo usarla come una qualunque altra funzione, ovvero quanto segue è codice legale

```go
shortMin := Min[int16] // func(int16, int16) int16
```

## Struct generiche

Finalmente ora possiamo anche definire strutture dati generiche come un albero con valori su ogni nodo.

```go
type Tree[T interface{}] struct {
    Left, Right *BinaryTree[T]
    Value       T
}
```

In realtà invece di dover scrivere ogni volta `interface{}` è stato aggiunto l'alias `any` quindi possiamo scrivere direttamente

```go
type Tree[T any] struct {
    Left, Right *BinaryTree[T]
    Value       T
}
```

Vediamo qualche altro esempio, possiamo anche avere un albero con valori solo sulle foglie, in particolare vediamo ora come possiamo anche definire dei metodi su tipi con _type parameters_.

```go
type BinaryTree[T any] interface{
    Has(value T) bool
}

type Leaf[T any] struct {
    value T
}

func (l Leaf[T]) Has(value T) {
    return l.value == value
}

type Branch[T any] struct {
    Left, Right BinaryTree[T]
}

func (b Branch[T]) Has(value T) {
    return b.Left.Has(value) || b.Right.Has(value)
}
```

Giusto per precisare, quando scriviamo `Leaf[T any]` stiamo introducendo un _type parameter_ che possiamo usare a destra nella definizione del tipo, ed anche quando scriviamo un metodo `func (l Leaf[T]) ...` stiamo reintroducendo la variabile T che infatti possiamo usare a destra. Infatti ad esempio

```go
func Zero[T any]() T {
    var zero T
    return zero
}
```

è una funzione ben definita, anzi è anche abbastanza utile quando si vuole ritornare un valore "vuoto" per un certo tipo ma il tipo ci arriva attraverso una _generics_.

## Type sets

Tornando all'esempio di prima della funzione `Min`, abbiamo visto che ogni tipo ha bisogno di un _type constraint_ anche se questo è semplicemente `any`.

Possiamo definire un _type constraint_ utilizzando o una classica interfaccia del Go oppure usando un _type set_ come nel caso di `Min`, vediamo com'è definito in particolare `Ordered` nella libreria standard (per la precisione per ora è nel pacchetto `golang.org/x/exp` che contiene il codice ancora considerato sperimentale)


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

Qui definiamo delle interfacce con delle union di vari tipi, al momento ci sono alcune restrizioni riguardo quali tipi sono ammessi per farne l'unione

- Tutti i tipi primitivi come `string`, `int`, `rune`, ...

- Un tipo primitivo preceduto da tilde semplicemente intende che rilassiamo il vincolo a tutti i tipi alias a quel tipo, quindi ad esempio `float64` non ammetterebbe 

    ```go
    type Liter float64
    ```

    mentre `~float64` sì in quanto il tipo `Liter` è solo un alias per `float64`.

- Altre interfacce rappresentanti un "constraint" ed in particolare solo interfacce senza metodi (per ora è stata aggiunta questa restrizione per evitare alcuni problemi teorici di complessità  dell'inferenza dei constraint)

## Type inference

La _type inference_ usa un algoritmo relativamente semplice unidirezionale, se ad esempio abbiamo una chiamata ad una funzione generica, se riusciamo a dedurre i _type parameters_ solo dagli argomenti della funzione allora possiamo omettere i _type parameters_ nella chiamata della funzione. Ad esempio questo non si può applicare alla funzione `Zero` di prima


```go
func Zero[T any]() T {
    var zero T
    return zero
}
```

In questo caso serve esplicitare il _type parameter_ anche se in realtà si potrebbe dedurre dal contesto (in questo casa stiamo già specificando il tipo del risultato ma questo poterebbe l'algoritmo di inferenza ad essere bidirezionale che complicherebbe molto le cose)

```go
var x int = Zero[int]()
```

## Quando usare le generics?

### Container

In Go in realtà esistono già da sempre alcune strutture dati "generiche" ovvero

- `[n]T` 
    
    Array di `n` elementi per il tipo `T`

- `[]T` 
    
    Slice per il tipo `T`

- `map[K]V` 
    
    Mappe con chiavi `K` e valori `V`

- `chan T` 
    
    Canali per elementi di tipo `T`

Prima delle generics c'è sempre stato il problema che non era possibile definire algoritmi generici per questi tipi di container, ora invece possiamo ed infatti alcune di questi sono "già in prova"

- Il modulo `golang.org/x/exp/slices` già offre

    - `func Index[E comparable](s []E, v E) int`

    - `func Equal[E comparable](s1, s2 []E) bool`

    - `func Sort[E constraints.Ordered](x []E)`

    - `func SortFunc[E any](x []E, less func(a, b E) bool)`

    - ...

- Invece il modulo `golang.org/x/exp/maps` ad esempio ha

    - `func Keys[M ~map[K]V, K comparable, V any](m M) []K`

    - `func Values[M ~map[K]V, K comparable, V any](m M) []V`

    - ...

Però ora che abbiamo le generics possiamo noi stessi definire tipi _container_ generici, ad esempio c'è già il modulo <https://github.com/zyedidia/generic> (1k stelle su GitHub) con varie strutture dati generiche come 

- `mapset` 

- `multimap` 

- `stack` (con un'interfaccia più simpatica rispetto al modo idiomatico del Go)

- `queue`

- `cache` (basato su `map[K]V` con una taglia massima che rimuove gli elementi usando la strategia LRU)

- `bimap`

- `hashmap` (implementazione alternativa di `map[K]V` con supporto per _copy-on-write_)


### Metodi generici

Consideriamo la seguente struttura dati generica

```go
package option

type Option[T any] struct{
    present bool
    value   T
}

func Some[T any](value T) Option[T] {
    return Option{ true, value }
}

func None[T any]() Option[T] {
    return Option{ present: false }
}

func (o Option[T]) Map(f func(T) T) Option[T] {
    if !o.present {
        return o
    }

    return Some(f(o.value))
}
```

```go
double := func (v int) int { 
    return v * 2
}

o1 := option.Some[int](10)
o2 := o1.Map(double) // Option[int]{ present: true, value: 20 }

o3 := option.None[int]()
o4 := o3.Map(double) // Option[int]{ present: false }
```

Questo sembrerebbe un buon utilizzo delle generics per introdurre il tipo `Option[T]` già molto usato in molti linguaggi funzionali e non. Ad esempio Rust che ha deciso di integrarli direttamente nel linguaggio prima con la macro `try!` e poi con l'operatore `?`.

Al momento però non è possibile introdurre generics nelle funzioni quindi la seguente funzione sarebbe illegale

```go
func (Option[T]) MapToOther[S any](f func(T) S) Option[s] {
    if !o.present {
        return o
    }

    return Some(f(o.value))
}
```

<!-- questo già richiede creare specializzazioni per `Option[T]` per ogni utilizzo di `T`, e se aggiungiamo `S` e ciò complicherebbe abbastanza il compilatore de Go che è noto per essere molto veloce rispetto ad altri per numero di righe di codice al secondo (in particolare il Go ha un compilatore single-pass che inizia a sputare codice macchina quando già sta leggendo i file di codice in input) [FACT CHECK]. -->

Il compilatore del Go per compilare codice con delle generics utilizza una tecnica chiamata **monomorfizzazione** ovvero per ogni utilizzo di una funzione generica, vede quali sono i _type parameter_ utilizzati e specializza quella funzione o tipo al caso particolare. 

Nel caso di strutture dati ad esempio

```go
package option

// se da qualche parte utilizziamo "Option[int]" allora viene generata questa struttura e queste funzioni.
type Option_int struct{
    present bool
    value   int
}

func Some_int(value int) Option_int {
    return Option{ true, value }
}

func None_int() Option_int {
    return Option{ present: false }
}

// se da qualche parte utilizziamo "Option[string]" allora viene generata questa struttura e queste funzioni.
type Option_string struct{
    present bool
    value   string
}

func Some_string(value string) Option_string {
    return Option{ true, value }
}

func None_string() Option_string {
    return Option{ present: false }
}
```

questo ha il vantaggio di creare specializzazioni per ogni caso specifico senza fare uso di puntatori (quindi il codice rimane abbastanza performante) però al prezzo di grandezza del binario generato.

(<https://go.googlesource.com/proposal/+/master/design/43651-type-parameters.md#no-parameterized-methods>)

In Go i metodi sui tipi sono stati introdotti come modo di astrazione via le interfacce. Detto in altri termini dato un tipo ed un'interfaccia possiamo facilmente vedere se questo verifica l'interfaccia e se abbiamo un valore di tipo quell'interfaccia dovremmo poter passare il tipo a quell'interfaccia anche da un altro modulo. Con questo principio in mente dovrebbe poter essere possibile definire la seguente interfaccia

```go
type Processor interface {
    Process[T any](v T) T
}
```

solo che a questo punto una funzione (non di per sé generica) potrebbe fare le chiamate

```go
func Example(v Processor) {
    fmt.Println(v.Process[int](3))
    fmt.Println(v.Process[string]("5"))
    fmt.Println(v.Process[[]string]([]string{"a", "b"}))
}
```

ad esempio su un tipo come

```go
type Foo struct{}

func (Foo) Process[T any](value T) T {
    return value
}
```

questo porterebbe a vari problemi sul come generare il codice per questo tipo in quanto fino a _runtime_ non sarebbero note quali chiamate generiche vengono instanziate a meno di non fare dell'analisi statica molto elaborata.

Un modo potrebbe essere fare come Rust e non permettere definire interfacce/trait con metodi/funzioni che introducono nuovi tipi parametrici, vedremo in Go 2...

## Quando non usare le generics?

Ci potrebbe venire in mente di scrivere una funzione per leggere tutto da un `io.Reader` aggiungendo il vincolo `io.Reader` al _type parameter_

```go
func ReadSome[T io.Reader](r T) ([]byte, error)
```

in questo caso però potevamo già scrivere

```go
func ReadSome(r io.Reader) ([]byte, error)
```

senza dover usare generics in quanto già le interfacce ci permettono di scrivere codice generico. Inoltre al momento l'implementazione delle generics monomorfizza solo per non-"_pointer types_" e genera un'unica implementazione se il _receiver_ è un _pointer type_ in quanto già il puntatore contiene dati sul tipo passato in input e non aggiungere specializzazioni non migliorerebbe la performance.

### Se le implementazioni differiscono

Se vogliamo scrive del codice generico per più tipo chi possiamo chiedere se l'implementazione è la stessa per tutti i tipi o se differisce, la regolare generale è che se l'implementazione è la stessa allora va bene usare un _type parameter_ altrimenti, se le implementazioni differiscono, come detto prima conviene usare semplicemente delle interfacce.

### Reflection

Ci sono alcuni casi in cui vorremmo implementare una qualche operazione per tipi che non possono avere metodi (ad esempio per dei tipi primitivi) e tale operazione è diversa per ogni tipo. In questo caso conviene usare la _reflection_ invece di interfacce o generics. Ad esempio `encoding/json` funziona in questo modo.

## Utilizzi interessanti delle generics

Le generics possono essere utilizzate anche solo per rendere il codice più sicuro dal punto di vista dai tipi (e per fare meno conversioni a _runtime_), ad esempio quando definiamo una struct generica nessuno ci obbliga ad utilizzare effettivamente il _type parameter_ che introduciamo.

```go
type DatabaseRef[T any] struct{ Id string }
    
type DatabaseTable[T any] struct {
	Table    string
	IdKey    string
	GetIdPtr func(entry T) *string
}

func (t DatabaseTable[T]) RefForId(id string) DatabaseRef[T] {
	return DatabaseRef[T]{id}
}

func (t DatabaseTable[T]) RefForValue(v T) DatabaseRef[T] {
	return t.RefForId(*t.GetIdPtr(v))
}

func DatabaseRead[T any](db Database, table DatabaseTable[T], ref DatabaseRef[T]) (*T, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE %s = ?`, table.Table, table.IdKey)

	result := db.Get(query, ref.Id)

	var value T
	if err := result.Scan(&value); err != nil {
		return nil, err
	}

	return &value, nil
}
```

più eventualmente anche altre funzioni per creare, modificare ed eliminare i dati nel database

```go
// create un'entrata nel database (generando un nuovo id)
func DatabaseCreate[T any](db Database, table DatabaseTable[T], value *T) (DatabaseRef[T], error)

// create un'entrata nel database inserendo il valore fornito
func DatabaseWrite[T any](db Database, table DatabaseTable[T], value *T) (DatabaseRef[T], error)

// aggiorna un'entrata nel db usando la ref ed il valore passato
func DatabaseUpdate[T any](db Database, table DatabaseTable[T], DatabaseRef[T], value *T) error

// elimina un'entrata nel db usando la ref fornita 
func DatabaseDelete[T any](db Database, table DatabaseTable[T], DatabaseRef[T]) error
```

Questo ci permette di creare delle _reference_ tipate che permettono di rendere _type safe_ l'API per interagire con il nostro database

```go
type User struct {
	Username  string
	FirstName string
	LastName  string
}

// rappresenta una tabella tipata con una primary key specifica
var UsersTable = DatabaseTable[User]{
	Table: "users",
	IdKey: "username",
	GetIdPtr: func(u User) *string {
		return &u.Username
	},
}
```

che potremo utilizzare ad esempio

```go
// user1 :: *User
user1 := &User{"j.smith", "John", "Smith"}

// ref1 :: DatabaseRef[User]
ref1,  _ := DatabaseWrite(db, UsersTable, u1)

// user2 :: *User
user2, _ := DatabaseRead(db, UsersTable, ref1)
```

Inoltre potenzialmente potremmo creare queste `DatabaseRef[T]` solo se l'entrata corrispondente esiste nel database, in questo modo renderemmo l'API di questa libreria completamente _type safe_ e rendendola anche privata potremmo anche rendere tali istanze costruibili solo interagendo con il database.

### Algebraic Data Types

In molti linguaggi funzionali sono presenti i cosiddetti _algebraic data types_ che permettono di codificare molte strutture dati in modo molto semplice e generico. Inoltre se il linguaggio è tipato permettono anche di creare delle interfacce di libreria molto pulite e sicure come ad esempio nel caso di _option_.

#### Option[T]

Vediamo meglio questo esempio, consideriamo la definizione `type Option[T] = Some(T) | None` (pseudo-codice? forse è valido in Scala boh) può essere codificato in Go come segue

```go
type private struct{}

type Option[T any] interface {
	isOption(private)
    Match(
		caseSome func(value T),
		caseNone func(),
	)
}

type some[T any] struct{ Value T }
func (some[T]) isOption(private) {}

func Some[T any](value T) Option[T] {
    return some{ value }
}

func (v some[T]) Match(caseSome func(value T), caseNone func()) {
    caseSome(v.Value)
}

type none[T any] struct{}
func (none[T]) isOption(private) {}

func None[T any]() Option[T] {
    return none{}
}

func (v none[T]) Match(caseSome func(value T), caseNone func()) {
	caseNone()
}
```

Possiamo così costruire istanze di `Option[T]` a piacimento usando le funzioni pubbliche `Some[T](value)` e `None[T]()` ma in quanto l'interfaccia ha un metodo che usa un tipo privato non possiamo definire tipi al di fuori che la rispettano. 

E se ci viene passata un'istanza di `Option` l'unica cosa che possiamo fare è chiamare la funzione match e valutare i vari casi.

#### Either[A, B]

o anche `type Either[A, B] = Left A | Right B`

```go
type Either[A, B any] interface {
	Match(
		caseLeft func(value A),
		caseRight func(value B),
	)
}

type Left[A, B any] struct{ Value A }

func (v Left[A, B]) Match(
	caseLeft func(value A),
	caseRight func(value B),
) {
	caseLeft(v.Value)
}

type Right[A, B any] struct{ Value B }

func (v Right[A, B]) Match(
	caseLeft func(value A),
	caseRight func(value B),
) {
	caseRight(v.Value)
}
```

## Bibliografia

Materiale per il talk sulle generics del Go:

- <https://go.dev/blog/intro-generics> &mdash; Basi sulle generics

- <https://go.dev/blog/when-generics> &mdash; Quanto usarle e quando no

- <https://go.googlesource.com/proposal/+/HEAD/design/43651-type-parameters.md> &mdash; Questa è proprio la proposal ufficiale in teoria
