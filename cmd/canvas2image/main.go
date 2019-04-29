package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/atomotic/iiifutils"
)

func main() {

	manifest := flag.String("manifest", "", "IIIF Manifest URL")
	canvas := flag.String("canvas", "", "IIIF Canvas URL")

	flag.Parse()
	if *manifest == "" || *canvas == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	image, err := iiifutils.ImageFromCanvas(*manifest, *canvas)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s/full/500,/0/default.jpg\n", image)
	}

}
