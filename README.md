# Solrize
Simple utility to import and export a collection from Solr

## Installation

```
$ go get -v -d
$ go install
$ solrize
```

## Usage
```Simple utility to import and export a collection from Solr

Usage:
  solrize [command]

Available Commands:
  export      Export a collection to .json file on disk
  help        Help about any command
  import      Import a .json file from disk to a collection

Flags:
  -c, --collection string   the name of the collection from solr
  -h, --help                help for solrize

Use "solrize [command] --help" for more information about a command.```