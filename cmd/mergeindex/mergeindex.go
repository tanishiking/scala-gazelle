package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/stackb/scala-gazelle/api/jarindex"
	"github.com/stackb/scala-gazelle/pkg/mergeindex"
)

const debug = false

var (
	outputFile       string
	predefinedLabels string
	preferredLabels  string
)

func main() {
	if debug {
		// index.ListFiles(".")
		log.Println("args:", os.Args)
	}

	log.SetPrefix("mergeindex: ")
	log.SetFlags(0) // don't print timestamps

	args := os.Args[1:]
	if len(args) == 1 && strings.HasPrefix(args[0], "@") {
		paramsFile := args[0][1:]
		var err error
		args, err = readParamsFile(paramsFile)
		if err != nil {
			log.Fatalln("failed to read params file:", paramsFile, err)
		}
	}
	files, err := parseFlags(args)
	if err != nil {
		log.Fatal(err)
	}

	if err := merge(files...); err != nil {
		log.Fatal(err)
	}
}

func readParamsFile(filename string) ([]string, error) {
	params := []string{}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		params = append(params, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return params, nil
}

func parseFlags(args []string) (files []string, err error) {
	fs := flag.NewFlagSet("mergeindex", flag.ExitOnError) // flag.ContinueOnError
	fs.StringVar(&predefinedLabels, "predefined", "", "a comma-separated list of labels to be considered predefined")
	fs.StringVar(&preferredLabels, "preferred", "", "a comma-separated list of labels to be considered preferred")
	fs.StringVar(&outputFile, "output_file", "", "the output file to write")
	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: mergeindex @PARAMS_FILE | mergeindex OPTIONS FILES")
		fs.PrintDefaults()
	}
	if err = fs.Parse(args); err != nil {
		return
	}

	if outputFile == "" {
		log.Fatal("-output_file is required")
	}

	// files = []string{}
	files = fs.Args()
	if len(files) == 0 {
		err = fmt.Errorf("positional args should be a non-empty list of .jarindex.json files to merge")
	}

	log.Println("flags predefinedLabels:", predefinedLabels)

	return
}

func merge(filenames ...string) error {
	jars := make([]*jarindex.JarFile, len(filenames))
	for i, filename := range filenames {
		jar, err := mergeindex.ReadJarFileProtoFile(filename)
		if err != nil {
			return err
		}
		jars[i] = jar
	}

	// spec is the final object to write as output
	// var spec index.IndexSpec
	predefined := strings.Split(predefinedLabels, ",")
	preferred := strings.Split(preferredLabels, ",")

	index, err := mergeindex.MergeJarFiles(func(format string, args ...interface{}) {
		log.Printf("warning: "+format, args...)
	}, predefined, jars)
	if err != nil {
		log.Fatal(err)
	}

	index.Preferred = preferred

	if err := mergeindex.WriteJarIndexProtoFile(outputFile, index); err != nil {
		return err
	}

	return nil
}
