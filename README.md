# csvabbrev

`csvabbrev` provides a simple way to shrink large CSV files without binary compression. It is geared towards CSV files where a rows tend to repeat values from the rows preceding them.

This README describes [the algorithm employed by](#how-it-works) and [the command-line usage of](#usage) `csvabbrev`.

# Advantages

Unlike gzip and other forms of binary compression, files compressed with `csvabbrev` can be joined with simple file concatenation. This is useful when working with many large CSV files.

Another advantage of `csvabbrev` is that it deals with individual *lines* of the input file. If you kill `gzip` while it is compressing a CSV file, who knows if the resulting file represent valid CSV. With `csvabbrev`, you know that sending a SIGTERM will cause `csvabbrev` to flush the current line before terminating. This is especially useful when piping an infinite stream of CSV data into `csvabbrev`.

# How it works

Suppose you have a CSV file like this:

```csv
Alex,Nichol,6
Alex,Nichol,7
Alex,Nichol,7
Alex,Nichol,5
Alex,Nichol,5
...
```

Clearly, even a half-baked compression algorithm could work some magic on this example. Sometimes, though, binary compression is too heavy-weight. Maybe you want a human-readable document, but you want to "compress" the data nonetheless.

`csvabbrev` shrinks this data by putting a `^` in place of repeated values:

```csv
Alex,Nichol,6
^,^,7
^,^,^
^,^,5
^,^,^
...
```

One issue with this strategy is that the original CSV file might already have `^` in some of its fields. To address this, `csvabbrev` inserts an extra `^` before any ^-only entry (e.g., `^`, `^^`, `^^^`) from the original document. Thus, `^` gets escaped to `^^`, `^^` to `^^^`, etc.

# Usage

To compile this, must have the Go programming language installed and configured correctly. Once you do, run `go install` in this repo to get the `csvabbrev` command. You can check for the presence of this command by viewing the usage:

    $ csvabbrev --help
    Usage: csvabbrev [-h | --help] [-i | --inflate] [file-in [file-out]]

The `csvabbrev` command can operate on standard input/output, or on files. If you do not specify files, it will assume that it is working with the standard input/output streams. Here are multiple ways that you can deflate (i.e. compress) a CSV file:

    $ cat original.csv | csvabbrev >compressed.csv
    $ csvabbrev original.csv compressed.csv
    $ csvabbrev original.csv >compressed.csv

Notice that, if you do not specify a second file argument, `csvabbrev` sends its output to its standard output.

To inflate (i.e. decompress) a file, the usage is almost exactly the same, but you must use the `-i` or `--inflate` flag:

  $ cat compressed.csv | csvabbrev -i >inflated.csv
  $ csvabbrev -i compressed.csv inflated.csv
  $ csvabbrev -i compressed.csv >inflated.csv

## Edge-case behavior

The `csvabbrev` command is designed to run on a continuous stream of input. This means that you can pipe a CSV-generating program into `csvabbrev`. If you do this, `csvabbrev` will forward the deflated output to its own standard output in real time (i.e. without any buffering).

The `csvabbrev` command attempts to shutdown gracefully. If `csvabbrev` is sent a `SIGTERM` or `SIGINT`, it will wait to finish writing its current line of CSV data before terminating.

If the incoming CSV data contains a malformed line, `csvabbrev` will print an error to standard error and terminate. The previous lines which `csvabbrev` has already processed will be written to the destination. However, no further lines will be processed before the command terminates.
