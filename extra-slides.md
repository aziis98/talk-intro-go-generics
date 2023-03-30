
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