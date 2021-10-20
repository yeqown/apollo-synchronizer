# apollo-synchronizer
Help developer to sync between local file and remote apollo portal web since portal web is so messy to use

### Features

- [x] download namespaces into local directory.
- [x] synchronize local files to remote apollo portal web.
- [ ] use terminal ui to display synchronization information.
- [ ] apply `force` and `overwrite` flag
- [ ] maybe some customize filter to dismiss some file/namespace?

### Quick start

```sh
go install github.com/yeqown/apollo-synchronizer/cmd/asy@latest
```

### Usage

```sh
# synchronize between one app in apollo with local file system.
$ ./asy -h
NAME:
   apollo-synchronizer - A new cli application

USAGE:
   asy [global options] command [command options] [arguments...]

AUTHOR:
   yeqown <yeqown@gmail.com>

COMMANDS:
   tool     To help developers create, delete and read resources from apollo portal.
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --up                       upload to apollo portal with local filesystem (default: false)
   --down                     download from apollo portal (default: true)
   --force, -f                indicates whether to create the target while it not exists. (default: false)
   --overwrite                indicates whether asy update the target while it exists. (default: false)
   --path value               specify the path to synchronize
   --apollo.portaladdr value  apollo portal address
   --apollo.appid value       the targeted remote app in apollo
   --apollo.secret value      openapi app's token
   --apollo.account value     user id in apollo (default: apollo)
   --apollo.env value         the environment of target remote app (default: DEV)
   --apollo.cluster value     the cluster of target remote app (default: default)
   --debug                    print debug logs (default: false)
   --help, -h                 show help (default: false)
   --version, -v              print the version (default: false) 
```

demo： 

```sh
# download configs from apollo [app=demo] [env=DEV] [cluster=default] 
./asy \
    --down \
    --debug --force --overwrite \
    --path=./debugdata \
    --apollo.portaladdr=http://127.0.0.1:8070 \
    --apollo.appid=demo \
    --apollo.secret=82a95a5722ae8649f64ca5859a13032acab4b2a3
```

> If synchronize the one [file/namespace] those not found in target place, 
> **asy** will create one automatically if you use [-f] [--force] option.


### Structure mapping

```sh
${FOLDER}       				= appID         // you can alse use `--appId` to specify.
├── filename.ext				= namespace.ext
├── ... more
└── service.yaml				= service.yaml
```

### References

- [Apollo Open Documentation](https://github.com/apolloconfig/apollo/wiki/Apollo%E5%BC%80%E6%94%BE%E5%B9%B3%E5%8F%B0)