# go-metron

An API client allowing Go programs to interact with Metron.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          

## Usage

Construct a new Metron client.

```go
import "github.com/jyggen/go-metron"

c, err := metron.NewClient(
	metron.WithAuthentication("username", "password"),
	metron.WithCaching(""),
)
```

### Caching

### Rate Limiting