// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: bigPackData.proto

/*
Package RPCProto is a generated protocol buffer package.

It is generated from these files:
	bigPackData.proto

It has these top-level messages:
	BigPackDataRequest
	BigPackDataResponse
*/
package RPCProto

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import bytes "bytes"

import strings "strings"
import reflect "reflect"
import sortkeys "github.com/gogo/protobuf/sortkeys"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type BigPackDataRequest struct {
	ContainerReq     map[string]string `protobuf:"bytes,1,rep,name=ContainerReq" json:"ContainerReq" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *BigPackDataRequest) Reset()                    { *m = BigPackDataRequest{} }
func (*BigPackDataRequest) ProtoMessage()               {}
func (*BigPackDataRequest) Descriptor() ([]byte, []int) { return fileDescriptorBigPackData, []int{0} }

func (m *BigPackDataRequest) GetContainerReq() map[string]string {
	if m != nil {
		return m.ContainerReq
	}
	return nil
}

type BigPackDataResponse struct {
	ContainerRes     map[string]string `protobuf:"bytes,1,rep,name=ContainerRes" json:"ContainerRes" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *BigPackDataResponse) Reset()                    { *m = BigPackDataResponse{} }
func (*BigPackDataResponse) ProtoMessage()               {}
func (*BigPackDataResponse) Descriptor() ([]byte, []int) { return fileDescriptorBigPackData, []int{1} }

func (m *BigPackDataResponse) GetContainerRes() map[string]string {
	if m != nil {
		return m.ContainerRes
	}
	return nil
}

