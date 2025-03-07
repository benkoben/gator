package main

import (
	"fmt"
	"log"

	"github.com/benkoben/gator/internal/config"
)

func main(){
    cfg, err := config.Read()
    if err != nil {
        log.Fatalf("could not read config file: %s", err)
    }

    fmt.Println(cfg)

    cfg.SetUser("benko")

    updatedCfg, err := config.Read()
    if err != nil {
        log.Fatalf("could not read config file: %s", err)
    }
    fmt.Println(updatedCfg)

}
