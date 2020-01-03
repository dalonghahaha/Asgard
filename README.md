# Asgard

This is Asgard Framework

## Remark

### protoc cmd

``` bash
protoc -I protos protos/base.proto --go_out=plugins=grpc:./rpc/

protoc -I protos protos/app.proto --go_out=plugins=grpc:./rpc/

protoc -I protos protos/job.proto --go_out=plugins=grpc:./rpc/
```
