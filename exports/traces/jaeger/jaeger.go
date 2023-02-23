package jaeger

import (
	"bytes"
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/fuyibing/log/exports/traces/jaeger/thrift"
	"time"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = context.Background
var _ = time.Now
var _ = bytes.Equal

type TagType int64

const (
	TagType_STRING TagType = 0
	TagType_DOUBLE TagType = 1
	TagType_BOOL   TagType = 2
	TagType_LONG   TagType = 3
	TagType_BINARY TagType = 4
)

func (p TagType) String() string {
	switch p {
	case TagType_STRING:
		return "STRING"
	case TagType_DOUBLE:
		return "DOUBLE"
	case TagType_BOOL:
		return "BOOL"
	case TagType_LONG:
		return "LONG"
	case TagType_BINARY:
		return "BINARY"
	}
	return "<UNSET>"
}

type SpanRefType int64

const (
	SpanRefType_CHILD_OF     SpanRefType = 0
	SpanRefType_FOLLOWS_FROM SpanRefType = 1
)

func (p SpanRefType) String() string {
	switch p {
	case SpanRefType_CHILD_OF:
		return "CHILD_OF"
	case SpanRefType_FOLLOWS_FROM:
		return "FOLLOWS_FROM"
	}
	return "<UNSET>"
}

func SpanRefTypeFromString(s string) (SpanRefType, error) {
	switch s {
	case "CHILD_OF":
		return SpanRefType_CHILD_OF, nil
	case "FOLLOWS_FROM":
		return SpanRefType_FOLLOWS_FROM, nil
	}
	return SpanRefType(0), fmt.Errorf("not a valid SpanRefType string")
}

func (p SpanRefType) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *SpanRefType) UnmarshalText(text []byte) error {
	q, err := SpanRefTypeFromString(string(text))
	if err != nil {
		return err
	}
	*p = q
	return nil
}

func (p *SpanRefType) Scan(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return errors.New("Scan value is not int64")
	}
	*p = SpanRefType(v)
	return nil
}

func (p *SpanRefType) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}

// Attributes:
//  - Key
//  - VType
//  - VStr
//  - VDouble
//  - VBool
//  - VLong
//  - VBinary
type Tag struct {
	Key     string   `thrift:"key,1,required" db:"key" json:"key"`
	VType   TagType  `thrift:"vType,2,required" db:"vType" json:"vType"`
	VStr    *string  `thrift:"vStr,3" db:"vStr" json:"vStr,omitempty"`
	VDouble *float64 `thrift:"vDouble,4" db:"vDouble" json:"vDouble,omitempty"`
	VBool   *bool    `thrift:"vBool,5" db:"vBool" json:"vBool,omitempty"`
	VLong   *int64   `thrift:"vLong,6" db:"vLong" json:"vLong,omitempty"`
	VBinary []byte   `thrift:"vBinary,7" db:"vBinary" json:"vBinary,omitempty"`
}

func (p *Tag) GetKey() string {
	return p.Key
}

func (p *Tag) GetVType() TagType {
	return p.VType
}

var Tag_VStr_DEFAULT string

func (p *Tag) GetVStr() string {
	if !p.IsSetVStr() {
		return Tag_VStr_DEFAULT
	}
	return *p.VStr
}

var Tag_VDouble_DEFAULT float64

func (p *Tag) GetVDouble() float64 {
	if !p.IsSetVDouble() {
		return Tag_VDouble_DEFAULT
	}
	return *p.VDouble
}

var Tag_VBool_DEFAULT bool

func (p *Tag) GetVBool() bool {
	if !p.IsSetVBool() {
		return Tag_VBool_DEFAULT
	}
	return *p.VBool
}

var Tag_VLong_DEFAULT int64

func (p *Tag) GetVLong() int64 {
	if !p.IsSetVLong() {
		return Tag_VLong_DEFAULT
	}
	return *p.VLong
}

var Tag_VBinary_DEFAULT []byte

func (p *Tag) GetVBinary() []byte {
	return p.VBinary
}
func (p *Tag) IsSetVStr() bool {
	return p.VStr != nil
}

func (p *Tag) IsSetVDouble() bool {
	return p.VDouble != nil
}

func (p *Tag) IsSetVBool() bool {
	return p.VBool != nil
}

func (p *Tag) IsSetVLong() bool {
	return p.VLong != nil
}

func (p *Tag) IsSetVBinary() bool {
	return p.VBinary != nil
}

func (p *Tag) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetKey bool = false
	var issetVType bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
				issetKey = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 2:
			if fieldTypeId == thrift.I32 {
				if err := p.ReadField2(ctx, iprot); err != nil {
					return err
				}
				issetVType = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 3:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField3(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 4:
			if fieldTypeId == thrift.DOUBLE {
				if err := p.ReadField4(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 5:
			if fieldTypeId == thrift.BOOL {
				if err := p.ReadField5(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 6:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField6(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 7:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField7(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetKey {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Key is not set"))
	}
	if !issetVType {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field VType is not set"))
	}
	return nil
}

func (p *Tag) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(ctx); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Key = v
	}
	return nil
}

func (p *Tag) ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(ctx); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		temp := TagType(v)
		p.VType = temp
	}
	return nil
}

func (p *Tag) ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(ctx); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.VStr = &v
	}
	return nil
}

func (p *Tag) ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadDouble(ctx); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		p.VDouble = &v
	}
	return nil
}

func (p *Tag) ReadField5(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(ctx); err != nil {
		return thrift.PrependError("error reading field 5: ", err)
	} else {
		p.VBool = &v
	}
	return nil
}

