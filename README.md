# jobq

jobq is a simple job queue.

jobq starts out with a set of initial jobs. As the initial jobs are processed by
the workers, new jobs can be submitted from within the workers if needed.

Think of it as a webcrawler-esque job queue:

1. Fetch from initial URL
2. Extract all URLs
3. For each found URL:
   - enqueue URL for crawling

## Example

A simple example on how jobq is used can be found under [examples/example.go](examples/simple/example.go).

## Credits

The implementation of jobq is based off of the algorithm by [qvl/httpsyet](https://github.com/qvl/httpsyet).
