# Test data

Golden and fixture files used by the test suite.

**Status: empty at M0.** No test in the repository reads from this directory yet;
`cmd/atron` tests are self-contained.

## Layout (planned)

```
testdata/
├── run/       # programs that must execute successfully
│   └── *.atrn + *.golden   (stdout + exit code)
├── errors/    # programs that must fail with a specific diagnostic
│   └── *.atrn + *.golden   (rendered diagnostic, with line/column)
├── spec/      # one file per specification example, chapter 4 of the guide
└── fuzz/      # seed corpora for FuzzLexer, FuzzParser, FuzzVM
```

## Rules

- Golden files are reviewed by a human when they change. Regenerating them with
  a `-update` flag and committing the result without reading the diff is not
  acceptable: it hides exactly the regressions the tests exist to catch.
- A golden file records the expected output exactly, including trailing
  newlines and exit status. Comparison is byte-for-byte.
- Error goldens include the full rendered diagnostic — code, message, source
  excerpt and caret line — because the rendering is part of the user-visible
  contract.
- Fixtures must be deterministic: no timestamps, random values, absolute paths
  or Go map iteration order in the expected output.
- Non-UTF-8 fixtures, if any, are named with a `.bin` suffix so tools do not try
  to render them.
