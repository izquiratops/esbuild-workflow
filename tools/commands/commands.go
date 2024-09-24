package commands

import (
	"fmt"
	"log"

	"github.com/evanw/esbuild/pkg/api"

	cssPlugin "github.com/izquiratops/dobunezumi/tools/plugins/css"
	htmlPlugin "github.com/izquiratops/dobunezumi/tools/plugins/html"
	httpPlugin "github.com/izquiratops/dobunezumi/tools/plugins/http"
	"github.com/izquiratops/dobunezumi/tools/utils/directory"
)

func Build(entryFilePath, distLocalPath string, enableMinify bool) {
	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{entryFilePath},
		Bundle:            true,
		Metafile:          true,
		MinifyWhitespace:  enableMinify,
		MinifyIdentifiers: enableMinify,
		MinifySyntax:      enableMinify,
		Outdir:            distLocalPath,
		Write:             true,
		Plugins: []api.Plugin{
			httpPlugin.Plugin(),
			htmlPlugin.Plugin(distLocalPath),
			cssPlugin.Plugin(),
		},
	})

	if len(result.Errors) > 0 {
		log.Fatal(result.Errors)
	}

	fmt.Printf("build successfully done\n%s", result.Metafile)
}

func Serve(entryFilePath, distLocalPath string, enableMinify bool) error {
	buildOptions := api.BuildOptions{
		EntryPoints:       []string{entryFilePath},
		Bundle:            true,
		MinifyWhitespace:  enableMinify,
		MinifyIdentifiers: enableMinify,
		MinifySyntax:      enableMinify,
		Outdir:            distLocalPath,
		Write:             true,
		Plugins: []api.Plugin{
			httpPlugin.Plugin(),
			htmlPlugin.Plugin(distLocalPath),
			cssPlugin.Plugin(),
		},
	}

	// Create a context for the build
	ctx, err := api.Context(buildOptions)
	if err != nil {
		log.Fatal("failed to create context:", err)
	}

	defer ctx.Dispose()

	// Watch for changes and rebuild
	watchErr := ctx.Watch(api.WatchOptions{})
	if watchErr != nil {
		log.Fatal("failed to start watch mode:", watchErr)
	}

	// Start the server
	result, serveErr := ctx.Serve(api.ServeOptions{
		Servedir: distLocalPath,
		Port:     8080,
	})
	if serveErr != nil {
		log.Fatal("failed to start server:", serveErr)
	}

	fmt.Printf("server started on http://localhost:%d\n", result.Port)

	// Returning from main() exits immediately in Go.
	// Block forever so we keep watching and don't exit.
	<-make(chan struct{})

	fmt.Printf("server stopped\n")

	return nil
}

func Clean(distLocalPath string) {
	directory.Clean(distLocalPath)

	fmt.Println("clean completed successfully.")
}