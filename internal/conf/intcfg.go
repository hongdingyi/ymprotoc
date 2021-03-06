package conf

var (
	tmpl = `

###################基础配置#####################

#依赖目录,不用改
includes:
  - $GOPATH/src/github.com/bilibili/kratos/third_party

#存放proto的目录
import_path: $GOPATH/src/youme.im/proto


#当前项目所依赖的proto文件
protos:
   - im/chat/ch.proto
   - im/message/mes.proto


####################编译配置####################
generate:
    go_options:
      extra_modifiers:            #由于mes.proto导入了ch.proto；所以要对ch.proto进行包名的映射   map  可添加多个
        im/chat/ch.proto: youme.im/im/moudledemo/proto/im/chat
    plugins:
      - name: go                  #plugin choice   cpp or java
        type: go                  #用go  还是 gogo
        flags: plugins=grpc       #parameter
        output: ./proto           #output path
`
)
