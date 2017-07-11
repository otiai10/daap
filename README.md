# DaaP

Docker as a Process

# Test

```sh
% git clone git@github.com:otiai10/daap.git
% cd daap
% ./test.sh
```

# Example

```go
func main() {

  ctx := context.Background()

  image := "your/image"
  args := Args{
    Machine: NewEnvMachine(),
    Env:     []string{
      "NODE_ENV=production",
    },
    Mounts:  []mount.Mount{
      mount.Mount{Type: mount.TypeBind, Source: "/tmp/dev/data", Target: "/var/data"},
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
