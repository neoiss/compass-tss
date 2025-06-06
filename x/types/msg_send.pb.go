// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: types/msg_send.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
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

type MsgSend struct {
	FromAddress github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,1,opt,name=from_address,json=fromAddress,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"from_address,omitempty"`
	ToAddress   github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,2,opt,name=to_address,json=toAddress,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"to_address,omitempty"`
	Amount      github_com_cosmos_cosmos_sdk_types.Coins      `protobuf:"bytes,3,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
}

func (m *MsgSend) Reset()         { *m = MsgSend{} }
func (m *MsgSend) String() string { return proto.CompactTextString(m) }
func (*MsgSend) ProtoMessage()    {}
func (*MsgSend) Descriptor() ([]byte, []int) {
	return fileDescriptor_376f985f642e301a, []int{0}
}
func (m *MsgSend) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSend) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSend.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSend) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSend.Merge(m, src)
}
func (m *MsgSend) XXX_Size() int {
	return m.Size()
}
func (m *MsgSend) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSend.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSend proto.InternalMessageInfo

func (m *MsgSend) GetFromAddress() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.FromAddress
	}
	return nil
}

func (m *MsgSend) GetToAddress() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.ToAddress
	}
	return nil
}

func (m *MsgSend) GetAmount() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Amount
	}
	return nil
}

func init() {
	proto.RegisterType((*MsgSend)(nil), "types.MsgSend")
}

func init() { proto.RegisterFile("types/msg_send.proto", fileDescriptor_376f985f642e301a) }

var fileDescriptor_376f985f642e301a = []byte{
	// 333 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x91, 0xbf, 0x4f, 0x02, 0x31,
	0x1c, 0xc5, 0xef, 0x20, 0x62, 0x3c, 0x58, 0x20, 0xc4, 0x20, 0x43, 0x21, 0x4e, 0x84, 0x84, 0x56,
	0x60, 0x33, 0x2e, 0xe0, 0xaa, 0x0b, 0x6e, 0x3a, 0x90, 0x5e, 0xaf, 0xde, 0x5d, 0xf4, 0xfa, 0x25,
	0xd7, 0x42, 0xf4, 0x5f, 0x70, 0x72, 0xf6, 0x4f, 0x70, 0xe2, 0x8f, 0x70, 0x60, 0x64, 0x74, 0x42,
	0x03, 0x03, 0xff, 0x83, 0x93, 0xe9, 0xb5, 0xfe, 0x18, 0x1d, 0x5c, 0xee, 0x5e, 0x5e, 0xfa, 0x3e,
	0xaf, 0x79, 0xf5, 0xaa, 0xea, 0x7e, 0xc2, 0x25, 0x49, 0x64, 0x38, 0x96, 0x5c, 0x04, 0x78, 0x92,
	0x82, 0x82, 0xca, 0x4e, 0xe6, 0xd6, 0x11, 0x03, 0x99, 0x80, 0x24, 0x3e, 0x95, 0x9c, 0xcc, 0xba,
	0x3e, 0x57, 0xb4, 0x4b, 0x18, 0xc4, 0xc2, 0x1c, 0xab, 0x57, 0x43, 0x08, 0x21, 0x93, 0x44, 0x2b,
	0xeb, 0x96, 0x69, 0x12, 0x0b, 0x20, 0xd9, 0xd7, 0x58, 0x87, 0x2f, 0x39, 0x6f, 0xf7, 0x5c, 0x86,
	0x17, 0x5c, 0x04, 0x95, 0xb1, 0x57, 0xba, 0x4e, 0x21, 0x19, 0xd3, 0x20, 0x48, 0xb9, 0x94, 0x35,
	0xb7, 0xe9, 0xb6, 0x4a, 0xc3, 0x93, 0x8f, 0x55, 0xa3, 0x13, 0xc6, 0x2a, 0x9a, 0xfa, 0x98, 0x41,
	0x42, 0x6c, 0xb3, 0xf9, 0x75, 0x64, 0x70, 0x43, 0xb2, 0x2b, 0xe1, 0x01, 0x63, 0x03, 0x13, 0x7c,
	0xda, 0xce, 0xdb, 0x05, 0x9f, 0xb3, 0xa8, 0xdf, 0x1b, 0x15, 0x35, 0xd1, 0xfa, 0x95, 0x2b, 0xcf,
	0x53, 0xf0, 0x8d, 0xcf, 0xfd, 0x03, 0x7e, 0x4f, 0xc1, 0x17, 0x9c, 0x79, 0x05, 0x9a, 0xc0, 0x54,
	0xa8, 0x5a, 0xbe, 0x99, 0x6f, 0x15, 0x7b, 0x07, 0xd8, 0x30, 0xb0, 0xde, 0x08, 0xdb, 0x8d, 0xf0,
	0x29, 0xc4, 0x62, 0x78, 0xb4, 0x58, 0x35, 0x9c, 0xe7, 0xb7, 0x46, 0xeb, 0x0f, 0xbd, 0x3a, 0x20,
	0x47, 0x16, 0x7d, 0xbc, 0xff, 0xb0, 0x9d, 0xb7, 0xcb, 0x2a, 0x82, 0x94, 0x45, 0x34, 0x16, 0xc4,
	0x4e, 0x37, 0x3c, 0x5b, 0xac, 0x91, 0xbb, 0x5c, 0x23, 0xf7, 0x7d, 0x8d, 0xdc, 0xc7, 0x0d, 0x72,
	0x96, 0x1b, 0xe4, 0xbc, 0x6e, 0x90, 0x73, 0xd9, 0x0b, 0x63, 0x75, 0x4b, 0x4d, 0xc7, 0x4f, 0x4e,
	0x2b, 0x01, 0x01, 0x27, 0xb3, 0x3e, 0xb9, 0xfb, 0xed, 0xeb, 0x4e, 0xbf, 0x90, 0xbd, 0x4d, 0xff,
	0x33, 0x00, 0x00, 0xff, 0xff, 0x9b, 0x71, 0xca, 0x31, 0x03, 0x02, 0x00, 0x00,
}

func (m *MsgSend) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSend) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSend) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMsgSend(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.ToAddress) > 0 {
		i -= len(m.ToAddress)
		copy(dAtA[i:], m.ToAddress)
		i = encodeVarintMsgSend(dAtA, i, uint64(len(m.ToAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.FromAddress) > 0 {
		i -= len(m.FromAddress)
		copy(dAtA[i:], m.FromAddress)
		i = encodeVarintMsgSend(dAtA, i, uint64(len(m.FromAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintMsgSend(dAtA []byte, offset int, v uint64) int {
	offset -= sovMsgSend(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgSend) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.FromAddress)
	if l > 0 {
		n += 1 + l + sovMsgSend(uint64(l))
	}
	l = len(m.ToAddress)
	if l > 0 {
		n += 1 + l + sovMsgSend(uint64(l))
	}
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovMsgSend(uint64(l))
		}
	}
	return n
}

func sovMsgSend(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMsgSend(x uint64) (n int) {
	return sovMsgSend(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgSend) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsgSend
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
			return fmt.Errorf("proto: MsgSend: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSend: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FromAddress", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgSend
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsgSend
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgSend
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FromAddress = append(m.FromAddress[:0], dAtA[iNdEx:postIndex]...)
			if m.FromAddress == nil {
				m.FromAddress = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ToAddress", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgSend
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsgSend
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgSend
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ToAddress = append(m.ToAddress[:0], dAtA[iNdEx:postIndex]...)
			if m.ToAddress == nil {
				m.ToAddress = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgSend
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
				return ErrInvalidLengthMsgSend
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMsgSend
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = append(m.Amount, types.Coin{})
			if err := m.Amount[len(m.Amount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsgSend(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsgSend
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
func skipMsgSend(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMsgSend
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
					return 0, ErrIntOverflowMsgSend
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
					return 0, ErrIntOverflowMsgSend
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
				return 0, ErrInvalidLengthMsgSend
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMsgSend
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMsgSend
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMsgSend        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMsgSend          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMsgSend = fmt.Errorf("proto: unexpected end of group")
)
