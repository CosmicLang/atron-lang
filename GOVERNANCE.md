# Governance

This document describes how decisions about Atron are made at milestone M0 and
how that model is expected to evolve. It is deliberately short: the project has
one maintainer and no released code, so a heavy process would be fiction.

## Current model: benevolent maintainer

At M0 Atron is governed by a **single maintainer** (the *benevolent maintainer*
model). The maintainer:

- sets the milestone plan and approves milestone completion,
- merges all changes to `main`,
- decides language questions that the specification does not settle,
- appoints committers and reviewers when the contributor base justifies it.

Decisions are made in public in issues and pull requests whenever the discussion
is not sensitive. The maintainer is expected to explain *why* a decision was
made, and to record decisions of lasting consequence as ADRs (see below).

### Why so small

A governance model larger than the contributor base produces process theatre and
slows the project without adding legitimacy. This section will be replaced when
the number of regular contributors makes a committee meaningful, and not before.

## Decision making

1. **Proposals.** Anything that changes the language, its observable behaviour,
   the bytecode format or the public API of `pkg/atron` starts as an issue that
   states the problem, the alternatives and the trade-offs.
2. **Discussion.** Discussion happens in the issue. The goal is consensus among
   the people who will maintain the result.
3. **Record.** Accepted decisions of lasting consequence are written as an ADR
   in `docs/decisions/`, numbered sequentially, and never deleted — a superseded
   ADR is marked as such and points to its replacement.
4. **Implementation.** Only then does code land, in the milestone where it
   belongs.
5. **Deadlock.** If consensus fails, the maintainer decides and records the
   reasoning in the issue. This is the point of having a maintainer.

Reversing an accepted ADR requires a new ADR that explains what changed.

### What needs an ADR

- language syntax or semantics,
- the bytecode instruction set or chunk format,
- package boundaries and the public API surface,
- the choice of a significant dependency,
- anything a future contributor would otherwise have to re-litigate.

### What does not

Bug fixes, refactors that preserve behaviour, tests, documentation and CI
changes are reviewed in the pull request itself.

## Becoming a committer

A contributor may be invited to become a committer after a sustained record of
high-quality contributions — typically several merged pull requests across more
than one area, with good judgement in review. Committers gain the ability to
merge pull requests and to approve milestone work. There is no application
process and no election: the maintainer invites, and the invitation is announced
in the repository.

## Milestones and releases

Work proceeds in milestones M0–M10 as defined by the Atron Master Implementation
Guide. A milestone is complete only when its implementation, tests,
documentation and quality gates all pass; the evidence is recorded in
`docs/PROGRESS.md`. A milestone that is not complete is reported as blocked or
partial rather than as done.

Releases follow Semantic Versioning 2.0. Before 1.0, minor versions may break
compatibility; every breaking change is listed in `CHANGELOG.md`. Version 1.0.0
freezes the language syntax and the `pkg/atron` API.

## Trademark and name policy

The name **Atron** and any project logo are intended to identify this project
and its official distributions.

**Status: not cleared.** No trademark search or registration has been performed,
and no logo exists. Until clearance happens, the project makes no claim of
ownership over the name beyond common-law use in this repository.

Pending clearance, the following policy applies:

- **Allowed without permission:** referring to the project by name in articles,
  talks, tutorials and comparisons; stating that software is "compatible with
  Atron" or "built for Atron"; redistributing unmodified release binaries;
  using the name in the package name of a third-party library that targets
  Atron, provided the package does not imply official endorsement.
- **Requires written permission:** using the name or logo in the name of a fork
  or a modified distribution, in a domain name, in merchandise, or in any way
  that suggests the project endorses a product or service.
- **Never allowed:** claiming authorship of the project, or presenting a modified
  build as an official release.

A fork must be renamed. This protects users from unknowingly running a modified
toolchain under the official name.

### Known name collision

An unrelated Python tool published as `atron-cli` also installs an executable
named `atron`. This is recorded in ADR-0003. The collision affects installation
instructions and tooling discoverability, not the license or the governance of
this project. Any change of project name requires an ADR and a migration plan for
the import path of `pkg/atron`.

## Changing this document

Changes to `GOVERNANCE.md` require the maintainer's approval and are announced in
the repository before taking effect. When the project adopts a committee model,
this document will specify its composition, quorum and voting rules at that
time, rather than speculating about them now.
