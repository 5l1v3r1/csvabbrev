# csvabbrev

**csvabbrev** provides a simple way to shrink large CSV files without binary compression. It is geared towards CSV files where a rows tend to repeat values from the rows preceding them.

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

**csvabbrev** shrinks this data by putting a `"` in place of repeated values:

```csv
Alex,Nichol,6
",",7
",","
",",5
",","
...
```

One issue with this strategy is that the original CSV file might already have `"` in some of its fields. To address this, **csvabbrev** inserts an extra `"` before any quote-only entry (e.g., `"`, `"""`, `""`) from the original document. Thus, `"` gets escaped to `""`, `""` to `"""`, etc.
