# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: api/smarket.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2
from google.protobuf import duration_pb2 as google_dot_protobuf_dot_duration__pb2
from fleet import fleet_pb2 as fleet_dot_fleet__pb2
from rideshare import rideshare_pb2 as rideshare_dot_rideshare__pb2
from adservice import adservice_pb2 as adservice_dot_adservice__pb2
from library import library_pb2 as library_dot_library__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='api/smarket.proto',
  package='api',
  syntax='proto3',
  serialized_options=_b('Z\003api'),
  serialized_pb=_b('\n\x11\x61pi/smarket.proto\x12\x03\x61pi\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x1egoogle/protobuf/duration.proto\x1a\x11\x66leet/fleet.proto\x1a\x19rideshare/rideshare.proto\x1a\x19\x61\x64service/adservice.proto\x1a\x15library/library.proto\"#\n\x08Response\x12\n\n\x02ok\x18\x01 \x01(\x08\x12\x0b\n\x03\x65rr\x18\x02 \x01(\t\"S\n\x0f\x43onfirmResponse\x12\n\n\x02ok\x18\x01 \x01(\x08\x12\'\n\x04wait\x18\x02 \x01(\x0b\x32\x19.google.protobuf.Duration\x12\x0b\n\x03\x65rr\x18\x03 \x01(\t\"\xf5\x02\n\x06Supply\x12\n\n\x02id\x18\x01 \x01(\x06\x12\x11\n\tsender_id\x18\x02 \x01(\x06\x12\x11\n\ttarget_id\x18\x03 \x01(\x06\x12\x1d\n\x04type\x18\x04 \x01(\x0e\x32\x0f.api.MarketType\x12\x13\n\x0bsupply_name\x18\x05 \x01(\t\x12&\n\x02ts\x18\x06 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\x10\n\x08\x61rg_json\x18\x07 \x01(\t\x12%\n\targ_Fleet\x18\n \x01(\x0b\x32\x10.api.fleet.FleetH\x00\x12\x31\n\rarg_RideShare\x18\x0b \x01(\x0b\x32\x18.api.rideshare.RideShareH\x00\x12\x31\n\rarg_AdService\x18\x0c \x01(\x0b\x32\x18.api.adservice.AdServiceH\x00\x12\x31\n\x0e\x61rg_LibService\x18\r \x01(\x0b\x32\x17.api.library.LibServiceH\x00\x42\x0b\n\targ_oneof\"\xf5\x02\n\x06\x44\x65mand\x12\n\n\x02id\x18\x01 \x01(\x06\x12\x11\n\tsender_id\x18\x02 \x01(\x06\x12\x11\n\ttarget_id\x18\x03 \x01(\x06\x12\x1d\n\x04type\x18\x04 \x01(\x0e\x32\x0f.api.MarketType\x12\x13\n\x0b\x64\x65mand_name\x18\x05 \x01(\t\x12&\n\x02ts\x18\x06 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\x10\n\x08\x61rg_json\x18\x07 \x01(\t\x12%\n\targ_Fleet\x18\n \x01(\x0b\x32\x10.api.fleet.FleetH\x00\x12\x31\n\rarg_RideShare\x18\x0b \x01(\x0b\x32\x18.api.rideshare.RideShareH\x00\x12\x31\n\rarg_AdService\x18\x0c \x01(\x0b\x32\x18.api.adservice.AdServiceH\x00\x12\x31\n\x0e\x61rg_LibService\x18\r \x01(\x0b\x32\x17.api.library.LibServiceH\x00\x42\x0b\n\targ_oneof\"\x82\x01\n\x06Target\x12\n\n\x02id\x18\x01 \x01(\x06\x12\x11\n\tsender_id\x18\x02 \x01(\x06\x12\x11\n\ttarget_id\x18\x03 \x01(\x06\x12\x1d\n\x04type\x18\x04 \x01(\x0e\x32\x0f.api.MarketType\x12\'\n\x04wait\x18\x05 \x01(\x0b\x32\x19.google.protobuf.Duration\"M\n\x07\x43hannel\x12\x11\n\tclient_id\x18\x01 \x01(\x06\x12\x1d\n\x04type\x18\x02 \x01(\x0e\x32\x0f.api.MarketType\x12\x10\n\x08\x61rg_json\x18\x03 \x01(\t*P\n\nMarketType\x12\x08\n\x04NONE\x10\x00\x12\x0e\n\nRIDE_SHARE\x10\x01\x12\x0e\n\nAD_SERVICE\x10\x02\x12\x0f\n\x0bLIB_SERVICE\x10\x03\x12\x07\n\x03\x45ND\x10\n2\xaa\x04\n\x07SMarket\x12.\n\x0eRegisterDemand\x12\x0b.api.Demand\x1a\r.api.Response\"\x00\x12.\n\x0eRegisterSupply\x12\x0b.api.Supply\x1a\r.api.Response\"\x00\x12-\n\rProposeDemand\x12\x0b.api.Demand\x1a\r.api.Response\"\x00\x12-\n\rProposeSupply\x12\x0b.api.Supply\x1a\r.api.Response\"\x00\x12\x34\n\rReserveSupply\x12\x0b.api.Target\x1a\x14.api.ConfirmResponse\"\x00\x12\x34\n\rReserveDemand\x12\x0b.api.Target\x1a\x14.api.ConfirmResponse\"\x00\x12\x33\n\x0cSelectSupply\x12\x0b.api.Target\x1a\x14.api.ConfirmResponse\"\x00\x12\x33\n\x0cSelectDemand\x12\x0b.api.Target\x1a\x14.api.ConfirmResponse\"\x00\x12\'\n\x07\x43onfirm\x12\x0b.api.Target\x1a\r.api.Response\"\x00\x12\x30\n\x0fSubscribeDemand\x12\x0c.api.Channel\x1a\x0b.api.Demand\"\x00\x30\x01\x12\x30\n\x0fSubscribeSupply\x12\x0c.api.Channel\x1a\x0b.api.Supply\"\x00\x30\x01\x42\x05Z\x03\x61pib\x06proto3')
  ,
  dependencies=[google_dot_protobuf_dot_timestamp__pb2.DESCRIPTOR,google_dot_protobuf_dot_duration__pb2.DESCRIPTOR,fleet_dot_fleet__pb2.DESCRIPTOR,rideshare_dot_rideshare__pb2.DESCRIPTOR,adservice_dot_adservice__pb2.DESCRIPTOR,library_dot_library__pb2.DESCRIPTOR,])

