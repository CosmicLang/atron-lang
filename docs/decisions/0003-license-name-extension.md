# ADR-0003: License, project name and source file extension

Status: accepted — Date: 2026-10-17

## Context

Three identity decisions have to be made before the first public commit, because
changing any of them later is expensive: the license governs every
contribution, the name appears in the module path and in every installation
instruction, and the extension appears in every source file and editor
configuration.

The guide recommends Apache-2.0 over MIT, adopts `.atrn` after checking that the
more obvious extensions are taken, and warns that the name `atron` collides with
an unrelated tool. None of these is a legal opinion, and no trademark search has
been performed.

## Decision

### License: Apache-2.0

The project is licensed under the Apache License, Version 2.0. The official text
is stored verbatim in `LICENSE`; it was not reconstructed from memory. Every
source file carries `// SPDX-License-Identifier: Apache-2.0`. Contributions are
accepted under the same license, certified by a DCO `Signed-off-by` trailer
instead of a CLA. `NOTICE` records attribution and third-party components.

### Name: Atron

The project is named **Atron**, and the command is `atron`. The Go module path is
`github.com/CosmicLang/atron-lang`.

The name is used with the knowledge that it is **not unique**: an unrelated
JavaScript package (`atron-js`), an unrelated Python tool (`atron-cli`) that
installs an executable named `atron`, an old TI calculator program, and various
academic and robotics projects use it too. None of them is a programming
language, but `atron-cli` collides at the level of the executable name.

### Extension: `.atrn`

Atron source files use `.atrn`. `.atr` (Atari 8-bit disk images), `.atn` (Adobe
Photoshop actions, among others) and `.atx` (ESRI ArcGIS attribute indexes) were
already in use and were rejected. `.atrc` is reserved for a future serialized
bytecode format and is not used for source.

## Alternatives considered

**MIT instead of Apache-2.0.** MIT is shorter and simpler. It was rejected
because it has no explicit patent grant and no retaliation clause; for a language
implementation, where contributor patents are a plausible risk, Apache-2.0's
patent grant is worth the extra text. Apache-2.0 is also compatible with the
usual Go ecosystem expectations.

**A unique, searchable name.** This is the technically safer option and remains
open: the project could be renamed with an ADR and a module path change. It was
not chosen because a rename now costs nothing and buys nothing concrete, while
the collision that actually matters (`atron-cli`) is addressed by documenting it
rather than by renaming. See "Open issue" below — this decision is the weakest of
the three and the most likely to be revisited.

**`.atr`, `.atn` or `.atx`.** Rejected as already taken; a colliding extension
breaks file association, syntax detection and GitHub Linguist classification.

**No extension check at all.** Rejected: discovering a collision after release
means changing every source file in the wild.

## Consequences

Positive:

- Patent grant and retaliation clause are in place from the first commit.
- The extension is unclaimed in public extension databases at the time of
  writing, so editors can associate it without ambiguity.
- The DCO keeps contribution friction low while preserving provenance.

Negative:

- Apache-2.0 requires preserving `NOTICE` and stating changes; contributors must
  be told this, which `CONTRIBUTING.md` does.
- The name collision means installation instructions must warn users, and a
  user with `atron-cli` installed will find the wrong `atron` on `PATH`. This is
  documented in `README.md` and `GOVERNANCE.md`.
- `.atrn` is four characters and not obviously pronounceable or guessable; it
  will not be found by users typing `.atr`.

## Open issue: clearance before publication

The guide requires the maintainer to verify the name and extension before public
distribution. **That verification has not been completed.** Before the repository
is made public or a release is tagged, the following must be done and recorded
here:

1. Search GitHub for `path:*.atrn` and for the string "Atron language" to confirm
   no competing language uses the extension.
2. Check GitHub Linguist's `languages.yml` and public extension databases
   (FileInfo and similar) for `.atrn`.
3. Perform a trademark search for "Atron" in the relevant jurisdictions, or
   obtain legal advice.
4. Decide whether to keep the name, given the `atron-cli` executable collision,
   or to rename with a migration plan for the module path.

Until step 4 is recorded, this ADR should be read as *accepted for development,
not cleared for public distribution*. `GOVERNANCE.md` states the interim
trademark policy, and `SECURITY.md` and `CODE_OF_CONDUCT.md` state that they are
incomplete until a contact address exists.
