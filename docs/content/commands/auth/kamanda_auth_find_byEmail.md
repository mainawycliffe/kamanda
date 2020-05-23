---
title: "kamanda auth find byEmail"
slug: kamanda_auth_find_byEmail
url: /commands/kamanda_auth_find_byemail/
summary: "Find a Firebase Auth user by email address"
---
## kamanda auth find byEmail

Find a Firebase Auth user by email address

### Synopsis

Find a Firebase Auth user by email address

```
kamanda auth find byEmail [flags]
```

### Examples

```
kamanda auth find by-email email@example.com
```

### Options

```
  -h, --help   help for byEmail
```

### Options inherited from parent commands

```
      --config string    config file (default is $HOME/.kamanda/config.json)
  -o, --output string    The format in which data will be outputted in [text, json, yaml]
  -P, --project string   The firebase project to use (default "default")
      --token string     firebase token to use for authentication
```

### SEE ALSO

* [kamanda auth find](/commands/kamanda_auth_find/)	 - Find a a Firebase Auth user by their Firebase UID.

###### Auto generated by spf13/cobra on 25-Apr-2020