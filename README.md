# DaaP

Docker as a Process

[![Build Status](https://travis-ci.org/otiai10/daap.svg?branch=master)](https://travis-ci.org/otiai10/daap)

# Example

```go
container := NewContainer("debian:latest", Args{
  Machine: &MachineConfig{
    Host:     "tcp://192.168.0.2",
    CertPath: "/Users/otiai10/.docker/machine/machines/foobar",
  },
})

ctx := context.Background()
container.PullImage(ctx)
container.Create(ctx)
container.Start(ctx)
container.Exec(ctx, &Execution{Inline: "uname -a"})
```
