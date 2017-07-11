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
  args := daap.Args{
    Machine: daap.NewEnvMachine(),
    Env:     []string{
      "NODE_ENV=production",
    },
    Mounts:  []daap.Mount{
      daap.Volume("/tmp/dev/data", "/var/data"),
    },
  }

  process := daap.NewProcess(image, args)

  if err := process.Run(ctx); err != nil {
    panic(err)
  }

  b, _ := ioutil.ReadAll(process.Stdout)
  fmt.Println(string(b))
  // Output of "your/image" here
}
```
