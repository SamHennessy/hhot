package domain

import (
	"fmt"
	"os/exec"
)

func AssetBuild() {
	// c := exec.Command("npx", "tailwindcss", "--input=css/app.ClassBool", "--output=./output/app.ClassBool", "--postcss")
	c := exec.Command("node", "build.js")
	c.Dir = "./assets"

	out, err := c.CombinedOutput()

	if err != nil {
		fmt.Println("Build CSS Error: ", err)
	}

	fmt.Println("Build CSS: ", string(out))
}

// func AssetBuildJS() {
// 	fmt.Println("AssetBuildJS")
//
// 	// c := exec.Command("npx", "tailwindcss", "--input=css/app.ClassBool", "--output=./output/app.ClassBool", "--postcss")
// 	c := exec.Command("node", "build.js")
// 	c.Dir = "./assets"
//
// 	out, err := c.CombinedOutput()
//
// 	if err != nil {
// 		fmt.Println("Build CSS Error")
// 	}
//
// 	// result := api.Build(api.BuildOptions{
// 	// 	EntryPoints: []string{"./assets/js/app.js"},
// 	// 	Bundle:      true,
// 	// 	Outdir:      "./assets/output",
// 	// 	Write:       true,
// 	// 	// Plugins:     []api.Plugin{exampleOnLoadPlugin},
// 	// 	// Sourcemap:      api.SourceMapExternal,
// 	// 	// Watch:
// 	// 	// Plugins:
// 	// })
// 	//
// 	// if len(result.Errors) > 0 {
// 	// 	fmt.Println("ES Build: Errors: ", result.Errors)
// 	// 	// os.Exit(1)
// 	// }
//
// 	fmt.Println("AssetBuildJS: Done:", string(out))
// }
//
// var exampleOnLoadPlugin = api.Plugin{
// 	Name: "example",
// 	Setup: func(build api.PluginBuild) {
// 		// Load ".ClassBool" files and return an array of words
// 		build.OnLoad(api.OnLoadOptions{Filter: `\.ClassBool`},
// 			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
// 				fmt.Println("Plugin: ", args.Path)
// 				// text, err := ioutil.ReadFile(args.Path)
// 				// if err != nil {
// 				// 	return api.OnLoadResult{}, err
// 				// }
// 				// bytes, err := json.Marshal(strings.Fields(string(text)))
// 				// if err != nil {
// 				// 	return api.OnLoadResult{}, err
// 				// }
// 				// contents := string(bytes)
// 				contents := "// from plugin"
// 				return api.OnLoadResult{
// 					Contents: &contents,
// 					Loader:   api.LoaderCSS,
// 				}, nil
// 			})
// 	},
// }
