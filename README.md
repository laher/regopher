# regopher

Regopher is an attempt to bring better refactoring tools to go.

Nothing is ready to use. Just some experiments right now.

Seemingly JetBrains have implemented much of this stuff for Goland, but it's closed source. Also it's probably written in Java, so probably not likely to be features for reuse by most go editors.

## Approach

 * `regopher` uses [dst](https://github.com/dave/dst) (provisionally) to read and manipulate source
 * `regopher` uses [guru](https://godoc.org/golang.org/x/tools/cmd/guru) to infer semantic information from sources - subcommands such as `referrers`, `definition`, `freevars`

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
 Remove unnecessary else     |         |           | scan source / block
 Remove unused parameters    |         | referrers | ''
 Remove unused types/..      |         | referrers | ''

#### As wrapper of existing tools
 * Expose method/variable (as above)
 * Move to another package (maybe possible to wrap rename)

### Fixes

Fixes suggest that the code is not currently valid. Not sure criteria yet for generating an AST for slightly-broken code.

Fix                                           | Feasibility |
----------------------------------------------|-------------|
Match signature to return values under cursor | ?


## Out of scope

Generally, there's no need to recreate things which already exists. 
I'll avoid things which have already been adopted by vim-go, vscode, etc.

 * gorename
 * go-iferr
 * goimpl/impl
 * struct tagging
 * keyify

extract-function is available via godoctor but seems not to be embraced by those at vim-go and vscode.
