# apollo-synctool
Help developer to sync between local file and remote apollo portal web since portal web is so messy to use


### Quick start

```sh
go install github.com/yeqown/apollo-synchronizer/cmd/asy@latest
```

### Usage

```sh
# synchronize between one app in apollo with local file system.
asy \ 
	--up / --down		: download or upload
	--secret    		: specified the registered app token
	[-f] --file 		: specify files to synchronize 
	[-d] --path 		: specify the folder to synchronize
	--appId     		: use conf of target app (appId and etc)
	[-f] --force		: automatically create [file/namespace] in case of target resource is not found.
	--overwrite 		: overwrite the target [file/namespace] if it has existed. 
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