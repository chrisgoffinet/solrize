// Copyright © 2018 Chris Goffinet <chris@threecomma.io>

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	solr "github.com/rtt/Go-Solr"
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export [url]",
	Short: "Export a collection to .json file on disk",
	Args:  cobra.MinimumNArgs(1),
	Run:   exportRun,
}

func exportRun(cmd *cobra.Command, args []string) {
	var (
		host       string
		exportFile string
	)

	host = args[0]
	batchSize, err := cmd.Flags().GetInt("batch-size")
	if err != nil {
		log.Fatal(err)
	}
	collection, err := cmd.Flags().GetString("collection")
	if err != nil {
		log.Fatal(err)
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		log.Fatal(err)
	}
	if output == "" {
		exportFile = fmt.Sprintf("%s.json", collection)
	} else {
		exportFile = output
	}

	file, err := os.Create(exportFile)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Exporting [%s] collection to file: %s\n", collection, exportFile)
	url := fmt.Sprintf("%s/solr/%s", host, collection)
	start := 0
	s := &solr.Connection{URL: url}
	for {
		q := solr.Query{
			Params: solr.URLParamMap{
				"q": []string{"*:*"},
			},
			Rows:  batchSize,
			Start: start,
		}
		res, err := s.Select(&q)

		if err != nil {
			log.Fatal(err)
		}

		if res.Results.Len() == 0 {
			break
		}

		log.Printf("Dumping %d documents...\n", res.Results.Len())

		for i := 0; i < res.Results.Len(); i++ {
			delete(res.Results.Collection[i].Fields, "_version_")
			j, err := json.Marshal(res.Results.Collection[i].Fields)
			if err != nil {
				log.Fatal(err)
			}

			file.Write(j)
			file.WriteString("\n")
		}
		start += batchSize
	}
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	exportCmd.Flags().IntP("batch-size", "b", 500, "batch size to use when exporting")
	exportCmd.Flags().StringP("output", "o", "", "filename to use when writing (default: <collection>.json)")
}
