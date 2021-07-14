# gonia | mania star + pp calculator
a very fast and accurate star + pp calculator for mania. gonia has low memory usage and very fast calculation times and currently has the fastest (for mania).

* [installation (windows)](#installation-windows)
* [installation (linux)](#installation-linuxubuntu)
* [usage](#usage)
* [alternatives](#alternatives)

# installation (windows)
you'll need [go](https://golang.org/doc/install) installed, to build the project. you'll also need [git](https://git-scm.com/downloads) installed to download the project. once all that is cleared up, you can go ahead and install the project

```
C:\> go get github.com/barrack-obama/gonia.git
C:\> cd $GOPATH/gonia/gonia
$GOPATH\gonia\gonia>
```

once you're in the folder, you can go ahead and build the porject.

```
$GOPATH\gonia\gonia> go build .
$GOPATH\gonia\gonia> ls gonia.exe

    Directory: $GOPATH\gonia\gonia

Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
-a---l         7/14/2021   5:08 PM        2284544 gonia.exe
```

then you can go down to [usage](#usage) for details on how to use it.

# installation (linux/ubuntu)
same as for windows, you'll need [git](https://www.atlassian.com/git/tutorials/install-git) and [go](https://www.vultr.com/docs/install-the-latest-version-of-golang-on-ubuntu) installed.

it's actually the same build instructions for linux as it is for windows.
```
~$ go get github.com/barrack-obama/gonia.git
~$ cd $GOPATH/gonia/gonia

$GOPATH/gonia/gonia$ go build .
$GOPATH/gonia/gonia$ ls gonia
gonia
```

# usage
you can see gonias usage by running it

```
$GOPATH/gonia/gonia$ ./gonia

parse conf: no arguments to parse                             <-- parsing errors
usage: ./gonia /path/to/beatmap.osu [score]s +[mods (int)]    <-- usage
example: ./gonia /home/simon/beatmaps/2220863.osu 993344s +64 <-- example

`path to beatmap` parameter is case sensitive, so if it says
`parse error: file not found`, check if your spelling is correct.
also score and mods are optional parameters, and the pp will display
as a perfect score with no mods.

```

# alternatives
[omppc](https://github.com/semyon422/omppc) written in lua

[maniera](https://github.com/NiceAesth/maniera) written in python

i suggest going with maniera, since it's most up to date with osu!'s star calculation, though it's not as accurate as gonia.









