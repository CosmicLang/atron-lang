# Examples

Demonstration programs written in Atron.

**Status: empty at M0.** The language cannot execute anything yet, so there is
no runnable example. Each file that lands here is added by the milestone that
makes it run, and is executed by CI from that point on.

## Planned examples

| File | Demonstrates | Milestone |
|---|---|---|
| `hello.atrn` | printing, the entry point, `println` | M4 |
| `fib.atrn` | recursion, functions, integers | M4 |
| `counter.atrn` | closures and captured variables | M5 |
| `vec.atrn` | lists, iteration, mutation | M5 |
| `wordfreq.atrn` | maps, string handling, ordering | M5 |
| `errors.atrn` | error values and the `?` operator | M6 |
| `modules/` | imports across files and the standard library | M7 |

## Rules

- Every example must run with the current `atron` binary and produce the output
  stated in its header comment. CI verifies this; an example that does not run
  is a broken build, not a documentation issue.
- Examples are small enough to read in one screen where possible. If one grows
  beyond that, it belongs in a tutorial under `docs/`, not here.
- No example may depend on network access, the wall clock or the host file
  system layout unless it exists specifically to demonstrate that capability,
  and then it must say so in its header.
