#!/usr/bin/env bash

echo 'Compiling Start.'

python -m grpc_tools.protoc  --proto_path=api --python_out=api --grpc_python_out=api api/common/common.proto
python -m grpc_tools.protoc  --proto_path=api --python_out=api --grpc_python_out=api api/adservice/adservice.proto
python -m grpc_tools.protoc  --proto_path=api --python_out=api --grpc_python_out=api api/fleet/fleet.proto
python -m grpc_tools.protoc  --proto_path=api --python_out=api --grpc_python_out=api api/library/library.proto
python -m grpc_tools.protoc  --proto_path=api --python_out=api --grpc_python_out=api api/marketing/marketing.proto
python -m grpc_tools.protoc  --proto_path=api --python_out=api --grpc_python_out=api api/ptransit/ptransit.proto
python -m grpc_tools.protoc  --proto_path=api --python_out=api --grpc_python_out=api api/rideshare/rideshare.proto
python -m grpc_tools.protoc  --proto_path=api --python_out=api --grpc_python_out=api api/routing/routing.proto

python -m grpc_tools.protoc  --proto_path=api --python_out=api --grpc_python_out=api api/synerex.proto

python -m grpc_tools.protoc  --proto_path=nodeapi --python_out=nodeapi --grpc_python_out=nodeapi nodeapi/nodeid.proto

echo 'Compiling Completed.'