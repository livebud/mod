# 0.0.6 / 2025-05-18

- support vendored stdlib packages like golang.org/x/net/dns/dnsmessage

# 0.0.5 / 2025-05-18

- expose mod.InStdlib and update std list

# 0.0.4 / 2025-03-25

- add `module.Contains()` function

# 0.0.3 / 2025-03-25

- **BREAKING**: remove embed support
- **BREAKING** rename `module.Directory()` to `mod.Dir()`
- add `mod.Abs(dir)` to lookup absolute directory
- add `module.ResolveDir(importPath)` to resolve a directory
- add `module.ResolveImport(dir)` to resolve an import path
- generalize `mod.Find` to support no dirs and many dirs

# 0.0.2 / 2023-10-07

- add fs.FS support to \*mod.Module
- add precommit hook

# 0.0.1 / 2023-08-14

- Initial commit