func (p *Tag) ReadField6(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 6: ", err)
	} else {
		p.VLong = &v
	}
	return nil
}

func (p *Tag) ReadField7(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(ctx); err != nil {
		return thrift.PrependError("error reading field 7: ", err)
	} else {
		p.VBinary = v
	}
	return nil
}

func (p *Tag) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "Tag"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField2(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField3(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField4(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField5(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField6(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField7(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Tag) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "key", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:key: ", p), err)
	}
	if err := oprot.WriteString(ctx, string(p.Key)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.key (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:key: ", p), err)
	}
	return err
}

func (p *Tag) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "vType", thrift.I32, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:vType: ", p), err)
	}
	if err := oprot.WriteI32(ctx, int32(p.VType)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.vType (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:vType: ", p), err)
	}
	return err
}

func (p *Tag) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if p.IsSetVStr() {
		if err := oprot.WriteFieldBegin(ctx, "vStr", thrift.STRING, 3); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:vStr: ", p), err)
		}
		if err := oprot.WriteString(ctx, string(*p.VStr)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.vStr (3) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(ctx); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 3:vStr: ", p), err)
		}
	}
	return err
}

func (p *Tag) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if p.IsSetVDouble() {
		if err := oprot.WriteFieldBegin(ctx, "vDouble", thrift.DOUBLE, 4); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:vDouble: ", p), err)
		}
		if err := oprot.WriteDouble(ctx, float64(*p.VDouble)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.vDouble (4) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(ctx); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 4:vDouble: ", p), err)
		}
	}
	return err
}

func (p *Tag) writeField5(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if p.IsSetVBool() {
		if err := oprot.WriteFieldBegin(ctx, "vBool", thrift.BOOL, 5); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:vBool: ", p), err)
		}
		if err := oprot.WriteBool(ctx, bool(*p.VBool)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.vBool (5) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(ctx); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 5:vBool: ", p), err)
		}
	}
	return err
}

func (p *Tag) writeField6(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if p.IsSetVLong() {
		if err := oprot.WriteFieldBegin(ctx, "vLong", thrift.I64, 6); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:vLong: ", p), err)
		}
		if err := oprot.WriteI64(ctx, int64(*p.VLong)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.vLong (6) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(ctx); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 6:vLong: ", p), err)
		}
	}
	return err
}

func (p *Tag) writeField7(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if p.IsSetVBinary() {
		if err := oprot.WriteFieldBegin(ctx, "vBinary", thrift.STRING, 7); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 7:vBinary: ", p), err)
		}
		if err := oprot.WriteBinary(ctx, p.VBinary); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.vBinary (7) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(ctx); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 7:vBinary: ", p), err)
		}
	}
	return err
}

func (p *Tag) Equals(other *Tag) bool {
	if p == other {
		return true
	} else if p == nil || other == nil {
		return false
	}
	if p.Key != other.Key {
		return false
	}
	if p.VType != other.VType {
		return false
	}
	if p.VStr != other.VStr {
		if p.VStr == nil || other.VStr == nil {
			return false
		}
		if (*p.VStr) != (*other.VStr) {
			return false
		}
	}
	if p.VDouble != other.VDouble {
		if p.VDouble == nil || other.VDouble == nil {
			return false
		}
		if (*p.VDouble) != (*other.VDouble) {
			return false
		}
	}
	if p.VBool != other.VBool {
		if p.VBool == nil || other.VBool == nil {
			return false
		}
		if (*p.VBool) != (*other.VBool) {
			return false
		}
	}
	if p.VLong != other.VLong {
		if p.VLong == nil || other.VLong == nil {
			return false
		}
		if (*p.VLong) != (*other.VLong) {
			return false
		}
	}
	if bytes.Compare(p.VBinary, other.VBinary) != 0 {
		return false
	}
	return true
}

func (p *Tag) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Tag(%+v)", *p)
}

// Attributes:
//  - Timestamp
//  - Fields
type Log struct {
	Timestamp int64  `thrift:"timestamp,1,required" db:"timestamp" json:"timestamp"`
	Fields    []*Tag `thrift:"fields,2,required" db:"fields" json:"fields"`
}

func (p *Log) GetTimestamp() int64 {
	return p.Timestamp
}

func (p *Log) GetFields() []*Tag {
	return p.Fields
}
func (p *Log) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetTimestamp bool = false
	var issetFields bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
				issetTimestamp = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 2:
			if fieldTypeId == thrift.LIST {
				if err := p.ReadField2(ctx, iprot); err != nil {
					return err
				}
				issetFields = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetTimestamp {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Timestamp is not set"))
	}
	if !issetFields {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Fields is not set"))
	}
	return nil
}

func (p *Log) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Timestamp = v
	}
	return nil
}

func (p *Log) ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin(ctx)
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*Tag, 0, size)
	p.Fields = tSlice
	for i := 0; i < size; i++ {
		_elem0 := &Tag{}
		if err := _elem0.Read(ctx, iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem0), err)
		}
		p.Fields = append(p.Fields, _elem0)
	}
	if err := iprot.ReadListEnd(ctx); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *Log) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "Log"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField2(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Log) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "timestamp", thrift.I64, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:timestamp: ", p), err)
	}
	if err := oprot.WriteI64(ctx, int64(p.Timestamp)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.timestamp (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:timestamp: ", p), err)
	}
	return err
}

func (p *Log) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "fields", thrift.LIST, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:fields: ", p), err)
	}
	if err := oprot.WriteListBegin(ctx, thrift.STRUCT, len(p.Fields)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.Fields {
		if err := v.Write(ctx, oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
		}
	}
	if err := oprot.WriteListEnd(ctx); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:fields: ", p), err)
	}
	return err
}

