# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: grpc_stdio.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x10grpc_stdio.proto\x12\x06plugin\x1a\x1bgoogle/protobuf/empty.proto\"u\n\tStdioData\x12*\n\x07\x63hannel\x18\x01 \x01(\x0e\x32\x19.plugin.StdioData.Channel\x12\x0c\n\x04\x64\x61ta\x18\x02 \x01(\x0c\".\n\x07\x43hannel\x12\x0b\n\x07INVALID\x10\x00\x12\n\n\x06STDOUT\x10\x01\x12\n\n\x06STDERR\x10\x02\x32G\n\tGRPCStdio\x12:\n\x0bStreamStdio\x12\x16.google.protobuf.Empty\x1a\x11.plugin.StdioData0\x01\x42\x08Z\x06pluginb\x06proto3')



_STDIODATA = DESCRIPTOR.message_types_by_name['StdioData']
_STDIODATA_CHANNEL = _STDIODATA.enum_types_by_name['Channel']
StdioData = _reflection.GeneratedProtocolMessageType('StdioData', (_message.Message,), {
  'DESCRIPTOR' : _STDIODATA,
  '__module__' : 'grpc_stdio_pb2'
  # @@protoc_insertion_point(class_scope:plugin.StdioData)
  })
_sym_db.RegisterMessage(StdioData)

_GRPCSTDIO = DESCRIPTOR.services_by_name['GRPCStdio']
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\006plugin'
  _STDIODATA._serialized_start=57
  _STDIODATA._serialized_end=174
  _STDIODATA_CHANNEL._serialized_start=128
  _STDIODATA_CHANNEL._serialized_end=174
  _GRPCSTDIO._serialized_start=176
  _GRPCSTDIO._serialized_end=247
# @@protoc_insertion_point(module_scope)
