![build status](https://travis-ci.org/finkf/ocerr.svg?branch=master)
# ocerr
Tools for ocr error examination written in [go](https://golang.org).

## Usage

```bash
ocerr cat *.gt.txt | ocerr align
```

```bash
ocerr cat *.gt.txt | ocerr align | ocerr split
```

For the output format check out the [align](testdata/align.gold.txt)
and [stat](testdata/stat.gold.txt) test files.
