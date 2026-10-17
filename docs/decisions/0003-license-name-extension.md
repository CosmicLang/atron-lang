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
  writing, so editors can associate it without ambiguity. This was re-verified
  in M0-R5: see "Clearance research" below for the queries and their results.
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
distribution. **That verification is only partly complete:** the non-legal checks
(steps 1 and 2) were performed in M0-R5 and are recorded below, while trademark
clearance and the final keep-or-rename decision (steps 3 and 4) remain open.
Before the repository is made public or a release is tagged, the following must
be finished and recorded here:

1. ~~Search GitHub for `path:*.atrn` and for the string "Atron language" to confirm
   no competing language uses the extension.~~ **Done (M0-R5).** No competing
   language and no existing `.atrn` file were found; see the research section.
2. ~~Check GitHub Linguist's `languages.yml` and public extension databases
   (FileInfo and similar) for `.atrn`.~~ **Done (M0-R5).** Absent from Linguist
   and from FileInfo; see the research section for the one ambiguous hit.
3. Perform a trademark search for "Atron" in the relevant jurisdictions, or
   obtain legal advice. **Outstanding.**
4. Decide whether to keep the name, given the `atron-cli` executable collision,
   or to rename with a migration plan for the module path. **Outstanding** —
   deliberately left to the maintainer, since M0-R5 found no conflict that
   forces the decision.

### Related governance items, not part of this decision

Two items are often mentioned alongside clearance. They are different in kind
and must not be conflated:

| Item | Status | Detail |
|---|---|---|
| Moderation contact in `CODE_OF_CONDUCT.md` | **Resolved (M0-R2)** | The Enforcement section publishes a working contact address. An earlier revision of this ADR said the file was incomplete "until a contact address exists"; that statement is obsolete and has been removed. Nothing about the Code of Conduct blocks publication. |
| Private vulnerability reporting (PVR) | **Resolved (M0-R6)** | `GET /repos/CosmicLang/atron-lang/private-vulnerability-reporting` returns `{"enabled": true}` with HTTP 200; re-verified on 2026-10-17. `SECURITY.md` describes the operational private channel. Activation required the repository *Admin* role, so it was a settings action by the maintainer, not a documentation change. |

Name and extension clearance neither depends on nor unblocks the code-of-conduct
contact or PVR: the first is a naming question, the others are published-contact
and security-channel configurations. The project remains *not cleared for public
distribution* on the strength of step 3 alone, not because of either of these
rows.

Until step 3 is performed and step 4 is recorded, this ADR is *accepted for
development, not cleared for public distribution*. `GOVERNANCE.md` states the
interim trademark policy.

## Clearance research (M0-R5)

Performed on 2026-10-17 with read-only, publicly verifiable requests. The scope
is deliberately limited to what can be checked without legal advice: whether a
competing *language* uses the name, whether the extension is registered, and
whether the executable name collides.

### Confirmed conflicts

| Finding | Evidence | Assessment |
|---|---|---|
| A PyPI package installs an executable named `atron` | `atron-cli` 0.0.1, `entry_points.txt` inside `atron_cli-0.0.1-py2.7.egg` declares `atron=atron_cli.atron:cli`. https://pypi.org/pypi/atron-cli/json | **Material, expected, accepted.** A user who installs it finds the wrong `atron` on `PATH`. Not a language. Documented in `README.md` and `GOVERNANCE.md`; the mitigation is a warning, not a rename. |
| The bare name `atron` is taken in npm | `npm view atron` → 1.0.13-a, described as an Electron/Angular framework. https://registry.npmjs.org/atron | **Immaterial for this project.** The command is distributed as a Go binary, not an npm package, and `atron-lang` is the module path. |
| `atron-js` exists in npm | 1.1.2, string/number utility functions. https://registry.npmjs.org/atron-js | Immaterial: a utility library, not a language. |
| Many unrelated repositories use the name | `GET /search/repositories?q=atron+in:name` → 305 results, e.g. `RISC-NYUAD/ATRON` (USV debris-collection system, Python), `skylix-dev/Atron` (React desktop apps), `atron-cc/atron-cli` | **Expected.** None is a programming language. The name is not unique and was never claimed to be. |

### Inconclusive, recorded as such

| Check | Observation | Why it is inconclusive |
|---|---|---|
| Trademark search for "Atron" | Not performed | Requires legal advice and jurisdiction scope. **No clearance is claimed.** |
| `file-extension.org` `.atrn` page | Returns HTTP 200 with a page titled "File extension ATRN" | **Not evidence of a real format.** The page has no file type, no developer and no opening program, and is dated 01/01/1970. A control request for a nonsense extension (`qqzzxx`) also returns 200, so the site auto-generates pages. Recorded to avoid a later reader mistaking it for a registered format. |
| `fileinfo.com` `.atrn` | HTTP 404. The sibling `.atr` page resolves to "ATR File". | Consistent with `.atrn` being unregistered, but absence on one site is not proof of global non-use. |

### No conflict found

| Check | Method | Result |
|---|---|---|
| Competing language named Atron | `GET /search/repositories?q="atron language"` and `q=atron+topic:programming-language` | 0 results each. No programming language uses the name. |
| `.atrn` in use anywhere on GitHub | Code search `extension:atrn` → **0**; `path:.atrn` → **0**; `atron+in:path+extension:atrn` → **0** | Zero is a real result, not a failed query: the controls `extension:go` and `extension:rs` return 45,219,840 and 31,260,672 on the same endpoint. |
| `.atrn` in GitHub Linguist | `languages.yml` has 9,676 lines and a valid `Go:` entry at line 2833; no `atrn` or `Atron` match. Extensions beginning `.at` present are only `.atomignore` and `.ats`. | The extension is unclaimed by Linguist, so syntax highlighting will need an upstream contribution rather than conflicting with an existing language. |
| `.atrn` as a Go module | `https://proxy.golang.org/github.com/atron/atron/@v/list` → 404 | No module occupies the namespace. |
| `atron` in Homebrew | Searched all 8,661 formulae in `formula.json`; no name contains `atron` | No formula collision. |
| PyPI `atron` (bare) | `https://pypi.org/pypi/atron/json` → 404 | Only `atron-cli` exists. |
| `.atr`, `.atn`, `.atx` already taken | `.atr` → "ATR File" on fileinfo.com | Confirms the original rejection of these extensions was correct. |

### Conclusion

Nothing found requires renaming the project before publication, and the
`atron-cli` collision is a documented packaging caveat rather than a conflict
that forces a new name. Because no trademark search has been performed, this
section is **research, not clearance**: the license, project name, source
extension and module path are unchanged, and the decision to keep the name
remains the maintainer's — see step 4 above.