func (p *Log) Equals(other *Log) bool {
	if p == other {
		return true
	} else if p == nil || other == nil {
		return false
	}
	if p.Timestamp != other.Timestamp {
		return false
	}
	if len(p.Fields) != len(other.Fields) {
		return false
	}
	for i, _tgt := range p.Fields {
		_src1 := other.Fields[i]
		if !_tgt.Equals(_src1) {
			return false
		}
	}
	return true
}

func (p *Log) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Log(%+v)", *p)
}

// Attributes:
//  - RefType
//  - TraceIdLow
//  - TraceIdHigh
//  - SpanId
type SpanRef struct {
	RefType     SpanRefType `thrift:"refType,1,required" db:"refType" json:"refType"`
	TraceIdLow  int64       `thrift:"traceIdLow,2,required" db:"traceIdLow" json:"traceIdLow"`
	TraceIdHigh int64       `thrift:"traceIdHigh,3,required" db:"traceIdHigh" json:"traceIdHigh"`
	SpanId      int64       `thrift:"spanId,4,required" db:"spanId" json:"spanId"`
}

func (p *SpanRef) GetRefType() SpanRefType {
	return p.RefType
}

func (p *SpanRef) GetTraceIdLow() int64 {
	return p.TraceIdLow
}

func (p *SpanRef) GetTraceIdHigh() int64 {
	return p.TraceIdHigh
}

func (p *SpanRef) GetSpanId() int64 {
	return p.SpanId
}
func (p *SpanRef) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetRefType bool = false
	var issetTraceIdLow bool = false
	var issetTraceIdHigh bool = false
	var issetSpanId bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.I32 {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
				issetRefType = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 2:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField2(ctx, iprot); err != nil {
					return err
				}
				issetTraceIdLow = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 3:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField3(ctx, iprot); err != nil {
					return err
				}
				issetTraceIdHigh = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 4:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField4(ctx, iprot); err != nil {
					return err
				}
				issetSpanId = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetRefType {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field RefType is not set"))
	}
	if !issetTraceIdLow {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field TraceIdLow is not set"))
	}
	if !issetTraceIdHigh {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field TraceIdHigh is not set"))
	}
	if !issetSpanId {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field SpanId is not set"))
	}
	return nil
}

func (p *SpanRef) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(ctx); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		temp := SpanRefType(v)
		p.RefType = temp
	}
	return nil
}

func (p *SpanRef) ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.TraceIdLow = v
	}
	return nil
}

func (p *SpanRef) ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.TraceIdHigh = v
	}
	return nil
}

func (p *SpanRef) ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		p.SpanId = v
	}
	return nil
}

func (p *SpanRef) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "SpanRef"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField2(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField3(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField4(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *SpanRef) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "refType", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:refType: ", p), err)
	}
	if err := oprot.WriteI32(ctx, int32(p.RefType)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.refType (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:refType: ", p), err)
	}
	return err
}

func (p *SpanRef) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "traceIdLow", thrift.I64, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:traceIdLow: ", p), err)
	}
	if err := oprot.WriteI64(ctx, int64(p.TraceIdLow)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.traceIdLow (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:traceIdLow: ", p), err)
	}
	return err
}

func (p *SpanRef) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "traceIdHigh", thrift.I64, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:traceIdHigh: ", p), err)
	}
	if err := oprot.WriteI64(ctx, int64(p.TraceIdHigh)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.traceIdHigh (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:traceIdHigh: ", p), err)
	}
	return err
}

func (p *SpanRef) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "spanId", thrift.I64, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:spanId: ", p), err)
	}
	if err := oprot.WriteI64(ctx, int64(p.SpanId)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.spanId (4) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:spanId: ", p), err)
	}
	return err
}

func (p *SpanRef) Equals(other *SpanRef) bool {
	if p == other {
		return true
	} else if p == nil || other == nil {
		return false
	}
	if p.RefType != other.RefType {
		return false
	}
	if p.TraceIdLow != other.TraceIdLow {
		return false
	}
	if p.TraceIdHigh != other.TraceIdHigh {
		return false
	}
	if p.SpanId != other.SpanId {
		return false
	}
	return true
}

func (p *SpanRef) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SpanRef(%+v)", *p)
}

// Attributes:
//  - TraceIdLow
//  - TraceIdHigh
//  - SpanId
//  - ParentSpanId
//  - OperationName
//  - References
//  - Flags
//  - StartTime
//  - Duration
//  - Tags
//  - Logs
type Span struct {
	TraceIdLow    int64      `thrift:"traceIdLow,1,required" db:"traceIdLow" json:"traceIdLow"`
	TraceIdHigh   int64      `thrift:"traceIdHigh,2,required" db:"traceIdHigh" json:"traceIdHigh"`
	SpanId        int64      `thrift:"spanId,3,required" db:"spanId" json:"spanId"`
	ParentSpanId  int64      `thrift:"parentSpanId,4,required" db:"parentSpanId" json:"parentSpanId"`
	OperationName string     `thrift:"operationName,5,required" db:"operationName" json:"operationName"`
	References    []*SpanRef `thrift:"references,6" db:"references" json:"references,omitempty"`
	Flags         int32      `thrift:"flags,7,required" db:"flags" json:"flags"`
	StartTime     int64      `thrift:"startTime,8,required" db:"startTime" json:"startTime"`
	Duration      int64      `thrift:"duration,9,required" db:"duration" json:"duration"`
	Tags          []*Tag     `thrift:"tags,10" db:"tags" json:"tags,omitempty"`
	Logs          []*Log     `thrift:"logs,11" db:"logs" json:"logs,omitempty"`
}

