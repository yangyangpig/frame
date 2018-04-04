#!/bin/bash

protoc -I=. -I=../../../ -I=../../../github.com/gogo/protobuf/protobuf --gogo_out=../bgf/ bgf.proto
