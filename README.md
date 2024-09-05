# go-metron

An API client allowing Go programs to interact with Metron.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          

## Usage

```go
import "github.com/jyggen/go-metron"
```

Construct a new Metron client.

```go
c, err := metron.NewClient(
	metron.WithAuthentication("username", "password"),
	metron.WithCaching(""),
)
```