func init() {
	proto.RegisterType((*BigPackDataRequest)(nil), "RPCProto.BigPackDataRequest")
	proto.RegisterType((*BigPackDataResponse)(nil), "RPCProto.BigPackDataResponse")
}
func (this *BigPackDataRequest) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*BigPackDataRequest)
	if !ok {
		that2, ok := that.(BigPackDataRequest)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *BigPackDataRequest")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *BigPackDataRequest but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *BigPackDataRequest but is not nil && this == nil")
	}
	if len(this.ContainerReq) != len(that1.ContainerReq) {
		return fmt.Errorf("ContainerReq this(%v) Not Equal that(%v)", len(this.ContainerReq), len(that1.ContainerReq))
	}
	for i := range this.ContainerReq {
		if this.ContainerReq[i] != that1.ContainerReq[i] {
			return fmt.Errorf("ContainerReq this[%v](%v) Not Equal that[%v](%v)", i, this.ContainerReq[i], i, that1.ContainerReq[i])
		}
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return fmt.Errorf("XXX_unrecognized this(%v) Not Equal that(%v)", this.XXX_unrecognized, that1.XXX_unrecognized)
	}
	return nil
}
func (this *BigPackDataRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*BigPackDataRequest)
	if !ok {
		that2, ok := that.(BigPackDataRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.ContainerReq) != len(that1.ContainerReq) {
		return false
	}
	for i := range this.ContainerReq {
		if this.ContainerReq[i] != that1.ContainerReq[i] {
			return false
		}
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *BigPackDataResponse) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*BigPackDataResponse)
	if !ok {
		that2, ok := that.(BigPackDataResponse)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *BigPackDataResponse")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *BigPackDataResponse but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *BigPackDataResponse but is not nil && this == nil")
	}
	if len(this.ContainerRes) != len(that1.ContainerRes) {
		return fmt.Errorf("ContainerRes this(%v) Not Equal that(%v)", len(this.ContainerRes), len(that1.ContainerRes))
	}
	for i := range this.ContainerRes {
		if this.ContainerRes[i] != that1.ContainerRes[i] {
			return fmt.Errorf("ContainerRes this[%v](%v) Not Equal that[%v](%v)", i, this.ContainerRes[i], i, that1.ContainerRes[i])
		}
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return fmt.Errorf("XXX_unrecognized this(%v) Not Equal that(%v)", this.XXX_unrecognized, that1.XXX_unrecognized)
	}
	return nil
}
func (this *BigPackDataResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*BigPackDataResponse)
	if !ok {
		that2, ok := that.(BigPackDataResponse)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.ContainerRes) != len(that1.ContainerRes) {
		return false
	}
	for i := range this.ContainerRes {
		if this.ContainerRes[i] != that1.ContainerRes[i] {
			return false
		}
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *BigPackDataRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&RPCProto.BigPackDataRequest{")
	keysForContainerReq := make([]string, 0, len(this.ContainerReq))
	for k, _ := range this.ContainerReq {
		keysForContainerReq = append(keysForContainerReq, k)
	}
	sortkeys.Strings(keysForContainerReq)
	mapStringForContainerReq := "map[string]string{"
	for _, k := range keysForContainerReq {
		mapStringForContainerReq += fmt.Sprintf("%#v: %#v,", k, this.ContainerReq[k])
	}
	mapStringForContainerReq += "}"
	if this.ContainerReq != nil {
		s = append(s, "ContainerReq: "+mapStringForContainerReq+",\n")
	}
	if this.XXX_unrecognized != nil {
		s = append(s, "XXX_unrecognized:"+fmt.Sprintf("%#v", this.XXX_unrecognized)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *BigPackDataResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&RPCProto.BigPackDataResponse{")
	keysForContainerRes := make([]string, 0, len(this.ContainerRes))
	for k, _ := range this.ContainerRes {
		keysForContainerRes = append(keysForContainerRes, k)
	}
	sortkeys.Strings(keysForContainerRes)
	mapStringForContainerRes := "map[string]string{"
	for _, k := range keysForContainerRes {
		mapStringForContainerRes += fmt.Sprintf("%#v: %#v,", k, this.ContainerRes[k])
	}
	mapStringForContainerRes += "}"
	if this.ContainerRes != nil {
		s = append(s, "ContainerRes: "+mapStringForContainerRes+",\n")
	}
	if this.XXX_unrecognized != nil {
		s = append(s, "XXX_unrecognized:"+fmt.Sprintf("%#v", this.XXX_unrecognized)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringBigPackData(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *BigPackDataRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BigPackDataRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ContainerReq) > 0 {
		for k, _ := range m.ContainerReq {
			dAtA[i] = 0xa
			i++
			v := m.ContainerReq[k]
			mapSize := 1 + len(k) + sovBigPackData(uint64(len(k))) + 1 + len(v) + sovBigPackData(uint64(len(v)))
			i = encodeVarintBigPackData(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintBigPackData(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintBigPackData(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *BigPackDataResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BigPackDataResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ContainerRes) > 0 {
		for k, _ := range m.ContainerRes {
			dAtA[i] = 0xa
			i++
			v := m.ContainerRes[k]
			mapSize := 1 + len(k) + sovBigPackData(uint64(len(k))) + 1 + len(v) + sovBigPackData(uint64(len(v)))
			i = encodeVarintBigPackData(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintBigPackData(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintBigPackData(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintBigPackData(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedBigPackDataRequest(r randyBigPackData, easy bool) *BigPackDataRequest {
	this := &BigPackDataRequest{}
	if r.Intn(10) != 0 {
		v1 := r.Intn(10)
		this.ContainerReq = make(map[string]string)
		for i := 0; i < v1; i++ {
			this.ContainerReq[randStringBigPackData(r)] = randStringBigPackData(r)
		}
	}
	if !easy && r.Intn(10) != 0 {
		this.XXX_unrecognized = randUnrecognizedBigPackData(r, 2)
	}
	return this
}

func NewPopulatedBigPackDataResponse(r randyBigPackData, easy bool) *BigPackDataResponse {
	this := &BigPackDataResponse{}
	if r.Intn(10) != 0 {
		v2 := r.Intn(10)
		this.ContainerRes = make(map[string]string)
		for i := 0; i < v2; i++ {
			this.ContainerRes[randStringBigPackData(r)] = randStringBigPackData(r)
		}
	}
	if !easy && r.Intn(10) != 0 {
		this.XXX_unrecognized = randUnrecognizedBigPackData(r, 2)
	}
	return this
}

type randyBigPackData interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneBigPackData(r randyBigPackData) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringBigPackData(r randyBigPackData) string {
	v3 := r.Intn(100)
	tmps := make([]rune, v3)
	for i := 0; i < v3; i++ {
		tmps[i] = randUTF8RuneBigPackData(r)
	}
	return string(tmps)
}
func randUnrecognizedBigPackData(r randyBigPackData, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldBigPackData(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldBigPackData(dAtA []byte, r randyBigPackData, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateBigPackData(dAtA, uint64(key))
		v4 := r.Int63()
		if r.Intn(2) == 0 {
			v4 *= -1
		}
		dAtA = encodeVarintPopulateBigPackData(dAtA, uint64(v4))
	case 1:
		dAtA = encodeVarintPopulateBigPackData(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateBigPackData(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateBigPackData(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateBigPackData(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateBigPackData(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(uint64(v)&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *BigPackDataRequest) Size() (n int) {
	var l int
	_ = l
	if len(m.ContainerReq) > 0 {
		for k, v := range m.ContainerReq {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovBigPackData(uint64(len(k))) + 1 + len(v) + sovBigPackData(uint64(len(v)))
			n += mapEntrySize + 1 + sovBigPackData(uint64(mapEntrySize))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *BigPackDataResponse) Size() (n int) {
	var l int
	_ = l
	if len(m.ContainerRes) > 0 {
		for k, v := range m.ContainerRes {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovBigPackData(uint64(len(k))) + 1 + len(v) + sovBigPackData(uint64(len(v)))
			n += mapEntrySize + 1 + sovBigPackData(uint64(mapEntrySize))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovBigPackData(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozBigPackData(x uint64) (n int) {
	return sovBigPackData(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *BigPackDataRequest) String() string {
	if this == nil {
		return "nil"
	}
	keysForContainerReq := make([]string, 0, len(this.ContainerReq))
	for k, _ := range this.ContainerReq {
		keysForContainerReq = append(keysForContainerReq, k)
	}
	sortkeys.Strings(keysForContainerReq)
	mapStringForContainerReq := "map[string]string{"
	for _, k := range keysForContainerReq {
		mapStringForContainerReq += fmt.Sprintf("%v: %v,", k, this.ContainerReq[k])
	}
	mapStringForContainerReq += "}"
	s := strings.Join([]string{`&BigPackDataRequest{`,
		`ContainerReq:` + mapStringForContainerReq + `,`,
		`XXX_unrecognized:` + fmt.Sprintf("%v", this.XXX_unrecognized) + `,`,
		`}`,
	}, "")
	return s
}
func (this *BigPackDataResponse) String() string {
	if this == nil {
		return "nil"
	}
	keysForContainerRes := make([]string, 0, len(this.ContainerRes))
	for k, _ := range this.ContainerRes {
		keysForContainerRes = append(keysForContainerRes, k)
	}
	sortkeys.Strings(keysForContainerRes)
	mapStringForContainerRes := "map[string]string{"
	for _, k := range keysForContainerRes {
		mapStringForContainerRes += fmt.Sprintf("%v: %v,", k, this.ContainerRes[k])
	}
	mapStringForContainerRes += "}"
	s := strings.Join([]string{`&BigPackDataResponse{`,
		`ContainerRes:` + mapStringForContainerRes + `,`,
		`XXX_unrecognized:` + fmt.Sprintf("%v", this.XXX_unrecognized) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringBigPackData(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *BigPackDataRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBigPackData
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BigPackDataRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BigPackDataRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContainerReq", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBigPackData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthBigPackData
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ContainerReq == nil {
				m.ContainerReq = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowBigPackData
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowBigPackData
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthBigPackData
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowBigPackData
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthBigPackData
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipBigPackData(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthBigPackData
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.ContainerReq[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBigPackData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBigPackData
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BigPackDataResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBigPackData
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BigPackDataResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BigPackDataResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContainerRes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBigPackData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthBigPackData
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ContainerRes == nil {
				m.ContainerRes = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowBigPackData
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowBigPackData
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthBigPackData
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowBigPackData
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthBigPackData
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipBigPackData(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthBigPackData
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.ContainerRes[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBigPackData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBigPackData
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipBigPackData(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBigPackData
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
					return 0, ErrIntOverflowBigPackData
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBigPackData
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthBigPackData
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowBigPackData
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipBigPackData(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthBigPackData = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBigPackData   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("bigPackData.proto", fileDescriptorBigPackData) }

var fileDescriptorBigPackData = []byte{
	// 273 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0xca, 0x4c, 0x0f,
	0x48, 0x4c, 0xce, 0x76, 0x49, 0x2c, 0x49, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x08,
	0x0a, 0x70, 0x0e, 0x00, 0xb1, 0xa4, 0x74, 0xd3, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3,
	0x73, 0xf5, 0xd3, 0xf3, 0xd3, 0xf3, 0xf5, 0xc1, 0x0a, 0x92, 0x4a, 0xd3, 0xc0, 0x3c, 0x30, 0x07,
	0xcc, 0x82, 0x68, 0x54, 0x5a, 0xcf, 0xc8, 0x25, 0xe4, 0x84, 0x30, 0x2e, 0x28, 0xb5, 0xb0, 0x34,
	0xb5, 0xb8, 0x44, 0x28, 0x82, 0x8b, 0xc7, 0x39, 0x3f, 0xaf, 0x24, 0x31, 0x33, 0x2f, 0xb5, 0x28,
	0x28, 0xb5, 0x50, 0x82, 0x51, 0x81, 0x59, 0x83, 0xdb, 0x48, 0x4f, 0x0f, 0x66, 0x8d, 0x1e, 0xa6,
	0x1e, 0x3d, 0x64, 0x0d, 0xae, 0x79, 0x25, 0x45, 0x95, 0x4e, 0x2c, 0x27, 0xee, 0xc9, 0x33, 0x04,
	0xa1, 0x98, 0x24, 0x65, 0xcf, 0x25, 0x88, 0xa1, 0x50, 0x48, 0x80, 0x8b, 0x39, 0x3b, 0xb5, 0x52,
	0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0xc4, 0x14, 0x12, 0xe1, 0x62, 0x2d, 0x4b, 0xcc, 0x29,
	0x4d, 0x95, 0x60, 0x02, 0x8b, 0x41, 0x38, 0x56, 0x4c, 0x16, 0x8c, 0x4a, 0x1b, 0x19, 0xb9, 0x84,
	0x51, 0x6c, 0x2f, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x15, 0x8a, 0x44, 0x71, 0x72, 0x31, 0xd4, 0xc9,
	0xfa, 0x38, 0x9c, 0x0c, 0xd1, 0x84, 0xec, 0xe6, 0x62, 0x5c, 0x6e, 0x2e, 0x46, 0x73, 0x73, 0x31,
	0xc9, 0x6e, 0x76, 0xd2, 0xb9, 0xf1, 0x50, 0x8e, 0xe1, 0xc1, 0x43, 0x39, 0xc6, 0x0f, 0x0f, 0xe5,
	0x18, 0x7f, 0x3c, 0x94, 0x63, 0x6c, 0x78, 0x24, 0xc7, 0xb8, 0xe2, 0x91, 0x1c, 0xe3, 0x8e, 0x47,
	0x72, 0x8c, 0x07, 0x1e, 0xc9, 0x31, 0x9e, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83,
	0x47, 0x72, 0x8c, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xdc, 0x02, 0x54, 0xa6, 0xe0, 0x01, 0x00,
	0x00,
}