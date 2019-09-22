package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"gopkg.in/mgo.v2"

	"github.com/go-numb/go-bitflyer/auth"
	"github.com/go-numb/go-bitflyer/v1"
	"github.com/go-numb/go-bitflyer/v1/public/executions"
	"github.com/go-numb/go-bitflyer/v1/types"
)

func main() {
	done := make(chan struct{})

	go getExec()

	<-done
}

func getExec() {
	bf := v1.NewClient(&v1.ClientOpts{
		&auth.AuthConfig{
			"", "",
		},
	})

	sec, err := mgo.Dial("http://localhost:28015")
	if err != nil {
		log.Fatal("mongo database can not set")
	}
	defer sec.Close()

	col := sec.DB("bffx").C("executions")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	id := 0

	rid, ex := getexec(id, bf)
	if ex == nil {
		log.Fatal("gets executions data is nil")
	}
	id = rid

	for _, e := range ex {
		col.Insert(e)
	}

	for {
		select {
		case <-ticker.C:
			rid, ex := getexec(id, bf)
			if ex == nil {
				log.Error("gets executions data is nil")
				continue
			}
			id = rid

			for _, e := range ex {
				col.Insert(e)
			}
		}
	}
}

func getexec(id int, bf *v1.Client) (int, []executions.Execution) {
	page := types.Pagination{
		Count: 499,
	}
	if id != 0 {
		page.After = id
	}

	exec, _, err := bf.Executions(&executions.Request{
		ProductCode: types.ProductCode("FX_BTC_JPY"),
	})
	if err != nil {
		return id, nil
	}

	var ex []executions.Execution
	ex = []executions.Execution(*exec)

	return 0, ex
}
