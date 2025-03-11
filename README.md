# Futaba SBUS serial protocol implementation on GO

## Installing
``` bash
go get github.com/romario5/go-serial-sbus
```


## Usage:
```go
import (
    "fmt"
    serial "github.com/romario5/go-serial-common"
    sbus "github.com/romario5/go-serial-sbus"
)

func main() {
    // Imagine we got some serial reader.
    rwc := GetReadWriteCloser()
    sbus := &sbus.SBUS{}
    packet := &serial.ChannelsPacket{}
    packet.Channels[0] = 1200
    err := sbus.WritePacket(rwc, packet)
    if err != nil {
        fmt.Println("Error on writing channels", err)
    }
}
```