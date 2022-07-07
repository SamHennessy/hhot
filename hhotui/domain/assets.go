package domain

import (
	"log"
	"os/exec"
)

func AssetBuild() {
	// c := exec.Command("npx", "tailwindcss", "--input=css/app.css", "--output=./output/app.css", "--postcss")
	c := exec.Command("node", "build")
	c.Dir = "./assets"

	out, err := c.CombinedOutput()

	if err != nil {
		log.Println("Build CSS Error: ", err)
	}

	log.Println("Build CSS: ", string(out))
}
