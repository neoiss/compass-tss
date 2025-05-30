// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: types/query_constant_values.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type QueryConstantValuesRequest struct {
	Height string `protobuf:"bytes,1,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *QueryConstantValuesRequest) Reset()         { *m = QueryConstantValuesRequest{} }
func (m *QueryConstantValuesRequest) String() string { return proto.CompactTextString(m) }
func (*QueryConstantValuesRequest) ProtoMessage()    {}
func (*QueryConstantValuesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_09a11eb66540db6e, []int{0}
}
func (m *QueryConstantValuesRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryConstantValuesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryConstantValuesRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryConstantValuesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryConstantValuesRequest.Merge(m, src)
}
func (m *QueryConstantValuesRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryConstantValuesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryConstantValuesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryConstantValuesRequest proto.InternalMessageInfo

func (m *QueryConstantValuesRequest) GetHeight() string {
	if m != nil {
		return m.Height
	}
	return ""
}

type QueryConstantValuesResponse struct {
	Int_64Values []*Int64Constants  `protobuf:"bytes,1,rep,name=int_64_values,json=int64Values,proto3" json:"int_64_values,omitempty"`
	BoolValues   []*BoolConstants   `protobuf:"bytes,2,rep,name=bool_values,json=boolValues,proto3" json:"bool_values,omitempty"`
	StringValues []*StringConstants `protobuf:"bytes,3,rep,name=string_values,json=stringValues,proto3" json:"string_values,omitempty"`
}

func (m *QueryConstantValuesResponse) Reset()         { *m = QueryConstantValuesResponse{} }
func (m *QueryConstantValuesResponse) String() string { return proto.CompactTextString(m) }
func (*QueryConstantValuesResponse) ProtoMessage()    {}
func (*QueryConstantValuesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_09a11eb66540db6e, []int{1}
}
func (m *QueryConstantValuesResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryConstantValuesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryConstantValuesResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryConstantValuesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryConstantValuesResponse.Merge(m, src)
}
func (m *QueryConstantValuesResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryConstantValuesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryConstantValuesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryConstantValuesResponse proto.InternalMessageInfo

func (m *QueryConstantValuesResponse) GetInt_64Values() []*Int64Constants {
	if m != nil {
		return m.Int_64Values
	}
	return nil
}

func (m *QueryConstantValuesResponse) GetBoolValues() []*BoolConstants {
	if m != nil {
		return m.BoolValues
	}
	return nil
}

func (m *QueryConstantValuesResponse) GetStringValues() []*StringConstants {
	if m != nil {
		return m.StringValues
	}
	return nil
}

type Int64Constants struct {
	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value int64  `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *Int64Constants) Reset()         { *m = Int64Constants{} }
func (m *Int64Constants) String() string { return proto.CompactTextString(m) }
func (*Int64Constants) ProtoMessage()    {}
func (*Int64Constants) Descriptor() ([]byte, []int) {
	return fileDescriptor_09a11eb66540db6e, []int{2}
}
func (m *Int64Constants) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Int64Constants) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Int64Constants.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Int64Constants) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Int64Constants.Merge(m, src)
}
func (m *Int64Constants) XXX_Size() int {
	return m.Size()
}
func (m *Int64Constants) XXX_DiscardUnknown() {
	xxx_messageInfo_Int64Constants.DiscardUnknown(m)
}

var xxx_messageInfo_Int64Constants proto.InternalMessageInfo

