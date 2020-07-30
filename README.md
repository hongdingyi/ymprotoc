#ymprotoc 

首先对生成的proto.yml进行依赖描述  然后就可编译proto

# 安装
```shell script
1、go get -u github.com/hongdingyi/ymprotoc

//不行采用以下方案
2、git clone https://github.com/hongdingyi/ymprotoc.git
go install

```

# 命令
请查看help

```shell script
ymprotoc -h
Usage:
  ymprotoc [command]

Available Commands:  
  help        Help about any command
  init        生成yaml文件
  build       编译proto

Flags:
  -h, --help   help for ymprotoc

Use "ymprotoc [command] --help" for more information about a command.
```


