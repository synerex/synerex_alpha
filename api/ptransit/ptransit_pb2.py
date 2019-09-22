# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: ptransit/ptransit.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import duration_pb2 as google_dot_protobuf_dot_duration__pb2
from api.common import common_pb2 as common_dot_common__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='ptransit/ptransit.proto',
  package='api.ptransit',
  syntax='proto3',
  serialized_options=_b('Z-github.com/synerex/synerex_alpha/api/ptransit'),
  serialized_pb=_b('\n\x17ptransit/ptransit.proto\x12\x0c\x61pi.ptransit\x1a\x1egoogle/protobuf/duration.proto\x1a\x13\x63ommon/common.proto\"\xa3\x03\n\tPTService\x12\x13\n\x0boperator_id\x18\x01 \x01(\x05\x12\x0f\n\x07line_id\x18\x02 \x01(\x05\x12\x17\n\x0fpast_station_id\x18\x03 \x01(\x05\x12\x18\n\x10station_group_id\x18\x04 \x01(\x05\x12\x17\n\x0fnext_station_id\x18\x05 \x01(\x05\x12\x19\n\x11next_station_name\x18\x06 \x01(\t\x12\x12\n\nvehicle_id\x18\x07 \x01(\x05\x12\r\n\x05\x61ngle\x18\x08 \x01(\x02\x12\r\n\x05speed\x18\t \x01(\x05\x12+\n\x10\x63urrent_location\x18\n \x01(\x0b\x32\x11.api.common.Place\x12\x36\n\x1cnext_arraival_timetable_time\x18\x0b \x01(\x0b\x32\x10.api.common.Time\x12-\n\x13past_departure_time\x18\x0c \x01(\x0b\x32\x10.api.common.Time\x12-\n\ndelay_time\x18\r \x01(\x0b\x32\x19.google.protobuf.Duration\x12\x14\n\x0cvehicle_type\x18\x0e \x01(\x05\")\n\x06PTgtfs\x12\x11\n\tgtfs_name\x18\x01 \x01(\t\x12\x0c\n\x04gtfs\x18\x02 \x01(\x0c\x42/Z-github.com/synerex/synerex_alpha/api/ptransitb\x06proto3')
  ,
  dependencies=[google_dot_protobuf_dot_duration__pb2.DESCRIPTOR,common_dot_common__pb2.DESCRIPTOR,])




_PTSERVICE = _descriptor.Descriptor(
  name='PTService',
  full_name='api.ptransit.PTService',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='operator_id', full_name='api.ptransit.PTService.operator_id', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='line_id', full_name='api.ptransit.PTService.line_id', index=1,
      number=2, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='past_station_id', full_name='api.ptransit.PTService.past_station_id', index=2,
      number=3, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='station_group_id', full_name='api.ptransit.PTService.station_group_id', index=3,
      number=4, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='next_station_id', full_name='api.ptransit.PTService.next_station_id', index=4,
      number=5, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='next_station_name', full_name='api.ptransit.PTService.next_station_name', index=5,
      number=6, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='vehicle_id', full_name='api.ptransit.PTService.vehicle_id', index=6,
      number=7, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='angle', full_name='api.ptransit.PTService.angle', index=7,
      number=8, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='speed', full_name='api.ptransit.PTService.speed', index=8,
      number=9, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='current_location', full_name='api.ptransit.PTService.current_location', index=9,
      number=10, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='next_arraival_timetable_time', full_name='api.ptransit.PTService.next_arraival_timetable_time', index=10,
      number=11, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='past_departure_time', full_name='api.ptransit.PTService.past_departure_time', index=11,
      number=12, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='delay_time', full_name='api.ptransit.PTService.delay_time', index=12,
      number=13, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='vehicle_type', full_name='api.ptransit.PTService.vehicle_type', index=13,
      number=14, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=95,
  serialized_end=514,
)


_PTGTFS = _descriptor.Descriptor(
  name='PTgtfs',
  full_name='api.ptransit.PTgtfs',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='gtfs_name', full_name='api.ptransit.PTgtfs.gtfs_name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='gtfs', full_name='api.ptransit.PTgtfs.gtfs', index=1,
      number=2, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=_b(""),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=516,
  serialized_end=557,
)

_PTSERVICE.fields_by_name['current_location'].message_type = common_dot_common__pb2._PLACE
_PTSERVICE.fields_by_name['next_arraival_timetable_time'].message_type = common_dot_common__pb2._TIME
_PTSERVICE.fields_by_name['past_departure_time'].message_type = common_dot_common__pb2._TIME
_PTSERVICE.fields_by_name['delay_time'].message_type = google_dot_protobuf_dot_duration__pb2._DURATION
DESCRIPTOR.message_types_by_name['PTService'] = _PTSERVICE
DESCRIPTOR.message_types_by_name['PTgtfs'] = _PTGTFS
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

PTService = _reflection.GeneratedProtocolMessageType('PTService', (_message.Message,), {
  'DESCRIPTOR' : _PTSERVICE,
  '__module__' : 'ptransit.ptransit_pb2'
  # @@protoc_insertion_point(class_scope:api.ptransit.PTService)
  })
_sym_db.RegisterMessage(PTService)

PTgtfs = _reflection.GeneratedProtocolMessageType('PTgtfs', (_message.Message,), {
  'DESCRIPTOR' : _PTGTFS,
  '__module__' : 'ptransit.ptransit_pb2'
  # @@protoc_insertion_point(class_scope:api.ptransit.PTgtfs)
  })
_sym_db.RegisterMessage(PTgtfs)


DESCRIPTOR._options = None
# @@protoc_insertion_point(module_scope)