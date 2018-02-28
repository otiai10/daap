# daap

[![Build Status](https://travis-ci.org/otiai10/daap.svg?branch=master)](https://travis-ci.org/otiai10/daap) [![codecov](https://codecov.io/gh/otiai10/daap/branch/master/graph/badge.svg)](https://codecov.io/gh/otiai10/daap)

Use Docker container as completely separated subprocess.

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

stream, _ := container.Exec(ctx, &Execution{Inline: "uname -a"})
for payload := range stream {
    fmt.Println(payload.Text())
}
```
