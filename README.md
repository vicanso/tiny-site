# forest
web framework by iris

## static

Create static asset's packr

```bash
packr
```

## test

```bash
GO_ENV=test VIPER_INIT_TEST=true go test -race -coverprofile=test.out ./... && go tool cover --html=test.out
```