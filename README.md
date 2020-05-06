# frontserver
golang으로 구성된 server 로직입니다. 

main 모드인 경우 단일 서버로 동작하며 로직도 서버에서 처리합니다.

front 모드인 경우 목적지 서버로 데이터를 전달하는 역활을 합니다. scale out으로 확장성 있도록 구성됩니다. 

## 필요 패키지 정보

### protobuf

  go get -d -u github.com/golang/protobuf/proto
 
  go get -d -u github.com/golang/protobuf/protoc-gen-go
 
  go install github.com/golang/protobuf/protoc-gen-go
 
### 클라이언트 UI관련 패키지 

 go get -d -u github.com/nsf/termbox-go

## Quick Links
* [설치방법](https://github.com/woong20123/frontserver/wiki/FrontServer_Set_Project)
* [설명](https://github.com/woong20123/frontserver/wiki/FrontServer_doc) 
* [ExampleServer](https://github.com/woong20123/frontserver/wiki/FrontServer_ExampleServer)


