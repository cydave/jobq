# jobq

jobq is a simple job queue.

jobq starts out with a set of initial jobs. As the initial jobs are processed by
the workers, new jobs can be submitted from within the workers if needed.

Think of it as a webcrawler-esque job queue:

1. Fetch from initial URL
2. Extract all URLs
3. For each found URL:
   - enqueue URL for crawling

**Note:** by default, jobq deduplicates jobs. Jobs that have been processed are
indexed in an internal map. Jobs with the same ID are
skipped if they are fed back in the job queue. If you have no need for
deduplication, use `jobq.NewDup[T]()` instead of `jobq.New[T]()`.

## Example

Find examples under [examples/](examples/)

## Credits

The implementation of jobq is based off of the algorithm by [qvl/httpsyet](https://github.com/qvl/httpsyet).
