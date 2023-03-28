
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