_MARKETTYPE = _descriptor.EnumDescriptor(
  name='MarketType',
  full_name='api.MarketType',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='NONE', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='RIDE_SHARE', index=1, number=1,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='AD_SERVICE', index=2, number=2,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='LIB_SERVICE', index=3, number=3,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='END', index=4, number=10,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=1273,
  serialized_end=1353,
)
_sym_db.RegisterEnumDescriptor(_MARKETTYPE)

MarketType = enum_type_wrapper.EnumTypeWrapper(_MARKETTYPE)
NONE = 0
RIDE_SHARE = 1
AD_SERVICE = 2
LIB_SERVICE = 3
END = 10



_RESPONSE = _descriptor.Descriptor(
  name='Response',
  full_name='api.Response',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='ok', full_name='api.Response.ok', index=0,
      number=1, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='err', full_name='api.Response.err', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
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
  serialized_start=187,
  serialized_end=222,
)


_CONFIRMRESPONSE = _descriptor.Descriptor(
  name='ConfirmResponse',
  full_name='api.ConfirmResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='ok', full_name='api.ConfirmResponse.ok', index=0,
      number=1, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='wait', full_name='api.ConfirmResponse.wait', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='err', full_name='api.ConfirmResponse.err', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
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
  serialized_start=224,
  serialized_end=307,
)