func (p *Span) GetTraceIdLow() int64 {
	return p.TraceIdLow
}

func (p *Span) GetTraceIdHigh() int64 {
	return p.TraceIdHigh
}

func (p *Span) GetSpanId() int64 {
	return p.SpanId
}

func (p *Span) GetParentSpanId() int64 {
	return p.ParentSpanId
}

func (p *Span) GetOperationName() string {
	return p.OperationName
}

var Span_References_DEFAULT []*SpanRef

func (p *Span) GetReferences() []*SpanRef {
	return p.References
}

func (p *Span) GetFlags() int32 {
	return p.Flags
}

func (p *Span) GetStartTime() int64 {
	return p.StartTime
}

func (p *Span) GetDuration() int64 {
	return p.Duration
}

var Span_Tags_DEFAULT []*Tag

func (p *Span) GetTags() []*Tag {
	return p.Tags
}

var Span_Logs_DEFAULT []*Log

func (p *Span) GetLogs() []*Log {
	return p.Logs
}
func (p *Span) IsSetReferences() bool {
	return p.References != nil
}

func (p *Span) IsSetTags() bool {
	return p.Tags != nil
}

func (p *Span) IsSetLogs() bool {
	return p.Logs != nil
}

func (p *Span) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetTraceIdLow bool = false
	var issetTraceIdHigh bool = false
	var issetSpanId bool = false
	var issetParentSpanId bool = false
	var issetOperationName bool = false
	var issetFlags bool = false
	var issetStartTime bool = false
	var issetDuration bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
				issetTraceIdLow = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 2:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField2(ctx, iprot); err != nil {
					return err
				}
				issetTraceIdHigh = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 3:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField3(ctx, iprot); err != nil {
					return err
				}
				issetSpanId = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 4:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField4(ctx, iprot); err != nil {
					return err
				}
				issetParentSpanId = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 5:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField5(ctx, iprot); err != nil {
					return err
				}
				issetOperationName = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 6:
			if fieldTypeId == thrift.LIST {
				if err := p.ReadField6(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 7:
			if fieldTypeId == thrift.I32 {
				if err := p.ReadField7(ctx, iprot); err != nil {
					return err
				}
				issetFlags = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 8:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField8(ctx, iprot); err != nil {
					return err
				}
				issetStartTime = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 9:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField9(ctx, iprot); err != nil {
					return err
				}
				issetDuration = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 10:
			if fieldTypeId == thrift.LIST {
				if err := p.ReadField10(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 11:
			if fieldTypeId == thrift.LIST {
				if err := p.ReadField11(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetTraceIdLow {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field TraceIdLow is not set"))
	}
	if !issetTraceIdHigh {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field TraceIdHigh is not set"))
	}
	if !issetSpanId {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field SpanId is not set"))
	}
	if !issetParentSpanId {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field ParentSpanId is not set"))
	}
	if !issetOperationName {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field OperationName is not set"))
	}
	if !issetFlags {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Flags is not set"))
	}
	if !issetStartTime {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field StartTime is not set"))
	}
	if !issetDuration {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Duration is not set"))
	}
	return nil
}

func (p *Span) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.TraceIdLow = v
	}
	return nil
}

func (p *Span) ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.TraceIdHigh = v
	}
	return nil
}

func (p *Span) ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.SpanId = v
	}
	return nil
}

func (p *Span) ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		p.ParentSpanId = v
	}
	return nil
}

func (p *Span) ReadField5(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(ctx); err != nil {
		return thrift.PrependError("error reading field 5: ", err)
	} else {
		p.OperationName = v
	}
	return nil
}

func (p *Span) ReadField6(ctx context.Context, iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin(ctx)
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*SpanRef, 0, size)
	p.References = tSlice
	for i := 0; i < size; i++ {
		_elem2 := &SpanRef{}
		if err := _elem2.Read(ctx, iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem2), err)
		}
		p.References = append(p.References, _elem2)
	}
	if err := iprot.ReadListEnd(ctx); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *Span) ReadField7(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(ctx); err != nil {
		return thrift.PrependError("error reading field 7: ", err)
	} else {
		p.Flags = v
	}
	return nil
}

func (p *Span) ReadField8(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 8: ", err)
	} else {
		p.StartTime = v
	}
	return nil
}

func (p *Span) ReadField9(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 9: ", err)
	} else {
		p.Duration = v
	}
	return nil
}

