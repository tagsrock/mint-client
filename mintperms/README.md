# mintperms

Convert between human-readable string representations and machine readable integer representations of permissions

Examples

```
> mintperms int call:0 send:1 name:1
Perms and SetBit (As Integers)
66      70

Perms and SetBit (As Bitmasks)
1000010 1000110


> mintperms string 66 70
send: true
call: false
name: true
```

