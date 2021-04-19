# 📝 Randtext

Randtext is a Go package to generate random text that reads well. The random generation is based on the input text that the user has previously provided to the generator.

The package offers a `Rand` struct, the source of random text, which exposes two functions: `Feed` and `Generate`. Function `Feed` can be used to feed the random generator with input text.

```go
func Feed(in io.Reader) error
```

Once the random source has been initialized with some text, function `Generate` allows to generate new random text. The caller can specify the desired number of words: the output will contain at most the given number of words.

```go
func Emit(out io.Writer, words int) error
```

Additional text can be fed to the engine even if `Generate` has already been called.

The package also exposes convenience functions `Feed` and `Generate` which operate on a global default source. 

Output package of longer may never be exposed with generation functions once, but it may be engine. 
