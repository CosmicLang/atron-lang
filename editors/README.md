# Editor support

Syntax highlighting and language integration for Atron.

**Status: empty at M0.** No grammar exists because the language is not
implemented; publishing a grammar now would describe a language that does not
exist and would drift immediately.

## Planned layout

```
editors/
├── textmate/     # TextMate grammar (atrn.tmLanguage.json)
└── vscode/       # VS Code extension wrapping the TextMate grammar
```

The TextMate grammar is the primary artifact: it can be consumed by VS Code,
Sublime Text, Atom-compatible editors and anything else that reads
`.tmLanguage.json`. The VS Code extension is a thin wrapper around it, not a
second grammar.

## When this arrives

- The grammar is generated from, or at minimum reviewed against, the token
  definitions in `internal/token` so that highlighting and lexing agree.
- It ships with the 0.1.0 release (milestone M10) at the latest, or as soon as
  the syntax is stable enough to be worth highlighting.
- A GitHub Linguist submission for the `.atrn` extension is a separate step and
  requires evidence of real-world usage.

## Near-term workaround

Until a grammar exists, editors can be configured manually. For example, in VS
Code, associate `.atrn` with a similar language in `settings.json`:

```json
{
  "files.associations": {
    "*.atrn": "javascript"
  }
}
```

This is a stopgap for reading files, not real support: there is no Atron
grammar, no completion and no diagnostics.
