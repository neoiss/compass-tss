// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: types/query_outbound_fee.proto

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

type QueryOutboundFeeRequest struct {
	Asset  string `protobuf:"bytes,1,opt,name=asset,proto3" json:"asset,omitempty"`
	Height string `protobuf:"bytes,2,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *QueryOutboundFeeRequest) Reset()         { *m = QueryOutboundFeeRequest{} }
func (m *QueryOutboundFeeRequest) String() string { return proto.CompactTextString(m) }
func (*QueryOutboundFeeRequest) ProtoMessage()    {}
func (*QueryOutboundFeeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec7cb76bf82f4d2e, []int{0}
}
func (m *QueryOutboundFeeRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryOutboundFeeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryOutboundFeeRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryOutboundFeeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryOutboundFeeRequest.Merge(m, src)
}
func (m *QueryOutboundFeeRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryOutboundFeeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryOutboundFeeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryOutboundFeeRequest proto.InternalMessageInfo

func (m *QueryOutboundFeeRequest) GetAsset() string {
	if m != nil {
		return m.Asset
	}
	return ""
}

func (m *QueryOutboundFeeRequest) GetHeight() string {
	if m != nil {
		return m.Height
	}
	return ""
}

type QueryOutboundFeeResponse struct {
	// the asset to display the outbound fee for
	Asset string `protobuf:"bytes,1,opt,name=asset,proto3" json:"asset"`
	// the asset's outbound fee, in (1e8-format) units of the asset
	OutboundFee string `protobuf:"bytes,2,opt,name=outbound_fee,json=outboundFee,proto3" json:"outbound_fee"`
	// Total RUNE the network has withheld as fees to later cover gas costs for this asset's outbounds
	FeeWithheldRune string `protobuf:"bytes,3,opt,name=fee_withheld_rune,json=feeWithheldRune,proto3" json:"fee_withheld_rune,omitempty"`
	// Total RUNE the network has spent to reimburse gas costs for this asset's outbounds
	FeeSpentRune string `protobuf:"bytes,4,opt,name=fee_spent_rune,json=feeSpentRune,proto3" json:"fee_spent_rune,omitempty"`
	// amount of RUNE by which the fee_withheld_rune exceeds the fee_spent_rune
	SurplusRune string `protobuf:"bytes,5,opt,name=surplus_rune,json=surplusRune,proto3" json:"surplus_rune,omitempty"`
	// dynamic multiplier basis points, based on the surplus_rune, affecting the size of the outbound_fee
	DynamicMultiplierBasisPoints string `protobuf:"bytes,6,opt,name=dynamic_multiplier_basis_points,json=dynamicMultiplierBasisPoints,proto3" json:"dynamic_multiplier_basis_points,omitempty"`
}

func (m *QueryOutboundFeeResponse) Reset()         { *m = QueryOutboundFeeResponse{} }
func (m *QueryOutboundFeeResponse) String() string { return proto.CompactTextString(m) }
func (*QueryOutboundFeeResponse) ProtoMessage()    {}
func (*QueryOutboundFeeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec7cb76bf82f4d2e, []int{1}
}
func (m *QueryOutboundFeeResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryOutboundFeeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryOutboundFeeResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryOutboundFeeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryOutboundFeeResponse.Merge(m, src)
}
func (m *QueryOutboundFeeResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryOutboundFeeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryOutboundFeeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryOutboundFeeResponse proto.InternalMessageInfo

func (m *QueryOutboundFeeResponse) GetAsset() string {
	if m != nil {
		return m.Asset
	}
	return ""
}

func (m *QueryOutboundFeeResponse) GetOutboundFee() string {
	if m != nil {
		return m.OutboundFee
	}
	return ""
}

func (m *QueryOutboundFeeResponse) GetFeeWithheldRune() string {
	if m != nil {
		return m.FeeWithheldRune
	}
	return ""
}

func (m *QueryOutboundFeeResponse) GetFeeSpentRune() string {
	if m != nil {
		return m.FeeSpentRune
	}
	return ""
}

func (m *QueryOutboundFeeResponse) GetSurplusRune() string {
	if m != nil {
		return m.SurplusRune
	}
	return ""
}

func (m *QueryOutboundFeeResponse) GetDynamicMultiplierBasisPoints() string {
	if m != nil {
		return m.DynamicMultiplierBasisPoints
	}
	return ""
}

type QueryOutboundFeesRequest struct {
	Height string `protobuf:"bytes,1,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *QueryOutboundFeesRequest) Reset()         { *m = QueryOutboundFeesRequest{} }
func (m *QueryOutboundFeesRequest) String() string { return proto.CompactTextString(m) }
func (*QueryOutboundFeesRequest) ProtoMessage()    {}
func (*QueryOutboundFeesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec7cb76bf82f4d2e, []int{2}
}
func (m *QueryOutboundFeesRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryOutboundFeesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryOutboundFeesRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryOutboundFeesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryOutboundFeesRequest.Merge(m, src)
}
func (m *QueryOutboundFeesRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryOutboundFeesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryOutboundFeesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryOutboundFeesRequest proto.InternalMessageInfo

func (m *QueryOutboundFeesRequest) GetHeight() string {
	if m != nil {
		return m.Height
	}
	return ""
}

type QueryOutboundFeesResponse struct {
	OutboundFees []*QueryOutboundFeeResponse `protobuf:"bytes,1,rep,name=outbound_fees,json=outboundFees,proto3" json:"outbound_fees,omitempty"`
}

func (m *QueryOutboundFeesResponse) Reset()         { *m = QueryOutboundFeesResponse{} }
func (m *QueryOutboundFeesResponse) String() string { return proto.CompactTextString(m) }
func (*QueryOutboundFeesResponse) ProtoMessage()    {}
func (*QueryOutboundFeesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec7cb76bf82f4d2e, []int{3}
}
func (m *QueryOutboundFeesResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryOutboundFeesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryOutboundFeesResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryOutboundFeesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryOutboundFeesResponse.Merge(m, src)
}
func (m *QueryOutboundFeesResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryOutboundFeesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryOutboundFeesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryOutboundFeesResponse proto.InternalMessageInfo

func (m *QueryOutboundFeesResponse) GetOutboundFees() []*QueryOutboundFeeResponse {
	if m != nil {
		return m.OutboundFees
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryOutboundFeeRequest)(nil), "types.QueryOutboundFeeRequest")
	proto.RegisterType((*QueryOutboundFeeResponse)(nil), "types.QueryOutboundFeeResponse")
	proto.RegisterType((*QueryOutboundFeesRequest)(nil), "types.QueryOutboundFeesRequest")
	proto.RegisterType((*QueryOutboundFeesResponse)(nil), "types.QueryOutboundFeesResponse")
}

func init() { proto.RegisterFile("types/query_outbound_fee.proto", fileDescriptor_ec7cb76bf82f4d2e) }

var fileDescriptor_ec7cb76bf82f4d2e = []byte{
	// 404 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xcf, 0x6e, 0xa2, 0x40,
	0x1c, 0xc7, 0x45, 0x57, 0x13, 0x47, 0xf6, 0x1f, 0x31, 0xbb, 0xec, 0x66, 0x03, 0xae, 0xd9, 0x83,
	0xd9, 0x03, 0x24, 0xfa, 0x06, 0x64, 0xff, 0x5c, 0x76, 0xd3, 0x96, 0x1e, 0x9a, 0xf4, 0x42, 0x40,
	0x7e, 0xc0, 0x24, 0x38, 0x83, 0xcc, 0x4c, 0x5b, 0xdf, 0xa2, 0x0f, 0xd2, 0x07, 0xe9, 0xd1, 0x63,
	0x4f, 0xa6, 0xd1, 0x9b, 0x4f, 0xd1, 0x30, 0x60, 0xc5, 0xd8, 0xde, 0x66, 0x3e, 0xf3, 0xf9, 0x7d,
	0x99, 0x7c, 0x07, 0x64, 0xf0, 0x45, 0x06, 0xcc, 0x9e, 0x0b, 0xc8, 0x17, 0x1e, 0x15, 0x3c, 0xa0,
	0x82, 0x84, 0x5e, 0x04, 0x60, 0x65, 0x39, 0xe5, 0x54, 0x6b, 0xcb, 0xf3, 0xaf, 0xfd, 0x98, 0xc6,
	0x54, 0x12, 0xbb, 0x58, 0x95, 0x87, 0xc3, 0xbf, 0xe8, 0xf3, 0x59, 0x31, 0x78, 0x52, 0xcd, 0xfd,
	0x01, 0x70, 0x61, 0x2e, 0x80, 0x71, 0xad, 0x8f, 0xda, 0x3e, 0x63, 0xc0, 0x75, 0x65, 0xa0, 0x8c,
	0xba, 0x6e, 0xb9, 0xd1, 0x3e, 0xa1, 0x4e, 0x02, 0x38, 0x4e, 0xb8, 0xde, 0x94, 0xb8, 0xda, 0x0d,
	0xef, 0x9a, 0x48, 0x3f, 0x4e, 0x62, 0x19, 0x25, 0x0c, 0x34, 0xf3, 0x20, 0xca, 0xe9, 0x6e, 0x57,
	0x66, 0x09, 0x76, 0xa9, 0x13, 0xa4, 0xd6, 0x6f, 0x5e, 0x66, 0x3b, 0x1f, 0xb6, 0x2b, 0xf3, 0x80,
	0xbb, 0x3d, 0xba, 0x4f, 0xd7, 0x7e, 0xa2, 0x8f, 0x11, 0x80, 0x77, 0x8d, 0x79, 0x92, 0x40, 0x1a,
	0x7a, 0xb9, 0x20, 0xa0, 0xb7, 0xe4, 0xad, 0xde, 0x47, 0x00, 0x17, 0x15, 0x77, 0x05, 0x01, 0xed,
	0x07, 0x7a, 0x57, 0xb8, 0x2c, 0x03, 0xc2, 0x4b, 0xf1, 0x8d, 0x14, 0xd5, 0x08, 0xe0, 0xbc, 0x80,
	0xd2, 0xfa, 0x8e, 0x54, 0x26, 0xf2, 0x2c, 0x15, 0xac, 0x74, 0xda, 0xd2, 0xe9, 0x55, 0x4c, 0x2a,
	0xbf, 0x91, 0x19, 0x2e, 0x88, 0x3f, 0xc3, 0x53, 0x6f, 0x26, 0x52, 0x8e, 0xb3, 0x14, 0x43, 0xee,
	0x05, 0x3e, 0xc3, 0xcc, 0xcb, 0x28, 0x26, 0x9c, 0xe9, 0x1d, 0x39, 0xf5, 0xad, 0xd2, 0xfe, 0x3f,
	0x5b, 0x4e, 0x21, 0x9d, 0x4a, 0x67, 0x38, 0x3e, 0x6e, 0x8b, 0xed, 0x8a, 0xdf, 0x57, 0xac, 0x1c,
	0x54, 0xec, 0xa3, 0x2f, 0x2f, 0xcc, 0x54, 0x15, 0xff, 0x42, 0x6f, 0xeb, 0x4d, 0x31, 0x5d, 0x19,
	0xb4, 0x46, 0xbd, 0xb1, 0x69, 0xc9, 0xd7, 0xb7, 0x5e, 0x7b, 0x1a, 0x57, 0xad, 0x35, 0xca, 0x9c,
	0x7f, 0xf7, 0x6b, 0x43, 0x59, 0xae, 0x0d, 0xe5, 0x71, 0x6d, 0x28, 0xb7, 0x1b, 0xa3, 0xb1, 0xdc,
	0x18, 0x8d, 0x87, 0x8d, 0xd1, 0xb8, 0x1c, 0xc7, 0x98, 0xa7, 0x7e, 0x60, 0x4d, 0xe9, 0xcc, 0xe6,
	0x09, 0xcd, 0xa7, 0x89, 0x8f, 0x89, 0x5c, 0x11, 0x1a, 0x82, 0x7d, 0x35, 0xb1, 0x6f, 0xea, 0xbc,
	0xf8, 0x68, 0xd0, 0x91, 0xff, 0xd8, 0xe4, 0x29, 0x00, 0x00, 0xff, 0xff, 0x75, 0xa3, 0xec, 0xfc,
	0xa2, 0x02, 0x00, 0x00,
}

func (m *QueryOutboundFeeRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryOutboundFeeRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryOutboundFeeRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Height) > 0 {
		i -= len(m.Height)
		copy(dAtA[i:], m.Height)
		i = encodeVarintQueryOutboundFee(dAtA, i, uint64(len(m.Height)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Asset) > 0 {
		i -= len(m.Asset)
		copy(dAtA[i:], m.Asset)
		i = encodeVarintQueryOutboundFee(dAtA, i, uint64(len(m.Asset)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryOutboundFeeResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryOutboundFeeResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryOutboundFeeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.DynamicMultiplierBasisPoints) > 0 {
		i -= len(m.DynamicMultiplierBasisPoints)
		copy(dAtA[i:], m.DynamicMultiplierBasisPoints)
		i = encodeVarintQueryOutboundFee(dAtA, i, uint64(len(m.DynamicMultiplierBasisPoints)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.SurplusRune) > 0 {
		i -= len(m.SurplusRune)
		copy(dAtA[i:], m.SurplusRune)
		i = encodeVarintQueryOutboundFee(dAtA, i, uint64(len(m.SurplusRune)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.FeeSpentRune) > 0 {
		i -= len(m.FeeSpentRune)
		copy(dAtA[i:], m.FeeSpentRune)
		i = encodeVarintQueryOutboundFee(dAtA, i, uint64(len(m.FeeSpentRune)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.FeeWithheldRune) > 0 {
		i -= len(m.FeeWithheldRune)
		copy(dAtA[i:], m.FeeWithheldRune)
		i = encodeVarintQueryOutboundFee(dAtA, i, uint64(len(m.FeeWithheldRune)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.OutboundFee) > 0 {
		i -= len(m.OutboundFee)
		copy(dAtA[i:], m.OutboundFee)
		i = encodeVarintQueryOutboundFee(dAtA, i, uint64(len(m.OutboundFee)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Asset) > 0 {
		i -= len(m.Asset)
		copy(dAtA[i:], m.Asset)
		i = encodeVarintQueryOutboundFee(dAtA, i, uint64(len(m.Asset)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryOutboundFeesRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryOutboundFeesRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryOutboundFeesRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Height) > 0 {
		i -= len(m.Height)
		copy(dAtA[i:], m.Height)
		i = encodeVarintQueryOutboundFee(dAtA, i, uint64(len(m.Height)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryOutboundFeesResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryOutboundFeesResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryOutboundFeesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.OutboundFees) > 0 {
		for iNdEx := len(m.OutboundFees) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.OutboundFees[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQueryOutboundFee(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintQueryOutboundFee(dAtA []byte, offset int, v uint64) int {
	offset -= sovQueryOutboundFee(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryOutboundFeeRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Asset)
	if l > 0 {
		n += 1 + l + sovQueryOutboundFee(uint64(l))
	}
	l = len(m.Height)
	if l > 0 {
		n += 1 + l + sovQueryOutboundFee(uint64(l))
	}
	return n
}

func (m *QueryOutboundFeeResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Asset)
	if l > 0 {
		n += 1 + l + sovQueryOutboundFee(uint64(l))
	}
	l = len(m.OutboundFee)
	if l > 0 {
		n += 1 + l + sovQueryOutboundFee(uint64(l))
	}
	l = len(m.FeeWithheldRune)
	if l > 0 {
		n += 1 + l + sovQueryOutboundFee(uint64(l))
	}
	l = len(m.FeeSpentRune)
	if l > 0 {
		n += 1 + l + sovQueryOutboundFee(uint64(l))
	}
	l = len(m.SurplusRune)
	if l > 0 {
		n += 1 + l + sovQueryOutboundFee(uint64(l))
	}
	l = len(m.DynamicMultiplierBasisPoints)
	if l > 0 {
		n += 1 + l + sovQueryOutboundFee(uint64(l))
	}
	return n
}

func (m *QueryOutboundFeesRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Height)
	if l > 0 {
		n += 1 + l + sovQueryOutboundFee(uint64(l))
	}
	return n
}

func (m *QueryOutboundFeesResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.OutboundFees) > 0 {
		for _, e := range m.OutboundFees {
			l = e.Size()
			n += 1 + l + sovQueryOutboundFee(uint64(l))
		}
	}
	return n
}

func sovQueryOutboundFee(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQueryOutboundFee(x uint64) (n int) {
	return sovQueryOutboundFee(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryOutboundFeeRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQueryOutboundFee
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
			return fmt.Errorf("proto: QueryOutboundFeeRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryOutboundFeeRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asset", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryOutboundFee
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
				return ErrInvalidLengthQueryOutboundFee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryOutboundFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Asset = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryOutboundFee
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
				return ErrInvalidLengthQueryOutboundFee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryOutboundFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Height = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQueryOutboundFee(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQueryOutboundFee
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
func (m *QueryOutboundFeeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQueryOutboundFee
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
			return fmt.Errorf("proto: QueryOutboundFeeResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryOutboundFeeResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asset", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryOutboundFee
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
				return ErrInvalidLengthQueryOutboundFee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryOutboundFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Asset = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OutboundFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryOutboundFee
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
				return ErrInvalidLengthQueryOutboundFee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryOutboundFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OutboundFee = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeWithheldRune", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryOutboundFee
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
				return ErrInvalidLengthQueryOutboundFee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryOutboundFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FeeWithheldRune = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeSpentRune", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryOutboundFee
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
				return ErrInvalidLengthQueryOutboundFee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryOutboundFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FeeSpentRune = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SurplusRune", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryOutboundFee
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
				return ErrInvalidLengthQueryOutboundFee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryOutboundFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SurplusRune = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DynamicMultiplierBasisPoints", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryOutboundFee
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
				return ErrInvalidLengthQueryOutboundFee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryOutboundFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DynamicMultiplierBasisPoints = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQueryOutboundFee(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQueryOutboundFee
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
func (m *QueryOutboundFeesRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQueryOutboundFee
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
			return fmt.Errorf("proto: QueryOutboundFeesRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryOutboundFeesRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryOutboundFee
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
				return ErrInvalidLengthQueryOutboundFee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryOutboundFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Height = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQueryOutboundFee(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQueryOutboundFee
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
func (m *QueryOutboundFeesResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQueryOutboundFee
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
			return fmt.Errorf("proto: QueryOutboundFeesResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryOutboundFeesResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OutboundFees", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryOutboundFee
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
				return ErrInvalidLengthQueryOutboundFee
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQueryOutboundFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OutboundFees = append(m.OutboundFees, &QueryOutboundFeeResponse{})
			if err := m.OutboundFees[len(m.OutboundFees)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQueryOutboundFee(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQueryOutboundFee
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
func skipQueryOutboundFee(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQueryOutboundFee
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
					return 0, ErrIntOverflowQueryOutboundFee
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
					return 0, ErrIntOverflowQueryOutboundFee
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
				return 0, ErrInvalidLengthQueryOutboundFee
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQueryOutboundFee
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQueryOutboundFee
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQueryOutboundFee        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQueryOutboundFee          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQueryOutboundFee = fmt.Errorf("proto: unexpected end of group")
)
