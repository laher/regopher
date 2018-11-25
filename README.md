# regopher

Regopher is an attempt to bring better refactoring tools to go.

Nothing is ready to use. Just some experiments right now.

Seemingly JetBrains have implemented much of this stuff for Goland, but it's closed source. Also it's probably written in Java, so probably not likely to be features for reuse by most go editors.

## Approach

 * `regopher` uses [dst](https://github.com/dave/dst) (provisionally) to read and manipulate source
 * `regopher` uses [guru](https://godoc.org/golang.org/x/tools/cmd/guru) to infer semantic information from sources - subcommands such as `referrers`, `definition`, `freevars`
 * In/out protocols modelled after `gorename`, `guru` and other Go tools. It should be straightforward to add these refactorings for any editor which already supports `gorename`.

I'm using TDD to gradually refine each refactoring (see `testdata/before/*.go` and `testdata/expected/*.go` ). **If you'd like to file a bug, please try to offer a PR with a failing test case.** *Or even better, offer a fix with it.*

## Editors

Preliminary work on a [vim plugin is here](https://github.com/laher/regopher.vim).

I'm really keen to get going on VS Code and Emacs support, amongst others. Please hit me up if you'd like to work on these integrations. If your editor plugin supports `gorename`, it should be easy to support `regopher` too.

## Status

Most of this work hasn't even started. This is just an idea with legs.

My first 2 refactors are available via vim, but very green. Please only use this tool on code which has been committed to version control or otherwise backed up.

I haven't even started working on multi-package refactoring (package-private refactoring only so far), or even avoiding name collisions. 

Having said all that, this is my first week and I'm pretty hyped about the rate of progress.

## Planned Refactors

Refactor                     | Status  | Use guru?  | Notes
-----------------------------|---------|------------|-------------------
 `params-to-struct`          | started | referrers? | 'Introduce parameter object'. Don't need guru for package-private funcs
 `results-to-struct`         | started | referrers? | 'Introduce result object'. NOTE: exclude last value if it's an error
 Extract variable            |
 Extract constant            |
 Move type to new file       |         |            | 
 Hide                        |         | n/a        | could use referrers to validate that it's unused. Is this just `gorename`?
 Extract function            |         | freevars   | This is already implemented by godoctor
 Remove unnecessary else     |         |            | scan source / block
 Remove unused parameters    |         | referrers  | ''
 Remove unused types/..      |         | referrers  | ''

### As wrapper of existing tools
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
I'm keen to try again here but it works fine for me, so I'll focus on others for now.

## Thanks

This is totally built on the shoulders of giants. Thanks to the Go team for existing tooling, to the `vim-go` team, 
Dave Brophy for `dst`, and the whole community.

## License

[MIT license](LICENSE)
