# 📝 Randtext

Randtext is a Go package to generate random text that reads well. The random generation is based on the input text that the user provides to the engine.

The package exposes two functions: `Feed` and `Emit`. Function `Feed` can be used to feed the random generator with input text.

```go
func Feed(in io.Reader) error
```

Once some text has been fed to the engine, function `Emit` allows to generate random text. The caller can specify the desired number of words. The output will never be longer, but it may be shorter.

```go
func Emit(out io.Writer, words int) error
```

Additional text can be fed to the engine even if `Emit` has already been called.
