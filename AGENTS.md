# OpenCode Agent Guidelines (genutil-go)

This repository contains advanced generic utilities for Go. An agent should treat this code as highly specialized and must adhere to these non-obvious patterns to avoid bugs or inefficiency.

## 🚨 Critical Gotchas & Patterns
*   **Error Handling:** Do not use standard `if err != nil` checks in core logic. Use the dedicated `errors/*` package:
    *   Use `defer Catch(&err)` around a block of code if multiple functions might return errors, as this ensures panics from failed operations are recovered and assigned to `err`.
    *   Use `Try()` to execute an operation that returns a value/error tuple, stripping the error part and relying on panic recovery for failure.
    *   If error wrapping is needed, use `Handle` instead of `Catch`.
*   **Slicing Mutability:** When using mutable iterators over slices (`slices.IterMut(&s)`), **you MUST pass a pointer to the slice**. Failing to do so may cause reallocation and invalidating the iterator's view of the data.
*   **Option Type Safety:** Never call `.Get()` on an `opt.Val`, `opt.Ref`, `opt.Opt` without first checking if it is non-empty using `.IsEmpty()`. This will result in a runtime panic. The preferred functional pattern for safe modification is using `.Morph()`.
    *   `option.Option` is deprecated new code should not use it.

## ⚙️ Workflow & Technical Quirks
### Testing and Verification
*   **Full Coverage:** Use `go test ./...` to run all tests across the entire repository.
*   **Targeted Tests:** When testing specific functionalities, utilize specialized packages (e.g., `errors/test`) which wrap results for assertion within the standard `testing` package context (e.g., `test.Result(os.Open("file")).Must(t)`).

### Iterators and Consumption
*   **Preferred Loop:** When iterating in a loop body, use `.Seq()` to convert an iterator to a native Go `iter.Seq`. This method is generally more performant than calling the raw `Next()` method repeatedly on the iterator object.
    ```go
    // Recommended:
    for n := range iterator.Range(0, 5).Seq() { ... }
    ```
*   **Collection:** Use dedicated methods like `.Collect()` or functional helpers (`iterator.CollectMap()`) rather than manual loops to gather elements into slices/maps.

### Maps (Nested Structures)
*   **Path Manipulation:** For complex structures (nested maps), always use the path functions: `maps.PutPath(m, []string{"a", "b"}, val)` for insertion, and `maps.GetPath(m, []string{"a","b"})` for retrieval. These handle deep path creation/checking automatically.

## 🗺️ Architecture Summary
*   **Core Abstractions:** The library builds upon several core generic primitives:
    *   **Tuple:** Generic type system providing fixed-size tuples (`Of2`, `Pair`). Note that a size 0 tuple is a unit type.
    *   **Slices/Maps:** Highly functional, offering transformations (`Map`, `Filter`), ranges (`Range`, `IncRangeBy`, etc.), and utilities for safe access (e.g., `maps.Keys(m)`).
    
## Notes on documentation
The project makes use of some less well known features in godoc:
* Document links in the form `[symbol]` or `[package.symbol]` render as links to the target symbol, in godoc's generated HTML documentation and in vscode mouseover popup (and probably other contexts). Note that if the referenced symbol has a type parameter, it must not be included within the square bracket. It can be placed adjacent outside the bracket, e.g. `[opt.Val][T]`.
* Example test cases which are included as part of godoc's generated HTML documentation. These are unit tests appearing in standard unit test files whose names match the patterns `Example()` for package examples, `ExampleTypename()` for type scoped examples, `ExampleFuncname()` for function examples and `ExampleTypename_Methodname()` for method examples. Any of the above patterns can include a `_suffix` suffix if there are multiple examples for the same entity.


## ⚠️ Summary of Forbidden Practices (Do Not Guess)
1.  DO NOT assume standard library error handling is sufficient; use the dedicated `errors/*` package logic.
2.  DO NOT iterate over mutable slices without passing a pointer (`&s`) to the underlying slice variable.
3.  DO NOT access an Option's value with `.Get()` unless guaranteed non-empty.
  
## 🚧 Work in progress
The `opt` package is a new package being developed to replace the `option` package that better handles the distinction between by value and by reference 
optional types.
