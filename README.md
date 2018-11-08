# Compilation
All mappings must be converted to block sequences to preserve order.
```yaml
foo:
    fooer: true
bar:
    barer:
    - thing
    - stuff
```
compiles to:
```yaml
- foo:
    - fooer: true
- bar:
    - barer:
      - thing
      - stuff
```
# Unordered List
```yaml
List Name:
  - item name
  - key: value
  - 1
  - true
  #...
```

* **Note**: First item must be a string, or a key-value pair that includes a string.

# Ordered List
```yaml
List Name:
  - 1: item name
  - 2: item name
  #... 
  - 10: item name
```

<!--
Object conversion:
  if <key/value>:
    if <int>: string    -> "<key>. <value>"
    if <string>: string -> "* **<key>**: <value>"
  if <string>           -> "<string>  "
--->