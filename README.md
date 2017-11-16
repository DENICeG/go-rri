go-rri
======

Simple lib for implementing an rri client

Example
-------
```
package main

import (
        "github.com/sebidude/go-rri/client"
        "log"
        "os"
        "time"
)

func main() {
        username := os.Getenv("RRIUSERNAME")
        password := os.Getenv("RRIPASSWORD")
        address := os.Getenv("RRIADDRESS")
        client, err := client.NewRriClient(username, password, address)
        if err != nil {
                log.Fatal(err)
        }
        err = client.Login()
        if err != nil {
                log.Println(err)
        }

        order := make(map[string]interface{})
        order["ACTION"] = "CHECK"
        order["VERSION"] = "2.0"
        order["DOMAIN"] = "mydomain.de"

        client.SendOrder(order)

        go client.Read()
        client.Logout()
        time.Sleep(4 * time.Second)
        client.Close()
}

```