func (p *Span) ReadField10(ctx context.Context, iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin(ctx)
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*Tag, 0, size)
	p.Tags = tSlice
	for i := 0; i < size; i++ {
		_elem3 := &Tag{}
		if err := _elem3.Read(ctx, iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem3), err)
		}
		p.Tags = append(p.Tags, _elem3)
	}
	if err := iprot.ReadListEnd(ctx); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *Span) ReadField11(ctx context.Context, iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin(ctx)
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*Log, 0, size)
	p.Logs = tSlice
	for i := 0; i < size; i++ {
		_elem4 := &Log{}
		if err := _elem4.Read(ctx, iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem4), err)
		}
		p.Logs = append(p.Logs, _elem4)
	}
	if err := iprot.ReadListEnd(ctx); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *Span) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "Span"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField2(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField3(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField4(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField5(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField6(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField7(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField8(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField9(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField10(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField11(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Span) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "traceIdLow", thrift.I64, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:traceIdLow: ", p), err)
	}
	if err := oprot.WriteI64(ctx, int64(p.TraceIdLow)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.traceIdLow (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:traceIdLow: ", p), err)
	}
	return err
}

func (p *Span) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "traceIdHigh", thrift.I64, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:traceIdHigh: ", p), err)
	}
	if err := oprot.WriteI64(ctx, int64(p.TraceIdHigh)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.traceIdHigh (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:traceIdHigh: ", p), err)
	}
	return err
}

func (p *Span) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "spanId", thrift.I64, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:spanId: ", p), err)
	}
	if err := oprot.WriteI64(ctx, int64(p.SpanId)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.spanId (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:spanId: ", p), err)
	}
	return err
}

func (p *Span) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "parentSpanId", thrift.I64, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:parentSpanId: ", p), err)
	}
	if err := oprot.WriteI64(ctx, int64(p.ParentSpanId)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.parentSpanId (4) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:parentSpanId: ", p), err)
	}
	return err
}

func (p *Span) writeField5(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "operationName", thrift.STRING, 5); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:operationName: ", p), err)
	}
	if err := oprot.WriteString(ctx, string(p.OperationName)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.operationName (5) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 5:operationName: ", p), err)
	}
	return err
}

func (p *Span) writeField6(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if p.IsSetReferences() {
		if err := oprot.WriteFieldBegin(ctx, "references", thrift.LIST, 6); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:references: ", p), err)
		}
		if err := oprot.WriteListBegin(ctx, thrift.STRUCT, len(p.References)); err != nil {
			return thrift.PrependError("error writing list begin: ", err)
		}
		for _, v := range p.References {
			if err := v.Write(ctx, oprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
			}
		}
		if err := oprot.WriteListEnd(ctx); err != nil {
			return thrift.PrependError("error writing list end: ", err)
		}
		if err := oprot.WriteFieldEnd(ctx); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 6:references: ", p), err)
		}
	}
	return err
}

func (p *Span) writeField7(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "flags", thrift.I32, 7); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 7:flags: ", p), err)
	}
	if err := oprot.WriteI32(ctx, int32(p.Flags)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.flags (7) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 7:flags: ", p), err)
	}
	return err
}

func (p *Span) writeField8(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "startTime", thrift.I64, 8); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 8:startTime: ", p), err)
	}
	if err := oprot.WriteI64(ctx, int64(p.StartTime)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.startTime (8) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 8:startTime: ", p), err)
	}
	return err
}

func (p *Span) writeField9(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "duration", thrift.I64, 9); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 9:duration: ", p), err)
	}
	if err := oprot.WriteI64(ctx, int64(p.Duration)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.duration (9) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 9:duration: ", p), err)
	}
	return err
}

func (p *Span) writeField10(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if p.IsSetTags() {
		if err := oprot.WriteFieldBegin(ctx, "tags", thrift.LIST, 10); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 10:tags: ", p), err)
		}
		if err := oprot.WriteListBegin(ctx, thrift.STRUCT, len(p.Tags)); err != nil {
			return thrift.PrependError("error writing list begin: ", err)
		}
		for _, v := range p.Tags {
			if err := v.Write(ctx, oprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
			}
		}
		if err := oprot.WriteListEnd(ctx); err != nil {
			return thrift.PrependError("error writing list end: ", err)
		}
		if err := oprot.WriteFieldEnd(ctx); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 10:tags: ", p), err)
		}
	}
	return err
}

func (p *Span) writeField11(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if p.IsSetLogs() {
		if err := oprot.WriteFieldBegin(ctx, "logs", thrift.LIST, 11); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 11:logs: ", p), err)
		}
		if err := oprot.WriteListBegin(ctx, thrift.STRUCT, len(p.Logs)); err != nil {
			return thrift.PrependError("error writing list begin: ", err)
		}
		for _, v := range p.Logs {
			if err := v.Write(ctx, oprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
			}
		}
		if err := oprot.WriteListEnd(ctx); err != nil {
			return thrift.PrependError("error writing list end: ", err)
		}
		if err := oprot.WriteFieldEnd(ctx); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 11:logs: ", p), err)
		}
	}
	return err
}

func (p *Span) Equals(other *Span) bool {
	if p == other {
		return true
	} else if p == nil || other == nil {
		return false
	}
	if p.TraceIdLow != other.TraceIdLow {
		return false
	}
	if p.TraceIdHigh != other.TraceIdHigh {
		return false
	}
	if p.SpanId != other.SpanId {
		return false
	}
	if p.ParentSpanId != other.ParentSpanId {
		return false
	}
	if p.OperationName != other.OperationName {
		return false
	}
	if len(p.References) != len(other.References) {
		return false
	}
	for i, _tgt := range p.References {
		_src5 := other.References[i]
		if !_tgt.Equals(_src5) {
			return false
		}
	}
	if p.Flags != other.Flags {
		return false
	}
	if p.StartTime != other.StartTime {
		return false
	}
	if p.Duration != other.Duration {
		return false
	}
	if len(p.Tags) != len(other.Tags) {
		return false
	}
	for i, _tgt := range p.Tags {
		_src6 := other.Tags[i]
		if !_tgt.Equals(_src6) {
			return false
		}
	}
	if len(p.Logs) != len(other.Logs) {
		return false
	}
	for i, _tgt := range p.Logs {
		_src7 := other.Logs[i]
		if !_tgt.Equals(_src7) {
			return false
		}
	}
	return true
}

func (p *Span) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Span(%+v)", *p)
}

// Attributes:
//  - ServiceName
//  - Tags
type Process struct {
	ServiceName string `thrift:"serviceName,1,required" db:"serviceName" json:"serviceName"`
	Tags        []*Tag `thrift:"tags,2" db:"tags" json:"tags,omitempty"`
}

func (p *Process) GetServiceName() string {
	return p.ServiceName
}

