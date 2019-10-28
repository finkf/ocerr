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

For the output format check out the [align](testdata/align.gold.txt)
and [stat](testdata/stat.gold.txt) test files.
