# Golangci-lint Configuration and Results

## Configuration

The project now has a comprehensive golangci-lint v2.5.0 configuration in `.golangci.yml`.

### Enabled Linters (19 total)

**Error Checking:**
- `errcheck` - Checks for unchecked errors
- `errorlint` - Finds Go 1.13+ error wrapping issues
- `nilerr` - Finds code returning nil when error is not nil

**Code Quality & Complexity:**
- `gocyclo` - Cyclomatic complexity (max: 15)
- `gocognit` - Cognitive complexity (max: 20)
- `maintidx` - Maintainability index

**Security:**
- `gosec` - Security vulnerability scanner

**Static Analysis:**
- `govet` - Official Go static analysis tool
- `ineffassign` - Detects ineffectual assignments
- `unused` - Checks for unused code
- `staticcheck` - Comprehensive static analysis

**Bugs & Correctness:**
- `bodyclose` - HTTP response body closure
- `sqlclosecheck` - SQL resource closure
- `rowserrcheck` - SQL rows.Err() checks
- `unconvert` - Unnecessary type conversions
- `unparam` - Unused function parameters
- `wastedassign` - Wasted assignments

**Performance:**
- `prealloc` - Slice preallocation opportunities

**Code Clarity:**
- `godot` - Comment punctuation
- `misspell` - Spelling errors
- `revive` - Extensible linter with multiple rules

**Testing:**
- `testifylint` - testify library usage
- `tparallel` - t.Parallel() usage

### Configuration Highlights

```yaml
run:
  timeout: 5m
  tests: true
  go: '1.25'
  
linters-settings:
  gocyclo:
    min-complexity: 15  # Higher for game logic
  
  gocognit:
    min-complexity: 20
  
  gosec:
    excludes:
      - G304  # Asset loading requires file paths

  errcheck:
    check-type-assertions: true
    exclude-functions:
      - (*github.com/veandco/go-sdl2/sdl.Window).Destroy
      - (*github.com/veandco/go-sdl2/sdl.Renderer).Destroy
      - (*github.com/veandco/go-sdl2/sdl.Texture).Destroy
```

## Current Linting Results

Running `golangci-lint run` found **85 issues**:

### Issue Breakdown

- **errcheck: 9 issues** - Unchecked error returns
  - Most are `Destroy()` calls and `file.Close()`
  - Examples: window.Destroy(), renderer.Destroy(), file.Close()

- **godot: 54 issues** - Comments not ending in periods
  - Style issue, low priority
  - Affects all comment blocks

- **gosec: 16 issues** - Security warnings
  - G115: Integer overflow conversions (SDL type conversions)
  - G304: File inclusion via variable (asset loading)
  - G301: Directory permissions (mkdir 0755)

- **revive: 6 issues** - Code style issues

### Priority Assessment

**High Priority:**
- errcheck issues for file operations (should handle file.Close() errors)

**Medium Priority:**
- gosec G115 warnings (integer conversions for SDL - review for safety)

**Low Priority:**
- godot comments (style only)
- gosec G304 (intentional for asset loading)
- gosec G301 (standard directory permissions)

## Running the Linter

```bash
# Run all linters
golangci-lint run

# Run with custom config
golangci-lint run --config .golangci.yml

# Run on specific files
golangci-lint run ./engine/core/...

# Fix auto-fixable issues
golangci-lint run --fix
```

## Next Steps

1. **Address errcheck issues** - Add proper error handling for file.Close()
2. **Review gosec warnings** - Ensure integer conversions are safe
3. **Optional: Fix godot** - Add periods to comment blocks
4. **Integrate into CI** - Add golangci-lint to GitHub Actions/CI pipeline

## Notes

- Test files are excluded from complexity checks (gocyclo, gocognit, maintidx)
- Example files are excluded from errcheck (for clarity)
- Configuration uses explicit linter selection (disable-all: true, then enable specific linters)
- Timeout set to 5 minutes for large codebases