_SUPPLY = _descriptor.Descriptor(
  name='Supply',
  full_name='api.Supply',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='api.Supply.id', index=0,
      number=1, type=6, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='sender_id', full_name='api.Supply.sender_id', index=1,
      number=2, type=6, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='target_id', full_name='api.Supply.target_id', index=2,
      number=3, type=6, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='type', full_name='api.Supply.type', index=3,
      number=4, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='supply_name', full_name='api.Supply.supply_name', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='ts', full_name='api.Supply.ts', index=5,
      number=6, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='arg_json', full_name='api.Supply.arg_json', index=6,
      number=7, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='arg_Fleet', full_name='api.Supply.arg_Fleet', index=7,
      number=10, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='arg_RideShare', full_name='api.Supply.arg_RideShare', index=8,
      number=11, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='arg_AdService', full_name='api.Supply.arg_AdService', index=9,
      number=12, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='arg_LibService', full_name='api.Supply.arg_LibService', index=10,
      number=13, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
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
    _descriptor.OneofDescriptor(
      name='arg_oneof', full_name='api.Supply.arg_oneof',
      index=0, containing_type=None, fields=[]),
  ],
  serialized_start=310,
  serialized_end=683,
)


_DEMAND = _descriptor.Descriptor(
  name='Demand',
  full_name='api.Demand',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='api.Demand.id', index=0,
      number=1, type=6, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='sender_id', full_name='api.Demand.sender_id', index=1,
      number=2, type=6, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='target_id', full_name='api.Demand.target_id', index=2,
      number=3, type=6, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='type', full_name='api.Demand.type', index=3,
      number=4, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='demand_name', full_name='api.Demand.demand_name', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='ts', full_name='api.Demand.ts', index=5,
      number=6, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='arg_json', full_name='api.Demand.arg_json', index=6,
      number=7, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='arg_Fleet', full_name='api.Demand.arg_Fleet', index=7,
      number=10, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='arg_RideShare', full_name='api.Demand.arg_RideShare', index=8,
      number=11, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='arg_AdService', full_name='api.Demand.arg_AdService', index=9,
      number=12, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='arg_LibService', full_name='api.Demand.arg_LibService', index=10,
      number=13, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
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
    _descriptor.OneofDescriptor(
      name='arg_oneof', full_name='api.Demand.arg_oneof',
      index=0, containing_type=None, fields=[]),
  ],
  serialized_start=686,
  serialized_end=1059,
)


_TARGET = _descriptor.Descriptor(
  name='Target',
  full_name='api.Target',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='api.Target.id', index=0,
      number=1, type=6, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='sender_id', full_name='api.Target.sender_id', index=1,
      number=2, type=6, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='target_id', full_name='api.Target.target_id', index=2,
      number=3, type=6, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='type', full_name='api.Target.type', index=3,
      number=4, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='wait', full_name='api.Target.wait', index=4,
      number=5, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
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
  serialized_start=1062,
  serialized_end=1192,
)


_CHANNEL = _descriptor.Descriptor(
  name='Channel',
  full_name='api.Channel',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='client_id', full_name='api.Channel.client_id', index=0,
      number=1, type=6, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='type', full_name='api.Channel.type', index=1,
      number=2, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='arg_json', full_name='api.Channel.arg_json', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
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
  serialized_start=1194,
  serialized_end=1271,
)

