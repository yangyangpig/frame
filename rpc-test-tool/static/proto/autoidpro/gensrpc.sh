#!/bin/bash

protoc -I=. -I=../../../ -I=../../../github.com/gogo/protobuf/protobuf --gogo_out=./ autoidpro.proto
