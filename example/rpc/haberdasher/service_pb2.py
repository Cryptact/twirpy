# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: example/rpc/haberdasher/service.proto
# Protobuf Python Version: 5.27.3
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    27,
    3,
    '',
    'example/rpc/haberdasher/service.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n%example/rpc/haberdasher/service.proto\x12\x14twitch.twirp.example\"0\n\x03Hat\x12\x0c\n\x04size\x18\x01 \x01(\x05\x12\r\n\x05\x63olor\x18\x02 \x01(\t\x12\x0c\n\x04name\x18\x03 \x01(\t\"\x16\n\x04Size\x12\x0e\n\x06inches\x18\x01 \x01(\x05\x32O\n\x0bHaberdasher\x12@\n\x07MakeHat\x12\x1a.twitch.twirp.example.Size\x1a\x19.twitch.twirp.example.HatB\nZ\x08/exampleb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'example.rpc.haberdasher.service_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\010/example'
  _globals['_HAT']._serialized_start=63
  _globals['_HAT']._serialized_end=111
  _globals['_SIZE']._serialized_start=113
  _globals['_SIZE']._serialized_end=135
  _globals['_HABERDASHER']._serialized_start=137
  _globals['_HABERDASHER']._serialized_end=216
# @@protoc_insertion_point(module_scope)
