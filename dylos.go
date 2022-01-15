package main

import (
        "bytes"
        "fmt"
        "io"
        "strings"
)
import "log"
import "github.com/tarm/serial"

func main() {
        config := &serial.Config{
                Name: "/dev/tty.usbserial-2140",
                Baud: 9600,
                ReadTimeout: 1,
                Size: 8,
        }

        data := make(chan string)
        stream, err := serial.OpenPort(config)
        if err != nil {
                log.Fatal(err)
        }

        small := make(chan string)
        large := make(chan string)

        go func() {
          for d := range small {
                 fmt.Println("Small: ", strings.TrimSpace(d))
          }
        }()

        go func() {
                for d := range large {
                        fmt.Println("Large: ", strings.TrimSpace(d))
                }
        }()

        go func() {
                temp := ""
               for d := range data {
                       if d != "" {
                               if d == "," {
                                       small <- strings.TrimSpace(temp)
                                       temp = ""
                               }
                               if d == "\n" {
                                       large <- strings.TrimSpace(temp)
                                       temp = ""
                               }
                               if d != "," {
                                       temp += d
                               }
                       }
               }
        }()


        for {
                buf := make([]byte, 8)
                _, err := stream.Read(buf)
                if err != nil && err != io.EOF {
                        log.Println(err)
                }
                buf = bytes.Trim(buf, "\x00")
                data <- string(buf)
        }

}

