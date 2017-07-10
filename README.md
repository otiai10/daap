# DaaP

Docker as a Process

```go
func main() {

  ctx := context.Background()

  image := "otiai10/foo"
  args := Args{
    Machine: NewEnvMachine(),
    Env:     []string{
      "NODE_ENV=production",
    },
    Mounts:  []mount.Mount{
      mount.Mount{Type: mount.TypeBind, Source: "/Users/otiai10/data/dev", Target: "/var/data"},
    },
  }

  proc := NewProcess(image, args)

  process := daap.NewProcess(image, args)

  if err := process.Run(ctx); err != nil {
    panic(err)
  }

  b, _ := ioutil.ReadAll(process.Stdout)
  fmt.Println(string(b))
}
```
