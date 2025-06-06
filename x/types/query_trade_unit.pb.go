// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: types/query_trade_unit.proto

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

type QueryTradeUnitRequest struct {
	Asset  string `protobuf:"bytes,1,opt,name=asset,proto3" json:"asset,omitempty"`
	Height string `protobuf:"bytes,2,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *QueryTradeUnitRequest) Reset()         { *m = QueryTradeUnitRequest{} }
func (m *QueryTradeUnitRequest) String() string { return proto.CompactTextString(m) }
func (*QueryTradeUnitRequest) ProtoMessage()    {}
func (*QueryTradeUnitRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5068f0a245ebc936, []int{0}
}
func (m *QueryTradeUnitRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTradeUnitRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTradeUnitRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTradeUnitRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTradeUnitRequest.Merge(m, src)
}
func (m *QueryTradeUnitRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryTradeUnitRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTradeUnitRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTradeUnitRequest proto.InternalMessageInfo

func (m *QueryTradeUnitRequest) GetAsset() string {
	if m != nil {
		return m.Asset
	}
	return ""
}

func (m *QueryTradeUnitRequest) GetHeight() string {
	if m != nil {
		return m.Height
	}
	return ""
}

type QueryTradeUnitResponse struct {
	// trade account asset with \"~\" separator
	Asset string `protobuf:"bytes,1,opt,name=asset,proto3" json:"asset"`
	// total units of trade asset
	Units string `protobuf:"bytes,2,opt,name=units,proto3" json:"units"`
	// total depth of trade asset
	Depth string `protobuf:"bytes,3,opt,name=depth,proto3" json:"depth"`
}

func (m *QueryTradeUnitResponse) Reset()         { *m = QueryTradeUnitResponse{} }
func (m *QueryTradeUnitResponse) String() string { return proto.CompactTextString(m) }
func (*QueryTradeUnitResponse) ProtoMessage()    {}
func (*QueryTradeUnitResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5068f0a245ebc936, []int{1}
}
func (m *QueryTradeUnitResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTradeUnitResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTradeUnitResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTradeUnitResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTradeUnitResponse.Merge(m, src)
}
func (m *QueryTradeUnitResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryTradeUnitResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTradeUnitResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTradeUnitResponse proto.InternalMessageInfo

func (m *QueryTradeUnitResponse) GetAsset() string {
	if m != nil {
		return m.Asset
	}
	return ""
}

func (m *QueryTradeUnitResponse) GetUnits() string {
	if m != nil {
		return m.Units
	}
	return ""
}

func (m *QueryTradeUnitResponse) GetDepth() string {
	if m != nil {
		return m.Depth
	}
	return ""
}

type QueryTradeUnitsRequest struct {
	Height string `protobuf:"bytes,1,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *QueryTradeUnitsRequest) Reset()         { *m = QueryTradeUnitsRequest{} }
func (m *QueryTradeUnitsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryTradeUnitsRequest) ProtoMessage()    {}
func (*QueryTradeUnitsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5068f0a245ebc936, []int{2}
}
func (m *QueryTradeUnitsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTradeUnitsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTradeUnitsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTradeUnitsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTradeUnitsRequest.Merge(m, src)
}
func (m *QueryTradeUnitsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryTradeUnitsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTradeUnitsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTradeUnitsRequest proto.InternalMessageInfo

func (m *QueryTradeUnitsRequest) GetHeight() string {
	if m != nil {
		return m.Height
	}
	return ""
}

type QueryTradeUnitsResponse struct {
	TradeUnits []*QueryTradeUnitResponse `protobuf:"bytes,1,rep,name=trade_units,json=tradeUnits,proto3" json:"trade_units,omitempty"`
}

func (m *QueryTradeUnitsResponse) Reset()         { *m = QueryTradeUnitsResponse{} }
func (m *QueryTradeUnitsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryTradeUnitsResponse) ProtoMessage()    {}
func (*QueryTradeUnitsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5068f0a245ebc936, []int{3}
}
func (m *QueryTradeUnitsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTradeUnitsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTradeUnitsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTradeUnitsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTradeUnitsResponse.Merge(m, src)
}
func (m *QueryTradeUnitsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryTradeUnitsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTradeUnitsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTradeUnitsResponse proto.InternalMessageInfo

func (m *QueryTradeUnitsResponse) GetTradeUnits() []*QueryTradeUnitResponse {
	if m != nil {
		return m.TradeUnits
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryTradeUnitRequest)(nil), "types.QueryTradeUnitRequest")
	proto.RegisterType((*QueryTradeUnitResponse)(nil), "types.QueryTradeUnitResponse")
	proto.RegisterType((*QueryTradeUnitsRequest)(nil), "types.QueryTradeUnitsRequest")
	proto.RegisterType((*QueryTradeUnitsResponse)(nil), "types.QueryTradeUnitsResponse")
}

func init() { proto.RegisterFile("types/query_trade_unit.proto", fileDescriptor_5068f0a245ebc936) }

var fileDescriptor_5068f0a245ebc936 = []byte{
	// 298 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x29, 0xa9, 0x2c, 0x48,
	0x2d, 0xd6, 0x2f, 0x2c, 0x4d, 0x2d, 0xaa, 0x8c, 0x2f, 0x29, 0x4a, 0x4c, 0x49, 0x8d, 0x2f, 0xcd,
	0xcb, 0x2c, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0xcb, 0x4a, 0x89, 0xa4, 0xe7,
	0xa7, 0xe7, 0x83, 0x45, 0xf4, 0x41, 0x2c, 0x88, 0xa4, 0x92, 0x2b, 0x97, 0x68, 0x20, 0x48, 0x5b,
	0x08, 0x48, 0x57, 0x68, 0x5e, 0x66, 0x49, 0x50, 0x6a, 0x61, 0x69, 0x6a, 0x71, 0x89, 0x90, 0x08,
	0x17, 0x6b, 0x62, 0x71, 0x71, 0x6a, 0x89, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x84, 0x23,
	0x24, 0xc6, 0xc5, 0x96, 0x91, 0x9a, 0x99, 0x9e, 0x51, 0x22, 0xc1, 0x04, 0x16, 0x86, 0xf2, 0x94,
	0xaa, 0xb9, 0xc4, 0xd0, 0x8d, 0x29, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x15, 0x92, 0x47, 0x31, 0xc7,
	0x89, 0xf3, 0xd5, 0x3d, 0x79, 0x88, 0x00, 0xcc, 0x48, 0x79, 0x2e, 0x56, 0x90, 0x63, 0x8b, 0x21,
	0x26, 0x42, 0x14, 0x80, 0x05, 0x82, 0x20, 0x14, 0x48, 0x41, 0x4a, 0x6a, 0x41, 0x49, 0x86, 0x04,
	0x33, 0x42, 0x01, 0x58, 0x20, 0x08, 0x42, 0x29, 0x19, 0xa0, 0x5b, 0x5e, 0x0c, 0xf3, 0x04, 0xc2,
	0xb9, 0x8c, 0x28, 0xce, 0x8d, 0xe4, 0x12, 0xc7, 0xd0, 0x01, 0x75, 0xaf, 0x1d, 0x17, 0x37, 0x22,
	0x04, 0x8b, 0x25, 0x18, 0x15, 0x98, 0x35, 0xb8, 0x8d, 0x64, 0xf5, 0xc0, 0x61, 0xa8, 0x87, 0xdd,
	0x8f, 0x41, 0x5c, 0x25, 0x70, 0x73, 0x9c, 0x7c, 0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e,
	0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58,
	0x8e, 0x21, 0xca, 0x28, 0x3d, 0xb3, 0x24, 0x27, 0x31, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0xbf, 0x24,
	0x23, 0xbf, 0x28, 0x39, 0x23, 0x31, 0x33, 0x0f, 0xcc, 0xca, 0xcb, 0x4f, 0x49, 0xd5, 0x2f, 0x33,
	0xd6, 0xaf, 0x40, 0x16, 0x07, 0x59, 0x98, 0xc4, 0x06, 0x8e, 0x25, 0x63, 0x40, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xf3, 0xa3, 0x50, 0x0d, 0xe2, 0x01, 0x00, 0x00,
}

func (m *QueryTradeUnitRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTradeUnitRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTradeUnitRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Height) > 0 {
		i -= len(m.Height)
		copy(dAtA[i:], m.Height)
		i = encodeVarintQueryTradeUnit(dAtA, i, uint64(len(m.Height)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Asset) > 0 {
		i -= len(m.Asset)
		copy(dAtA[i:], m.Asset)
		i = encodeVarintQueryTradeUnit(dAtA, i, uint64(len(m.Asset)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryTradeUnitResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTradeUnitResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTradeUnitResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Depth) > 0 {
		i -= len(m.Depth)
		copy(dAtA[i:], m.Depth)
		i = encodeVarintQueryTradeUnit(dAtA, i, uint64(len(m.Depth)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Units) > 0 {
		i -= len(m.Units)
		copy(dAtA[i:], m.Units)
		i = encodeVarintQueryTradeUnit(dAtA, i, uint64(len(m.Units)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Asset) > 0 {
		i -= len(m.Asset)
		copy(dAtA[i:], m.Asset)
		i = encodeVarintQueryTradeUnit(dAtA, i, uint64(len(m.Asset)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryTradeUnitsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTradeUnitsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTradeUnitsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Height) > 0 {
		i -= len(m.Height)
		copy(dAtA[i:], m.Height)
		i = encodeVarintQueryTradeUnit(dAtA, i, uint64(len(m.Height)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryTradeUnitsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTradeUnitsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTradeUnitsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TradeUnits) > 0 {
		for iNdEx := len(m.TradeUnits) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.TradeUnits[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQueryTradeUnit(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintQueryTradeUnit(dAtA []byte, offset int, v uint64) int {
	offset -= sovQueryTradeUnit(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryTradeUnitRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Asset)
	if l > 0 {
		n += 1 + l + sovQueryTradeUnit(uint64(l))
	}
	l = len(m.Height)
	if l > 0 {
		n += 1 + l + sovQueryTradeUnit(uint64(l))
	}
	return n
}

func (m *QueryTradeUnitResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Asset)
	if l > 0 {
		n += 1 + l + sovQueryTradeUnit(uint64(l))
	}
	l = len(m.Units)
	if l > 0 {
		n += 1 + l + sovQueryTradeUnit(uint64(l))
	}
	l = len(m.Depth)
	if l > 0 {
		n += 1 + l + sovQueryTradeUnit(uint64(l))
	}
	return n
}

func (m *QueryTradeUnitsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Height)
	if l > 0 {
		n += 1 + l + sovQueryTradeUnit(uint64(l))
	}
	return n
}

func (m *QueryTradeUnitsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.TradeUnits) > 0 {
		for _, e := range m.TradeUnits {
			l = e.Size()
			n += 1 + l + sovQueryTradeUnit(uint64(l))
		}
	}
	return n
}

func sovQueryTradeUnit(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQueryTradeUnit(x uint64) (n int) {
	return sovQueryTradeUnit(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryTradeUnitRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQueryTradeUnit
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
			return fmt.Errorf("proto: QueryTradeUnitRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTradeUnitRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asset", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryTradeUnit
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
				return ErrInvalidLengthQueryTradeUnit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryTradeUnit
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
					return ErrIntOverflowQueryTradeUnit
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
				return ErrInvalidLengthQueryTradeUnit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryTradeUnit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Height = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQueryTradeUnit(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQueryTradeUnit
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
func (m *QueryTradeUnitResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQueryTradeUnit
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
			return fmt.Errorf("proto: QueryTradeUnitResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTradeUnitResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asset", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryTradeUnit
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
				return ErrInvalidLengthQueryTradeUnit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryTradeUnit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Asset = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Units", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryTradeUnit
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
				return ErrInvalidLengthQueryTradeUnit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryTradeUnit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Units = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Depth", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryTradeUnit
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
				return ErrInvalidLengthQueryTradeUnit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryTradeUnit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Depth = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQueryTradeUnit(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQueryTradeUnit
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
func (m *QueryTradeUnitsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQueryTradeUnit
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
			return fmt.Errorf("proto: QueryTradeUnitsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTradeUnitsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryTradeUnit
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
				return ErrInvalidLengthQueryTradeUnit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQueryTradeUnit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Height = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQueryTradeUnit(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQueryTradeUnit
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
func (m *QueryTradeUnitsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQueryTradeUnit
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
			return fmt.Errorf("proto: QueryTradeUnitsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTradeUnitsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TradeUnits", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQueryTradeUnit
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
				return ErrInvalidLengthQueryTradeUnit
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQueryTradeUnit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TradeUnits = append(m.TradeUnits, &QueryTradeUnitResponse{})
			if err := m.TradeUnits[len(m.TradeUnits)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQueryTradeUnit(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQueryTradeUnit
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
func skipQueryTradeUnit(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQueryTradeUnit
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
					return 0, ErrIntOverflowQueryTradeUnit
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
					return 0, ErrIntOverflowQueryTradeUnit
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
				return 0, ErrInvalidLengthQueryTradeUnit
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQueryTradeUnit
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQueryTradeUnit
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQueryTradeUnit        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQueryTradeUnit          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQueryTradeUnit = fmt.Errorf("proto: unexpected end of group")
)
