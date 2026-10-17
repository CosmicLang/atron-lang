# ADR-0001: Runtime value representation

Status: accepted (benchmark validation deferred to M4) — Date: 2026-10-17

## Context

The virtual machine must represent dynamic values — nil, bool, int, float,
string, list, map, range, function, native function, struct type, instance,
error and module — without allocating for the scalars that dominate execution.
Two common options in Go:

1. `any` in every field. Simple, but every `int` or `bool` stored in a container
   allocates, and a 64-bit integer cannot be stored in the interface word.
2. A struct carrying the scalar payload in a `uint64` and using `any` only for
   the variants that genuinely need a pointer.

The Atron Master Implementation Guide (chapter 5.2, decision table in 3.2) picks
the second. This ADR records that decision, its alternatives and the measurement
that would overturn it.

## Decision

`object.Value` is a 32-byte struct:

```go
type Value struct {
    K Kind   // type tag
    N uint64 // nil, bool, int, or float via math.Float64bits
    R any    // *Str, *List, *Map, *Closure, *Instance, *Error, *Module, ...
}
```

Associated rules:

- Field access is restricted to package `object`. The rest of the tree uses
  constructors (`Int`, `Float`, `Bool`, `Str`) and accessors (`AsInt`,
  `AsFloat`, `AsBool`, `AsStr`). This is what makes the representation
  replaceable without touching the VM, the compiler or the standard library.
- `Value` is copied by value and owns no mutable state. Variants that hold a
  pointer refer to structures whose mutability is defined by their type: `List`
  and `Map` are mutable through the reference, `Str` is immutable.
- `Nil`, `True` and `False` are constant values (`Value{K: KNil}` and so on).
  They are immutable and safe for concurrent use; they are not mutable global
  state.
- Structural equality is depth-limited (1000) so cyclic structures cannot loop
  forever; exceeding the limit produces error `R0012`.
- Formatting of values lives in exactly one place and arrives in M4, together
  with the first consumer. No `String` method is added at M0, so two divergent
  formatters cannot appear.

## Alternatives considered

**Plain `any`.** Rejected: every `int` appended to a `List` would allocate, so
the 10M-iteration sum loop from the guide would allocate 10M times. Comparing
scalars would also require a type switch on every operation.

**Union via `unsafe.Pointer`.** Rejected: requires `unsafe`, complicates
garbage-collector reachability, and buys nothing over a 32-byte struct for the
scalar cases that matter.

**Interface-tagged variant types (`type Int int64`, …).** Rejected: reintroduces
allocation in containers and turns the VM into a web of type switches.

**Pointer to every value.** Rejected: allocates for `nil`, `true` and small
integers, which are the majority of operands.

## Consequences

Positive:

- `nil`, `bool`, `int` and `float` never allocate.
- The VM manipulates one tag word and one scalar payload word, keeping the
  execution loop simple.
- Encapsulation behind constructors and accessors allows the representation to
  change without rewriting consumers.

Negative and risks:

- 32 bytes per value, versus 16 for a bare interface. A 1024-slot stack costs
  32 KiB. That is acceptable, but it is a cost that must be justified by
  measurement rather than intuition.
- The struct is not comparable with `==` because of the `any` field, so all
  comparison must go through the structural equality function.
- The Go garbage collector sees `R` on every value, which keeps referenced
  objects alive as long as any copy of the `Value` exists. That is the intended
  behaviour, but it means old stack slots must be cleared to release large
  lists.

## Validation

The decision stays **provisional until M4**, when the real VM can run the
benchmarks defined by the guide:

- `BenchmarkAddLoop10M` — summation loop, 10M iterations.
- `BenchmarkFib30` — recursive `fib(30)`.
- `BenchmarkListAppend1M` — list growth.

Criterion: if the 32-byte struct loses to plain `any` by more than 10% on any of
the three, this ADR is superseded by one that records the measurement and the
new representation, keeping the constructor and accessor API intact. Results are
stored under `docs/bench/` and cited from `docs/PROGRESS.md`.

No benchmark was run at M0 because the VM does not exist. Any performance number
attributed to M0 would be fabricated.
