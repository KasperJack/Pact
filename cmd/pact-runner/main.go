package main

import (
    "fmt"
    "os"

)


func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: pact-runner <package.lua>")
        os.Exit(1)
    }

    path := os.Args[1]

    if err := runner.RunFile(path); err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)	 
    }
    fmt.Printf("%d", "hello")
}