var Process_Tags_DEFAULT []*Tag

func (p *Process) GetTags() []*Tag {
	return p.Tags
}
func (p *Process) IsSetTags() bool {
	return p.Tags != nil
}

func (p *Process) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetServiceName bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
				issetServiceName = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 2:
			if fieldTypeId == thrift.LIST {
				if err := p.ReadField2(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetServiceName {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field ServiceName is not set"))
	}
	return nil
}

func (p *Process) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(ctx); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.ServiceName = v
	}
	return nil
}

func (p *Process) ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin(ctx)
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*Tag, 0, size)
	p.Tags = tSlice
	for i := 0; i < size; i++ {
		_elem8 := &Tag{}
		if err := _elem8.Read(ctx, iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem8), err)
		}
		p.Tags = append(p.Tags, _elem8)
	}
	if err := iprot.ReadListEnd(ctx); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *Process) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "Process"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField2(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Process) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "serviceName", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:serviceName: ", p), err)
	}
	if err := oprot.WriteString(ctx, string(p.ServiceName)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.serviceName (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:serviceName: ", p), err)
	}
	return err
}

func (p *Process) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if p.IsSetTags() {
		if err := oprot.WriteFieldBegin(ctx, "tags", thrift.LIST, 2); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:tags: ", p), err)
		}
		if err := oprot.WriteListBegin(ctx, thrift.STRUCT, len(p.Tags)); err != nil {
			return thrift.PrependError("error writing list begin: ", err)
		}
		for _, v := range p.Tags {
			if err := v.Write(ctx, oprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
			}
		}
		if err := oprot.WriteListEnd(ctx); err != nil {
			return thrift.PrependError("error writing list end: ", err)
		}
		if err := oprot.WriteFieldEnd(ctx); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 2:tags: ", p), err)
		}
	}
	return err
}

func (p *Process) Equals(other *Process) bool {
	if p == other {
		return true
	} else if p == nil || other == nil {
		return false
	}
	if p.ServiceName != other.ServiceName {
		return false
	}
	if len(p.Tags) != len(other.Tags) {
		return false
	}
	for i, _tgt := range p.Tags {
		_src9 := other.Tags[i]
		if !_tgt.Equals(_src9) {
			return false
		}
	}
	return true
}

func (p *Process) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Process(%+v)", *p)
}

// Attributes:
//  - FullQueueDroppedSpans
//  - TooLargeDroppedSpans
//  - FailedToEmitSpans
type ClientStats struct {
	FullQueueDroppedSpans int64 `thrift:"fullQueueDroppedSpans,1,required" db:"fullQueueDroppedSpans" json:"fullQueueDroppedSpans"`
	TooLargeDroppedSpans  int64 `thrift:"tooLargeDroppedSpans,2,required" db:"tooLargeDroppedSpans" json:"tooLargeDroppedSpans"`
	FailedToEmitSpans     int64 `thrift:"failedToEmitSpans,3,required" db:"failedToEmitSpans" json:"failedToEmitSpans"`
}

func (p *ClientStats) GetFullQueueDroppedSpans() int64 {
	return p.FullQueueDroppedSpans
}

func (p *ClientStats) GetTooLargeDroppedSpans() int64 {
	return p.TooLargeDroppedSpans
}

func (p *ClientStats) GetFailedToEmitSpans() int64 {
	return p.FailedToEmitSpans
}
func (p *ClientStats) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetFullQueueDroppedSpans bool = false
	var issetTooLargeDroppedSpans bool = false
	var issetFailedToEmitSpans bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
				issetFullQueueDroppedSpans = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 2:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField2(ctx, iprot); err != nil {
					return err
				}
				issetTooLargeDroppedSpans = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 3:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField3(ctx, iprot); err != nil {
					return err
				}
				issetFailedToEmitSpans = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetFullQueueDroppedSpans {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field FullQueueDroppedSpans is not set"))
	}
	if !issetTooLargeDroppedSpans {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field TooLargeDroppedSpans is not set"))
	}
	if !issetFailedToEmitSpans {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field FailedToEmitSpans is not set"))
	}
	return nil
}

func (p *ClientStats) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.FullQueueDroppedSpans = v
	}
	return nil
}

func (p *ClientStats) ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.TooLargeDroppedSpans = v
	}
	return nil
}

func (p *ClientStats) ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.FailedToEmitSpans = v
	}
	return nil
}

func (p *ClientStats) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "ClientStats"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField2(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField3(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *ClientStats) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "fullQueueDroppedSpans", thrift.I64, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:fullQueueDroppedSpans: ", p), err)
	}
	if err := oprot.WriteI64(ctx, int64(p.FullQueueDroppedSpans)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.fullQueueDroppedSpans (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:fullQueueDroppedSpans: ", p), err)
	}
	return err
}

func (p *ClientStats) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "tooLargeDroppedSpans", thrift.I64, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:tooLargeDroppedSpans: ", p), err)
	}
	if err := oprot.WriteI64(ctx, int64(p.TooLargeDroppedSpans)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.tooLargeDroppedSpans (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:tooLargeDroppedSpans: ", p), err)
	}
	return err
}

func (p *ClientStats) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "failedToEmitSpans", thrift.I64, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:failedToEmitSpans: ", p), err)
	}
	if err := oprot.WriteI64(ctx, int64(p.FailedToEmitSpans)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.failedToEmitSpans (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:failedToEmitSpans: ", p), err)
	}
	return err
}

