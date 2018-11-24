# regopher

Refactoring tools for go [experimental].

Nothing is ready to use. Just some experiments right now.

## Approach

Regopher is an attempt to bring better refactoring tools to go.

 * Interprets `guru` results - `referrers`, `definition`, `freevars`
 * Uses `dst` (provisionally) to manipulate source

## Ideas

### Refactors

Refactor                     | Status  | Use guru? | Notes
-----------------------------|---------|-----------|-------------------
 Introduce parameter object  | started | referrers | Probably don't need guru for package-private funcs
 Introduce result object     |         | referrers | as above
 Extract variable            |
 Extract constant            |
 Move to a new file          |         |           | 
 Hide                        |         | n/a       | could use referrers to validate that it's unused. Is this just `gorename`?
 Extract function            |         | freevars  | This is already implemented by godoctor

#### As wrapper of existing tools
 * Expose method/variable (as above)
 * Move to another package (maybe possible to wrap rename)

### Fixes

Fixes suggest that the code is not currently valid. Not sure criteria yet for generating an AST for slightly-broken code.

Fix                                           | Feasibility |
----------------------------------------------|-------------|
Match signature to return values under cursor | ?


## Out of scope

 * gorename
 * iferr