_CONFIRMRESPONSE.fields_by_name['wait'].message_type = google_dot_protobuf_dot_duration__pb2._DURATION
_SUPPLY.fields_by_name['type'].enum_type = _MARKETTYPE
_SUPPLY.fields_by_name['ts'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_SUPPLY.fields_by_name['arg_Fleet'].message_type = fleet_dot_fleet__pb2._FLEET
_SUPPLY.fields_by_name['arg_RideShare'].message_type = rideshare_dot_rideshare__pb2._RIDESHARE
_SUPPLY.fields_by_name['arg_AdService'].message_type = adservice_dot_adservice__pb2._ADSERVICE
_SUPPLY.fields_by_name['arg_LibService'].message_type = library_dot_library__pb2._LIBSERVICE
_SUPPLY.oneofs_by_name['arg_oneof'].fields.append(
  _SUPPLY.fields_by_name['arg_Fleet'])
_SUPPLY.fields_by_name['arg_Fleet'].containing_oneof = _SUPPLY.oneofs_by_name['arg_oneof']
_SUPPLY.oneofs_by_name['arg_oneof'].fields.append(
  _SUPPLY.fields_by_name['arg_RideShare'])
_SUPPLY.fields_by_name['arg_RideShare'].containing_oneof = _SUPPLY.oneofs_by_name['arg_oneof']
_SUPPLY.oneofs_by_name['arg_oneof'].fields.append(
  _SUPPLY.fields_by_name['arg_AdService'])
_SUPPLY.fields_by_name['arg_AdService'].containing_oneof = _SUPPLY.oneofs_by_name['arg_oneof']
_SUPPLY.oneofs_by_name['arg_oneof'].fields.append(
  _SUPPLY.fields_by_name['arg_LibService'])
_SUPPLY.fields_by_name['arg_LibService'].containing_oneof = _SUPPLY.oneofs_by_name['arg_oneof']
_DEMAND.fields_by_name['type'].enum_type = _MARKETTYPE
_DEMAND.fields_by_name['ts'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_DEMAND.fields_by_name['arg_Fleet'].message_type = fleet_dot_fleet__pb2._FLEET
_DEMAND.fields_by_name['arg_RideShare'].message_type = rideshare_dot_rideshare__pb2._RIDESHARE
_DEMAND.fields_by_name['arg_AdService'].message_type = adservice_dot_adservice__pb2._ADSERVICE
_DEMAND.fields_by_name['arg_LibService'].message_type = library_dot_library__pb2._LIBSERVICE
_DEMAND.oneofs_by_name['arg_oneof'].fields.append(
  _DEMAND.fields_by_name['arg_Fleet'])
_DEMAND.fields_by_name['arg_Fleet'].containing_oneof = _DEMAND.oneofs_by_name['arg_oneof']
_DEMAND.oneofs_by_name['arg_oneof'].fields.append(
  _DEMAND.fields_by_name['arg_RideShare'])
_DEMAND.fields_by_name['arg_RideShare'].containing_oneof = _DEMAND.oneofs_by_name['arg_oneof']
_DEMAND.oneofs_by_name['arg_oneof'].fields.append(
  _DEMAND.fields_by_name['arg_AdService'])
_DEMAND.fields_by_name['arg_AdService'].containing_oneof = _DEMAND.oneofs_by_name['arg_oneof']
_DEMAND.oneofs_by_name['arg_oneof'].fields.append(
  _DEMAND.fields_by_name['arg_LibService'])
_DEMAND.fields_by_name['arg_LibService'].containing_oneof = _DEMAND.oneofs_by_name['arg_oneof']
_TARGET.fields_by_name['type'].enum_type = _MARKETTYPE
_TARGET.fields_by_name['wait'].message_type = google_dot_protobuf_dot_duration__pb2._DURATION
_CHANNEL.fields_by_name['type'].enum_type = _MARKETTYPE
DESCRIPTOR.message_types_by_name['Response'] = _RESPONSE
DESCRIPTOR.message_types_by_name['ConfirmResponse'] = _CONFIRMRESPONSE
DESCRIPTOR.message_types_by_name['Supply'] = _SUPPLY
DESCRIPTOR.message_types_by_name['Demand'] = _DEMAND
DESCRIPTOR.message_types_by_name['Target'] = _TARGET
DESCRIPTOR.message_types_by_name['Channel'] = _CHANNEL
DESCRIPTOR.enum_types_by_name['MarketType'] = _MARKETTYPE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Response = _reflection.GeneratedProtocolMessageType('Response', (_message.Message,), dict(
  DESCRIPTOR = _RESPONSE,
  __module__ = 'api.smarket_pb2'
  # @@protoc_insertion_point(class_scope:api.Response)
  ))
_sym_db.RegisterMessage(Response)

ConfirmResponse = _reflection.GeneratedProtocolMessageType('ConfirmResponse', (_message.Message,), dict(
  DESCRIPTOR = _CONFIRMRESPONSE,
  __module__ = 'api.smarket_pb2'
  # @@protoc_insertion_point(class_scope:api.ConfirmResponse)
  ))
_sym_db.RegisterMessage(ConfirmResponse)

Supply = _reflection.GeneratedProtocolMessageType('Supply', (_message.Message,), dict(
  DESCRIPTOR = _SUPPLY,
  __module__ = 'api.smarket_pb2'
  # @@protoc_insertion_point(class_scope:api.Supply)
  ))
_sym_db.RegisterMessage(Supply)

Demand = _reflection.GeneratedProtocolMessageType('Demand', (_message.Message,), dict(
  DESCRIPTOR = _DEMAND,
  __module__ = 'api.smarket_pb2'
  # @@protoc_insertion_point(class_scope:api.Demand)
  ))
_sym_db.RegisterMessage(Demand)

Target = _reflection.GeneratedProtocolMessageType('Target', (_message.Message,), dict(
  DESCRIPTOR = _TARGET,
  __module__ = 'api.smarket_pb2'
  # @@protoc_insertion_point(class_scope:api.Target)
  ))
_sym_db.RegisterMessage(Target)

Channel = _reflection.GeneratedProtocolMessageType('Channel', (_message.Message,), dict(
  DESCRIPTOR = _CHANNEL,
  __module__ = 'api.smarket_pb2'
  # @@protoc_insertion_point(class_scope:api.Channel)
  ))
_sym_db.RegisterMessage(Channel)


DESCRIPTOR._options = None

_SMARKET = _descriptor.ServiceDescriptor(
  name='SMarket',
  full_name='api.SMarket',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=1356,
  serialized_end=1910,
  methods=[
  _descriptor.MethodDescriptor(
    name='RegisterDemand',
    full_name='api.SMarket.RegisterDemand',
    index=0,
    containing_service=None,
    input_type=_DEMAND,
    output_type=_RESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='RegisterSupply',
    full_name='api.SMarket.RegisterSupply',
    index=1,
    containing_service=None,
    input_type=_SUPPLY,
    output_type=_RESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='ProposeDemand',
    full_name='api.SMarket.ProposeDemand',
    index=2,
    containing_service=None,
    input_type=_DEMAND,
    output_type=_RESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='ProposeSupply',
    full_name='api.SMarket.ProposeSupply',
    index=3,
    containing_service=None,
    input_type=_SUPPLY,
    output_type=_RESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='ReserveSupply',
    full_name='api.SMarket.ReserveSupply',
    index=4,
    containing_service=None,
    input_type=_TARGET,
    output_type=_CONFIRMRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='ReserveDemand',
    full_name='api.SMarket.ReserveDemand',
    index=5,
    containing_service=None,
    input_type=_TARGET,
    output_type=_CONFIRMRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='SelectSupply',
    full_name='api.SMarket.SelectSupply',
    index=6,
    containing_service=None,
    input_type=_TARGET,
    output_type=_CONFIRMRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='SelectDemand',
    full_name='api.SMarket.SelectDemand',
    index=7,
    containing_service=None,
    input_type=_TARGET,
    output_type=_CONFIRMRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='Confirm',
    full_name='api.SMarket.Confirm',
    index=8,
    containing_service=None,
    input_type=_TARGET,
    output_type=_RESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='SubscribeDemand',
    full_name='api.SMarket.SubscribeDemand',
    index=9,
    containing_service=None,
    input_type=_CHANNEL,
    output_type=_DEMAND,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='SubscribeSupply',
    full_name='api.SMarket.SubscribeSupply',
    index=10,
    containing_service=None,
    input_type=_CHANNEL,
    output_type=_SUPPLY,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_SMARKET)

DESCRIPTOR.services_by_name['SMarket'] = _SMARKET

# @@protoc_insertion_point(module_scope)