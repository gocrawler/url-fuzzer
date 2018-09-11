# Url_fuzzer

> A simple multi threaded and concurrent url fuzzer implementation in `Go`

### Please Make sure you have `sites.txt` and `lists.txt` on the same directory

## Installation

> **NOTICE:"" The default max speed is 20 connections. if you want to increase or decrease the default speed, please run as `./fuzz -max 50` or `go run fuzzer.go -max 50`



Copy this commands in your terminal 

```txt
go build -o fuzz fuzzer.go

```
After building it once you can run it by

```
./fuzz
```

# Files and assets

You must need to put all your sites in `sites.txt` as shown in the example 
**Be sure you have added `http://` or `https://` prefix on every url**

```
http://site.com
http://example.com
http://target.com
```

You must need to put all your paths in `lists.txt` as shown in the example 
**Be sure you have added `/`  prefix on every url**

```
/hello
/admin/
/shell.php
```


# Results

Your all `found` sites with paths will be saved in `GOT.voot` and after finishing the whole process you will get a html file with clickable links for every result in `index.html`


> You can also get the result in `http://localhost:1339`

> All errors are available on `http://localhost:1339/err`

# Wordlist

> For wordlists check the comments on this gist https://gist.github.com/AnikHasibul/a6de5d07c0b33af7c0ab7ed4abf9dad2

# Enjoy!