func (p *ClientStats) Equals(other *ClientStats) bool {
	if p == other {
		return true
	} else if p == nil || other == nil {
		return false
	}
	if p.FullQueueDroppedSpans != other.FullQueueDroppedSpans {
		return false
	}
	if p.TooLargeDroppedSpans != other.TooLargeDroppedSpans {
		return false
	}
	if p.FailedToEmitSpans != other.FailedToEmitSpans {
		return false
	}
	return true
}

func (p *ClientStats) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ClientStats(%+v)", *p)
}

// Attributes:
//  - Process
//  - Spans
//  - SeqNo
//  - Stats
type Batch struct {
	Process *Process     `thrift:"process,1,required" db:"process" json:"process"`
	Spans   []*Span      `thrift:"spans,2,required" db:"spans" json:"spans"`
	SeqNo   *int64       `thrift:"seqNo,3" db:"seqNo" json:"seqNo,omitempty"`
	Stats   *ClientStats `thrift:"stats,4" db:"stats" json:"stats,omitempty"`
}

var Batch_Process_DEFAULT *Process

func (p *Batch) GetProcess() *Process {
	if !p.IsSetProcess() {
		return Batch_Process_DEFAULT
	}
	return p.Process
}

func (p *Batch) GetSpans() []*Span {
	return p.Spans
}

var Batch_SeqNo_DEFAULT int64

func (p *Batch) GetSeqNo() int64 {
	if !p.IsSetSeqNo() {
		return Batch_SeqNo_DEFAULT
	}
	return *p.SeqNo
}

var Batch_Stats_DEFAULT *ClientStats

func (p *Batch) GetStats() *ClientStats {
	if !p.IsSetStats() {
		return Batch_Stats_DEFAULT
	}
	return p.Stats
}
func (p *Batch) IsSetProcess() bool {
	return p.Process != nil
}

func (p *Batch) IsSetSeqNo() bool {
	return p.SeqNo != nil
}

func (p *Batch) IsSetStats() bool {
	return p.Stats != nil
}

func (p *Batch) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetProcess bool = false
	var issetSpans bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRUCT {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
				issetProcess = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 2:
			if fieldTypeId == thrift.LIST {
				if err := p.ReadField2(ctx, iprot); err != nil {
					return err
				}
				issetSpans = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 3:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField3(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 4:
			if fieldTypeId == thrift.STRUCT {
				if err := p.ReadField4(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetProcess {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Process is not set"))
	}
	if !issetSpans {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Spans is not set"))
	}
	return nil
}

func (p *Batch) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	p.Process = &Process{}
	if err := p.Process.Read(ctx, iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Process), err)
	}
	return nil
}

func (p *Batch) ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin(ctx)
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*Span, 0, size)
	p.Spans = tSlice
	for i := 0; i < size; i++ {
		_elem10 := &Span{}
		if err := _elem10.Read(ctx, iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem10), err)
		}
		p.Spans = append(p.Spans, _elem10)
	}
	if err := iprot.ReadListEnd(ctx); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *Batch) ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.SeqNo = &v
	}
	return nil
}

func (p *Batch) ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
	p.Stats = &ClientStats{}
	if err := p.Stats.Read(ctx, iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Stats), err)
	}
	return nil
}

func (p *Batch) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "Batch"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField2(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField3(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField4(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Batch) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "process", thrift.STRUCT, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:process: ", p), err)
	}
	if err := p.Process.Write(ctx, oprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Process), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:process: ", p), err)
	}
	return err
}

func (p *Batch) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "spans", thrift.LIST, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:spans: ", p), err)
	}
	if err := oprot.WriteListBegin(ctx, thrift.STRUCT, len(p.Spans)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.Spans {
		if err := v.Write(ctx, oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
		}
	}
	if err := oprot.WriteListEnd(ctx); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:spans: ", p), err)
	}
	return err
}

func (p *Batch) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if p.IsSetSeqNo() {
		if err := oprot.WriteFieldBegin(ctx, "seqNo", thrift.I64, 3); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:seqNo: ", p), err)
		}
		if err := oprot.WriteI64(ctx, int64(*p.SeqNo)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.seqNo (3) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(ctx); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 3:seqNo: ", p), err)
		}
	}
	return err
}

func (p *Batch) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if p.IsSetStats() {
		if err := oprot.WriteFieldBegin(ctx, "stats", thrift.STRUCT, 4); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:stats: ", p), err)
		}
		if err := p.Stats.Write(ctx, oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Stats), err)
		}
		if err := oprot.WriteFieldEnd(ctx); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 4:stats: ", p), err)
		}
	}
	return err
}

func (p *Batch) Equals(other *Batch) bool {
	if p == other {
		return true
	} else if p == nil || other == nil {
		return false
	}
	if !p.Process.Equals(other.Process) {
		return false
	}
	if len(p.Spans) != len(other.Spans) {
		return false
	}
	for i, _tgt := range p.Spans {
		_src11 := other.Spans[i]
		if !_tgt.Equals(_src11) {
			return false
		}
	}
	if p.SeqNo != other.SeqNo {
		if p.SeqNo == nil || other.SeqNo == nil {
			return false
		}
		if (*p.SeqNo) != (*other.SeqNo) {
			return false
		}
	}
	if !p.Stats.Equals(other.Stats) {
		return false
	}
	return true
}

func (p *Batch) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Batch(%+v)", *p)
}

// Attributes:
//  - Ok
type BatchSubmitResponse struct {
	Ok bool `thrift:"ok,1,required" db:"ok" json:"ok"`
}

func (p *BatchSubmitResponse) GetOk() bool {
	return p.Ok
}
func (p *BatchSubmitResponse) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetOk bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.BOOL {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
				issetOk = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetOk {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Ok is not set"))
	}
	return nil
}

func (p *BatchSubmitResponse) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(ctx); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Ok = v
	}
	return nil
}

