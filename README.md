# regopher

Refactoring tools for go [experimental].

Nothing is ready to use. Just some experiments right now.

## Approach

Regopher is an attempt to bring better refactoring tools to go.

 * Interprets `guru` results - `referrers`, `definition`, `freevars`
 * Uses `dst` (provisionally) to manipulate source

## Ideas

### Refactors

 * Extract function
 * Introduce parameter object
 * Introduce result object
 * Extract variable
 * Extract constant
 * Move to another file

#### As wrapper of existing tools
 * Hide method/variable (it's really a subset of rename) 
 * Expose method/variable (as above)
 * Move to another package (maybe possible to wrap rename)

### Fixes
 * Match signature to return values under cursor


## Out of scope

 * gorename
 * iferr
