package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

func main() {
	t := time.Now()

	path := flag.String("path", "", "path to dbt project")
	system := flag.String("system", "", "name of the system")

	flag.Parse()

	tables := flag.Args()

	if path == nil {
		panic("need to provide path")
	}
	if system == nil {
		panic("need to provide system")
	}
	if len(tables) == 0 {
		panic("need to provide table names")
	}

	var p MetaParser = &YmlParser{}

	// get all files paths in provided dir
	paths, err := p.GetPaths(*path)
	if err != nil {
		panic(err.Error())
	}

	var wg sync.WaitGroup

	// loop yaml files to find needed tables and create metafiles
	for _, fpath := range paths {
		wg.Add(1)
		go func(path string) {
			err := p.Parse(path, system, &tables)
			if err != nil {
				fmt.Println(err.Error())
			}
			wg.Done()
		}(fpath)
	}

	wg.Wait()

	fmt.Println(time.Since(t).Seconds())

}
