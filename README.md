# gpq

[![License](https://img.shields.io/github/license/seankhliao/gpq.svg?style=flat-square)](LICENSE)
![Version](https://img.shields.io/github/v/tag/seankhliao/gpq?sort=semver&style=flat-square)

GoProxyQuery simple cli to query the go {index,proxy,sum} servers

## Usage

```bash
# show index
gpq index
# limit index output
gpq index -limit 10
# limit index age
gpq index -since 2019-01-01T00:00:00.000Z

# list the versions of a module
gpq proxy example.com/module
# show the info of a particular version
gpq proxy example.com/module v0.1.0
# download the zip file
gpq proxy -save example.com/module v0.1.0

gpq sum example.com/module v0.1.0
```

## TODO

- [ ] implement paging for index
- [ ] verify sums
- [x] smarter input: accept module@vers
- [ ] accept dir for sums?
- [ ] verbose mode
