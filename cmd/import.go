// Copyright © 2018 Chris Goffinet <chris@threecomma.io>

package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	solr "github.com/rtt/Go-Solr"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a .json file from disk to a collection",
	Args:  cobra.MinimumNArgs(1),
	Run:   importRun,
}

func process(url *string, batch []map[string]interface{}) {
	s := &solr.Connection{URL: *url}
	f := map[string]interface{}{
		"add": batch,
	}
	_, err := s.Update(f, true)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Importing %d documents...\n", len(batch))
	}
}

func worker(url *string, wg *sync.WaitGroup, bufCh chan map[string]interface{}, batchSize int) {
	var batch []map[string]interface{}
	defer wg.Done()

	for {
		select {
		case v, ok := <-bufCh:
			if !ok {
				// anything left to process, finish it
				if len(batch) > 0 {
					process(url, batch)
				}
				return
			}
			batch = append(batch, v)
			if len(batch) >= batchSize {
				process(url, batch)
				batch = batch[:0]
			}
		}
	}
}

func importRun(cmd *cobra.Command, args []string) {
	var (
		host = args[0]
		wg   = new(sync.WaitGroup)
	)

	batchSize, err := cmd.Flags().GetInt("batch-size")
	if err != nil {
		log.Fatal(err)
	}
	collection, err := cmd.Flags().GetString("collection")
	if err != nil {
		log.Fatal(err)
	}

	inputFile, err := cmd.Flags().GetString("input")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.Printf("Importing [%s] collection from file: %s\n", collection, inputFile)
	url := fmt.Sprintf("%s/solr/%s", host, collection)

	wg.Add(1)

	bufCh := make(chan map[string]interface{}, batchSize)
	go worker(&url, wg, bufCh, batchSize)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var doc = make(map[string]interface{})
		json.Unmarshal(scanner.Bytes(), &doc)
		bufCh <- doc
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	close(bufCh)
	wg.Wait()
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().IntP("batch-size", "b", 500, "batch size to use when importing")
	importCmd.Flags().StringP("input", "i", "", "filename to use when importing")
	importCmd.MarkFlagRequired("input")
}