func (m *Int64Constants) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Int64Constants) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type BoolConstants struct {
	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value bool   `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *BoolConstants) Reset()         { *m = BoolConstants{} }
func (m *BoolConstants) String() string { return proto.CompactTextString(m) }
func (*BoolConstants) ProtoMessage()    {}
func (*BoolConstants) Descriptor() ([]byte, []int) {
	return fileDescriptor_09a11eb66540db6e, []int{3}
}
func (m *BoolConstants) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BoolConstants) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BoolConstants.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BoolConstants) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BoolConstants.Merge(m, src)
}
func (m *BoolConstants) XXX_Size() int {
	return m.Size()
}
func (m *BoolConstants) XXX_DiscardUnknown() {
	xxx_messageInfo_BoolConstants.DiscardUnknown(m)
}

var xxx_messageInfo_BoolConstants proto.InternalMessageInfo

func (m *BoolConstants) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *BoolConstants) GetValue() bool {
	if m != nil {
		return m.Value
	}
	return false
}

type StringConstants struct {
	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *StringConstants) Reset()         { *m = StringConstants{} }
func (m *StringConstants) String() string { return proto.CompactTextString(m) }
func (*StringConstants) ProtoMessage()    {}
func (*StringConstants) Descriptor() ([]byte, []int) {
	return fileDescriptor_09a11eb66540db6e, []int{4}
}
func (m *StringConstants) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *StringConstants) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StringConstants.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *StringConstants) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StringConstants.Merge(m, src)
}
func (m *StringConstants) XXX_Size() int {
	return m.Size()
}
func (m *StringConstants) XXX_DiscardUnknown() {
	xxx_messageInfo_StringConstants.DiscardUnknown(m)
}

var xxx_messageInfo_StringConstants proto.InternalMessageInfo

func (m *StringConstants) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *StringConstants) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*QueryConstantValuesRequest)(nil), "types.QueryConstantValuesRequest")
	proto.RegisterType((*QueryConstantValuesResponse)(nil), "types.QueryConstantValuesResponse")
	proto.RegisterType((*Int64Constants)(nil), "types.Int64Constants")
	proto.RegisterType((*BoolConstants)(nil), "types.BoolConstants")
	proto.RegisterType((*StringConstants)(nil), "types.StringConstants")
}

func init() { proto.RegisterFile("types/query_constant_values.proto", fileDescriptor_09a11eb66540db6e) }

var fileDescriptor_09a11eb66540db6e = []byte{
	// 337 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0xd2, 0x41, 0x4b, 0x02, 0x41,
	0x14, 0x07, 0x70, 0x47, 0x53, 0xf2, 0x99, 0x05, 0x83, 0x89, 0x18, 0x2c, 0xb6, 0x27, 0x4f, 0xbb,
	0xa0, 0x26, 0x98, 0x37, 0x3b, 0x05, 0x5d, 0xda, 0xa0, 0x43, 0x17, 0x59, 0x6d, 0xd8, 0x1d, 0x58,
	0xe7, 0xe9, 0xce, 0x28, 0xf9, 0x2d, 0xfa, 0x58, 0x1d, 0x3a, 0x78, 0xec, 0x18, 0xfa, 0x45, 0x62,
	0x67, 0x67, 0x31, 0x23, 0xf0, 0xf6, 0xde, 0xcc, 0xfb, 0x3d, 0xfe, 0xbb, 0x0c, 0x5c, 0xab, 0xf5,
	0x9c, 0x49, 0x77, 0xb1, 0x64, 0xf1, 0x7a, 0x3c, 0x45, 0x21, 0x95, 0x2f, 0xd4, 0x78, 0xe5, 0x47,
	0x4b, 0x26, 0x9d, 0x79, 0x8c, 0x0a, 0x69, 0x51, 0x8f, 0x34, 0x6b, 0x01, 0x06, 0xa8, 0x4f, 0xdc,
	0xa4, 0x4a, 0x2f, 0xed, 0x1e, 0x34, 0x1f, 0x13, 0x7b, 0x67, 0xe8, 0xb3, 0x96, 0x1e, 0x5b, 0x2c,
	0x99, 0x54, 0xb4, 0x0e, 0xa5, 0x90, 0xf1, 0x20, 0x54, 0x0d, 0xd2, 0x22, 0xed, 0xb2, 0x67, 0x3a,
	0xfb, 0x93, 0xc0, 0xd5, 0xbf, 0x4c, 0xce, 0x51, 0x48, 0x46, 0x07, 0x50, 0xe5, 0x42, 0x8d, 0xfb,
	0x3d, 0x93, 0xa4, 0x41, 0x5a, 0x85, 0x76, 0xa5, 0x73, 0xe9, 0xe8, 0x28, 0xce, 0xbd, 0x50, 0xfd,
	0x5e, 0x46, 0xa5, 0x57, 0xe1, 0x49, 0x9f, 0xae, 0xa0, 0x37, 0x50, 0x99, 0x20, 0x46, 0x19, 0xcc,
	0x6b, 0x58, 0x33, 0x70, 0x84, 0x18, 0xed, 0x1d, 0x24, 0x83, 0x86, 0x0d, 0xa1, 0x2a, 0x55, 0xcc,
	0x45, 0x90, 0xc1, 0x82, 0x86, 0x75, 0x03, 0x9f, 0xf4, 0xdd, 0x9e, 0x9e, 0xa5, 0xc3, 0x29, 0xb6,
	0x6f, 0xe1, 0xfc, 0x30, 0x12, 0xa5, 0x70, 0x22, 0xfc, 0x19, 0x33, 0x9f, 0xad, 0x6b, 0x5a, 0x83,
	0xa2, 0xde, 0xdd, 0xc8, 0xb7, 0x48, 0xbb, 0xe0, 0xa5, 0x8d, 0x3d, 0x80, 0xea, 0x41, 0xaa, 0xe3,
	0xf4, 0x34, 0xa3, 0x43, 0xb8, 0xf8, 0x93, 0xeb, 0x38, 0x2e, 0x1b, 0x3c, 0x7a, 0xf8, 0xd8, 0x5a,
	0x64, 0xb3, 0xb5, 0xc8, 0xf7, 0xd6, 0x22, 0xef, 0x3b, 0x2b, 0xb7, 0xd9, 0x59, 0xb9, 0xaf, 0x9d,
	0x95, 0x7b, 0xe9, 0x04, 0x5c, 0x45, 0xfe, 0xc4, 0x99, 0xe2, 0xcc, 0x55, 0x21, 0xc6, 0xd3, 0xd0,
	0xe7, 0x42, 0x57, 0x02, 0x5f, 0x99, 0xbb, 0xea, 0xba, 0x6f, 0xbf, 0xcf, 0x93, 0xff, 0x33, 0x29,
	0xe9, 0xd7, 0xd0, 0xfd, 0x09, 0x00, 0x00, 0xff, 0xff, 0xde, 0xf6, 0x20, 0xfd, 0x4f, 0x02, 0x00,
	0x00,
}

func (m *QueryConstantValuesRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryConstantValuesRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryConstantValuesRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Height) > 0 {
		i -= len(m.Height)
		copy(dAtA[i:], m.Height)
		i = encodeVarintQueryConstantValues(dAtA, i, uint64(len(m.Height)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryConstantValuesResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryConstantValuesResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryConstantValuesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.StringValues) > 0 {
		for iNdEx := len(m.StringValues) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.StringValues[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQueryConstantValues(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.BoolValues) > 0 {
		for iNdEx := len(m.BoolValues) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.BoolValues[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQueryConstantValues(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Int_64Values) > 0 {
		for iNdEx := len(m.Int_64Values) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Int_64Values[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQueryConstantValues(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Int64Constants) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Int64Constants) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Int64Constants) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Value != 0 {
		i = encodeVarintQueryConstantValues(dAtA, i, uint64(m.Value))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintQueryConstantValues(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *BoolConstants) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BoolConstants) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BoolConstants) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Value {
		i--
		if m.Value {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintQueryConstantValues(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *StringConstants) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StringConstants) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *StringConstants) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Value) > 0 {
		i -= len(m.Value)
		copy(dAtA[i:], m.Value)
		i = encodeVarintQueryConstantValues(dAtA, i, uint64(len(m.Value)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintQueryConstantValues(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQueryConstantValues(dAtA []byte, offset int, v uint64) int {
	offset -= sovQueryConstantValues(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryConstantValuesRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Height)
	if l > 0 {
		n += 1 + l + sovQueryConstantValues(uint64(l))
	}
	return n
}

func (m *QueryConstantValuesResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Int_64Values) > 0 {
		for _, e := range m.Int_64Values {
			l = e.Size()
			n += 1 + l + sovQueryConstantValues(uint64(l))
		}
	}
	if len(m.BoolValues) > 0 {
		for _, e := range m.BoolValues {
			l = e.Size()
			n += 1 + l + sovQueryConstantValues(uint64(l))
		}
	}
	if len(m.StringValues) > 0 {
		for _, e := range m.StringValues {
			l = e.Size()
			n += 1 + l + sovQueryConstantValues(uint64(l))
		}
	}
	return n
}

func (m *Int64Constants) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovQueryConstantValues(uint64(l))
	}
	if m.Value != 0 {
		n += 1 + sovQueryConstantValues(uint64(m.Value))
	}
	return n
}

func (m *BoolConstants) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovQueryConstantValues(uint64(l))
	}
	if m.Value {
		n += 2
	}
	return n
}

func (m *StringConstants) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovQueryConstantValues(uint64(l))
	}
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovQueryConstantValues(uint64(l))
	}
	return n
}

func sovQueryConstantValues(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQueryConstantValues(x uint64) (n int) {
	return sovQueryConstantValues(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryConstantValuesRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQueryConstantValues
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryConstantValuesRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryConstantValuesRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryConstantValues
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Height = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQueryConstantValues(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryConstantValuesResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQueryConstantValues
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryConstantValuesResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryConstantValuesResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Int_64Values", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryConstantValues
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Int_64Values = append(m.Int_64Values, &Int64Constants{})
			if err := m.Int_64Values[len(m.Int_64Values)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BoolValues", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryConstantValues
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BoolValues = append(m.BoolValues, &BoolConstants{})
			if err := m.BoolValues[len(m.BoolValues)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StringValues", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryConstantValues
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StringValues = append(m.StringValues, &StringConstants{})
			if err := m.StringValues[len(m.StringValues)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQueryConstantValues(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Int64Constants) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQueryConstantValues
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Int64Constants: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Int64Constants: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryConstantValues
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			m.Value = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryConstantValues
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Value |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQueryConstantValues(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BoolConstants) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQueryConstantValues
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BoolConstants: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BoolConstants: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryConstantValues
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryConstantValues
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Value = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipQueryConstantValues(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *StringConstants) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQueryConstantValues
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: StringConstants: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StringConstants: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryConstantValues
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryConstantValues
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQueryConstantValues(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQueryConstantValues
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQueryConstantValues(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQueryConstantValues
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQueryConstantValues
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQueryConstantValues
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQueryConstantValues
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQueryConstantValues
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQueryConstantValues
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQueryConstantValues        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQueryConstantValues          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQueryConstantValues = fmt.Errorf("proto: unexpected end of group")
)