func (p *BatchSubmitResponse) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "BatchSubmitResponse"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *BatchSubmitResponse) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "ok", thrift.BOOL, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:ok: ", p), err)
	}
	if err := oprot.WriteBool(ctx, bool(p.Ok)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.ok (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:ok: ", p), err)
	}
	return err
}

func (p *BatchSubmitResponse) Equals(other *BatchSubmitResponse) bool {
	if p == other {
		return true
	} else if p == nil || other == nil {
		return false
	}
	if p.Ok != other.Ok {
		return false
	}
	return true
}

func (p *BatchSubmitResponse) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BatchSubmitResponse(%+v)", *p)
}

type Collector interface {
	// Parameters:
	//  - Batches
	SubmitBatches(ctx context.Context, batches []*Batch) (_r []*BatchSubmitResponse, _err error)
}

type CollectorClient struct {
	c    thrift.TClient
	meta thrift.ResponseMeta
}

func (p *CollectorClient) Client_() thrift.TClient {
	return p.c
}

func (p *CollectorClient) LastResponseMeta_() thrift.ResponseMeta {
	return p.meta
}

func (p *CollectorClient) SetLastResponseMeta_(meta thrift.ResponseMeta) {
	p.meta = meta
}

// Parameters:
//  - Batches
func (p *CollectorClient) SubmitBatches(ctx context.Context, batches []*Batch) (_r []*BatchSubmitResponse, _err error) {
	var _args12 CollectorSubmitBatchesArgs
	_args12.Batches = batches
	var _result14 CollectorSubmitBatchesResult
	var _meta13 thrift.ResponseMeta
	_meta13, _err = p.Client_().Call(ctx, "submitBatches", &_args12, &_result14)
	p.SetLastResponseMeta_(_meta13)
	if _err != nil {
		return
	}
	return _result14.GetSuccess(), nil
}

type CollectorProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	handler      Collector
}

func (p *CollectorProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *CollectorProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

func (p *CollectorProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func (p *CollectorProcessor) Process(ctx context.Context, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqId, err2 := iprot.ReadMessageBegin(ctx)
	if err2 != nil {
		return false, thrift.WrapTException(err2)
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(ctx, seqId, iprot, oprot)
	}
	iprot.Skip(ctx, thrift.STRUCT)
	iprot.ReadMessageEnd(ctx)
	x16 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(ctx, name, thrift.EXCEPTION, seqId)
	x16.Write(ctx, oprot)
	oprot.WriteMessageEnd(ctx)
	oprot.Flush(ctx)
	return false, x16

}

// HELPER FUNCTIONS AND STRUCTURES

// Attributes:
//  - Batches
type CollectorSubmitBatchesArgs struct {
	Batches []*Batch `thrift:"batches,1" db:"batches" json:"batches"`
}

func (p *CollectorSubmitBatchesArgs) GetBatches() []*Batch {
	return p.Batches
}
func (p *CollectorSubmitBatchesArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.LIST {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *CollectorSubmitBatchesArgs) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin(ctx)
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*Batch, 0, size)
	p.Batches = tSlice
	for i := 0; i < size; i++ {
		_elem17 := &Batch{}
		if err := _elem17.Read(ctx, iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem17), err)
		}
		p.Batches = append(p.Batches, _elem17)
	}
	if err := iprot.ReadListEnd(ctx); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *CollectorSubmitBatchesArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "submitBatches_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *CollectorSubmitBatchesArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "batches", thrift.LIST, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:batches: ", p), err)
	}
	if err := oprot.WriteListBegin(ctx, thrift.STRUCT, len(p.Batches)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.Batches {
		if err := v.Write(ctx, oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
		}
	}
	if err := oprot.WriteListEnd(ctx); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:batches: ", p), err)
	}
	return err
}

func (p *CollectorSubmitBatchesArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("CollectorSubmitBatchesArgs(%+v)", *p)
}

// Attributes:
//  - Success
type CollectorSubmitBatchesResult struct {
	Success []*BatchSubmitResponse `thrift:"success,0" db:"success" json:"success,omitempty"`
}

func (p *CollectorSubmitBatchesResult) GetSuccess() []*BatchSubmitResponse {
	return p.Success
}
func (p *CollectorSubmitBatchesResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *CollectorSubmitBatchesResult) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if fieldTypeId == thrift.LIST {
				if err := p.ReadField0(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *CollectorSubmitBatchesResult) ReadField0(ctx context.Context, iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin(ctx)
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*BatchSubmitResponse, 0, size)
	p.Success = tSlice
	for i := 0; i < size; i++ {
		_elem18 := &BatchSubmitResponse{}
		if err := _elem18.Read(ctx, iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem18), err)
		}
		p.Success = append(p.Success, _elem18)
	}
	if err := iprot.ReadListEnd(ctx); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *CollectorSubmitBatchesResult) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "submitBatches_result"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField0(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *CollectorSubmitBatchesResult) writeField0(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin(ctx, "success", thrift.LIST, 0); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err)
		}
		if err := oprot.WriteListBegin(ctx, thrift.STRUCT, len(p.Success)); err != nil {
			return thrift.PrependError("error writing list begin: ", err)
		}
		for _, v := range p.Success {
			if err := v.Write(ctx, oprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
			}
		}
		if err := oprot.WriteListEnd(ctx); err != nil {
			return thrift.PrependError("error writing list end: ", err)
		}
		if err := oprot.WriteFieldEnd(ctx); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err)
		}
	}
	return err
}

func (p *CollectorSubmitBatchesResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("CollectorSubmitBatchesResult(%+v)", *p)
}
