![build status](https://travis-ci.org/finkf/gocr.svg?branch=master)
# gocr
Tools for ocr error examination written in [go](https://golang.org).

## Usage

```bash
gocr cat *.gt.txt | gocr align
```

```bash
gocr cat *.gt.txt | gocr align | gocr split
```

Check out the [align](testdata/align_gold.txt) and
[stat](testdata/stat_gold.txt) output files.