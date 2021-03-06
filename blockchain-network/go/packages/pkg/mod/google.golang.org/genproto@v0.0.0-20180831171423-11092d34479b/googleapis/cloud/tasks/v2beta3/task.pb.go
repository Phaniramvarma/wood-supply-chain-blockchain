// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/tasks/v2beta3/task.proto

package tasks // import "google.golang.org/genproto/googleapis/cloud/tasks/v2beta3"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import status "google.golang.org/genproto/googleapis/rpc/status"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The view specifies a subset of [Task][google.cloud.tasks.v2beta3.Task] data.
//
// When a task is returned in a response, not all
// information is retrieved by default because some data, such as
// payloads, might be desirable to return only when needed because
// of its large size or because of the sensitivity of data that it
// contains.
type Task_View int32

const (
	// Unspecified. Defaults to BASIC.
	Task_VIEW_UNSPECIFIED Task_View = 0
	// The basic view omits fields which can be large or can contain
	// sensitive data.
	//
	// This view does not include the
	// [body in AppEngineHttpRequest][google.cloud.tasks.v2beta3.AppEngineHttpRequest.body].
	// Bodies are desirable to return only when needed, because they
	// can be large and because of the sensitivity of the data that you
	// choose to store in it.
	Task_BASIC Task_View = 1
	// All information is returned.
	//
	// Authorization for [FULL][google.cloud.tasks.v2beta3.Task.View.FULL] requires
	// `cloudtasks.tasks.fullView` [Google IAM](https://cloud.google.com/iam/)
	// permission on the [Queue][google.cloud.tasks.v2beta3.Queue] resource.
	Task_FULL Task_View = 2
)

var Task_View_name = map[int32]string{
	0: "VIEW_UNSPECIFIED",
	1: "BASIC",
	2: "FULL",
}
var Task_View_value = map[string]int32{
	"VIEW_UNSPECIFIED": 0,
	"BASIC":            1,
	"FULL":             2,
}

func (x Task_View) String() string {
	return proto.EnumName(Task_View_name, int32(x))
}
func (Task_View) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_task_949e38cc41c5c152, []int{0, 0}
}

// A unit of scheduled work.
type Task struct {
	// Optionally caller-specified in [CreateTask][google.cloud.tasks.v2beta3.CloudTasks.CreateTask].
	//
	// The task name.
	//
	// The task name must have the following format:
	// `projects/PROJECT_ID/locations/LOCATION_ID/queues/QUEUE_ID/tasks/TASK_ID`
	//
	// * `PROJECT_ID` can contain letters ([A-Za-z]), numbers ([0-9]),
	//    hyphens (-), colons (:), or periods (.).
	//    For more information, see
	//    [Identifying projects](https://cloud.google.com/resource-manager/docs/creating-managing-projects#identifying_projects)
	// * `LOCATION_ID` is the canonical ID for the task's location.
	//    The list of available locations can be obtained by calling
	//    [ListLocations][google.cloud.location.Locations.ListLocations].
	//    For more information, see https://cloud.google.com/about/locations/.
	// * `QUEUE_ID` can contain letters ([A-Za-z]), numbers ([0-9]), or
	//   hyphens (-). The maximum length is 100 characters.
	// * `TASK_ID` can contain only letters ([A-Za-z]), numbers ([0-9]),
	//   hyphens (-), or underscores (_). The maximum length is 500 characters.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Required.
	//
	// The task's payload is used by the task's target to process the task.
	// A payload is valid only if it is compatible with the queue's target.
	//
	// Types that are valid to be assigned to PayloadType:
	//	*Task_AppEngineHttpRequest
	PayloadType isTask_PayloadType `protobuf_oneof:"payload_type"`
	// The time when the task is scheduled to be attempted.
	//
	// For App Engine queues, this is when the task will be attempted or retried.
	//
	// `schedule_time` will be truncated to the nearest microsecond.
	ScheduleTime *timestamp.Timestamp `protobuf:"bytes,4,opt,name=schedule_time,json=scheduleTime,proto3" json:"schedule_time,omitempty"`
	// Output only. The time that the task was created.
	//
	// `create_time` will be truncated to the nearest second.
	CreateTime *timestamp.Timestamp `protobuf:"bytes,5,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// Output only. The number of attempts dispatched.
	//
	// This count includes tasks which have been dispatched but haven't
	// received a response.
	DispatchCount int32 `protobuf:"varint,6,opt,name=dispatch_count,json=dispatchCount,proto3" json:"dispatch_count,omitempty"`
	// Output only. The number of attempts which have received a response.
	ResponseCount int32 `protobuf:"varint,7,opt,name=response_count,json=responseCount,proto3" json:"response_count,omitempty"`
	// Output only. The status of the task's first attempt.
	//
	// Only [dispatch_time][google.cloud.tasks.v2beta3.Attempt.dispatch_time] will be set.
	// The other [Attempt][google.cloud.tasks.v2beta3.Attempt] information is not retained by Cloud Tasks.
	FirstAttempt *Attempt `protobuf:"bytes,8,opt,name=first_attempt,json=firstAttempt,proto3" json:"first_attempt,omitempty"`
	// Output only. The status of the task's last attempt.
	LastAttempt *Attempt `protobuf:"bytes,9,opt,name=last_attempt,json=lastAttempt,proto3" json:"last_attempt,omitempty"`
	// Output only. The view specifies which subset of the [Task][google.cloud.tasks.v2beta3.Task] has
	// been returned.
	View                 Task_View `protobuf:"varint,10,opt,name=view,proto3,enum=google.cloud.tasks.v2beta3.Task_View" json:"view,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Task) Reset()         { *m = Task{} }
func (m *Task) String() string { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()    {}
func (*Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_task_949e38cc41c5c152, []int{0}
}
func (m *Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Task.Unmarshal(m, b)
}
func (m *Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Task.Marshal(b, m, deterministic)
}
func (dst *Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Task.Merge(dst, src)
}
func (m *Task) XXX_Size() int {
	return xxx_messageInfo_Task.Size(m)
}
func (m *Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Task proto.InternalMessageInfo

func (m *Task) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type isTask_PayloadType interface {
	isTask_PayloadType()
}

type Task_AppEngineHttpRequest struct {
	AppEngineHttpRequest *AppEngineHttpRequest `protobuf:"bytes,3,opt,name=app_engine_http_request,json=appEngineHttpRequest,proto3,oneof"`
}

func (*Task_AppEngineHttpRequest) isTask_PayloadType() {}

func (m *Task) GetPayloadType() isTask_PayloadType {
	if m != nil {
		return m.PayloadType
	}
	return nil
}

func (m *Task) GetAppEngineHttpRequest() *AppEngineHttpRequest {
	if x, ok := m.GetPayloadType().(*Task_AppEngineHttpRequest); ok {
		return x.AppEngineHttpRequest
	}
	return nil
}

func (m *Task) GetScheduleTime() *timestamp.Timestamp {
	if m != nil {
		return m.ScheduleTime
	}
	return nil
}

func (m *Task) GetCreateTime() *timestamp.Timestamp {
	if m != nil {
		return m.CreateTime
	}
	return nil
}

func (m *Task) GetDispatchCount() int32 {
	if m != nil {
		return m.DispatchCount
	}
	return 0
}

func (m *Task) GetResponseCount() int32 {
	if m != nil {
		return m.ResponseCount
	}
	return 0
}

func (m *Task) GetFirstAttempt() *Attempt {
	if m != nil {
		return m.FirstAttempt
	}
	return nil
}

func (m *Task) GetLastAttempt() *Attempt {
	if m != nil {
		return m.LastAttempt
	}
	return nil
}

func (m *Task) GetView() Task_View {
	if m != nil {
		return m.View
	}
	return Task_VIEW_UNSPECIFIED
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Task) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Task_OneofMarshaler, _Task_OneofUnmarshaler, _Task_OneofSizer, []interface{}{
		(*Task_AppEngineHttpRequest)(nil),
	}
}

func _Task_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Task)
	// payload_type
	switch x := m.PayloadType.(type) {
	case *Task_AppEngineHttpRequest:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.AppEngineHttpRequest); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Task.PayloadType has unexpected type %T", x)
	}
	return nil
}

func _Task_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Task)
	switch tag {
	case 3: // payload_type.app_engine_http_request
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(AppEngineHttpRequest)
		err := b.DecodeMessage(msg)
		m.PayloadType = &Task_AppEngineHttpRequest{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Task_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Task)
	// payload_type
	switch x := m.PayloadType.(type) {
	case *Task_AppEngineHttpRequest:
		s := proto.Size(x.AppEngineHttpRequest)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// The status of a task attempt.
type Attempt struct {
	// Output only. The time that this attempt was scheduled.
	//
	// `schedule_time` will be truncated to the nearest microsecond.
	ScheduleTime *timestamp.Timestamp `protobuf:"bytes,1,opt,name=schedule_time,json=scheduleTime,proto3" json:"schedule_time,omitempty"`
	// Output only. The time that this attempt was dispatched.
	//
	// `dispatch_time` will be truncated to the nearest microsecond.
	DispatchTime *timestamp.Timestamp `protobuf:"bytes,2,opt,name=dispatch_time,json=dispatchTime,proto3" json:"dispatch_time,omitempty"`
	// Output only. The time that this attempt response was received.
	//
	// `response_time` will be truncated to the nearest microsecond.
	ResponseTime *timestamp.Timestamp `protobuf:"bytes,3,opt,name=response_time,json=responseTime,proto3" json:"response_time,omitempty"`
	// Output only. The response from the target for this attempt.
	//
	// If `response_time` is unset, then the task has not been attempted or is
	// currently running and the `response_status` field is meaningless.
	ResponseStatus       *status.Status `protobuf:"bytes,4,opt,name=response_status,json=responseStatus,proto3" json:"response_status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Attempt) Reset()         { *m = Attempt{} }
func (m *Attempt) String() string { return proto.CompactTextString(m) }
func (*Attempt) ProtoMessage()    {}
func (*Attempt) Descriptor() ([]byte, []int) {
	return fileDescriptor_task_949e38cc41c5c152, []int{1}
}
func (m *Attempt) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Attempt.Unmarshal(m, b)
}
func (m *Attempt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Attempt.Marshal(b, m, deterministic)
}
func (dst *Attempt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Attempt.Merge(dst, src)
}
func (m *Attempt) XXX_Size() int {
	return xxx_messageInfo_Attempt.Size(m)
}
func (m *Attempt) XXX_DiscardUnknown() {
	xxx_messageInfo_Attempt.DiscardUnknown(m)
}

var xxx_messageInfo_Attempt proto.InternalMessageInfo

func (m *Attempt) GetScheduleTime() *timestamp.Timestamp {
	if m != nil {
		return m.ScheduleTime
	}
	return nil
}

func (m *Attempt) GetDispatchTime() *timestamp.Timestamp {
	if m != nil {
		return m.DispatchTime
	}
	return nil
}

func (m *Attempt) GetResponseTime() *timestamp.Timestamp {
	if m != nil {
		return m.ResponseTime
	}
	return nil
}

func (m *Attempt) GetResponseStatus() *status.Status {
	if m != nil {
		return m.ResponseStatus
	}
	return nil
}

func init() {
	proto.RegisterType((*Task)(nil), "google.cloud.tasks.v2beta3.Task")
	proto.RegisterType((*Attempt)(nil), "google.cloud.tasks.v2beta3.Attempt")
	proto.RegisterEnum("google.cloud.tasks.v2beta3.Task_View", Task_View_name, Task_View_value)
}

func init() {
	proto.RegisterFile("google/cloud/tasks/v2beta3/task.proto", fileDescriptor_task_949e38cc41c5c152)
}

var fileDescriptor_task_949e38cc41c5c152 = []byte{
	// 538 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x94, 0xcd, 0x6e, 0xda, 0x40,
	0x10, 0xc7, 0xe3, 0xc4, 0xf9, 0x60, 0x03, 0x14, 0xad, 0x22, 0xc5, 0x42, 0x55, 0x8b, 0xa8, 0x50,
	0x39, 0xd9, 0x2d, 0x39, 0x55, 0x39, 0xa0, 0x40, 0x41, 0x20, 0x45, 0x15, 0x32, 0x49, 0x2a, 0xf5,
	0x62, 0x2d, 0x66, 0x63, 0xac, 0x98, 0xdd, 0xad, 0x77, 0x9c, 0x28, 0x8f, 0xd0, 0xc7, 0xec, 0x9b,
	0x54, 0xde, 0x0f, 0x54, 0xa9, 0x29, 0xb4, 0x37, 0xcf, 0xec, 0xef, 0xff, 0x9f, 0x59, 0xcf, 0xd8,
	0xa8, 0x93, 0x70, 0x9e, 0x64, 0x34, 0x88, 0x33, 0x5e, 0x2c, 0x03, 0x20, 0xf2, 0x41, 0x06, 0x8f,
	0xbd, 0x05, 0x05, 0x72, 0xa1, 0x22, 0x5f, 0xe4, 0x1c, 0x38, 0x6e, 0x6a, 0xcc, 0x57, 0x98, 0xaf,
	0x30, 0xdf, 0x60, 0xcd, 0xd7, 0xc6, 0x82, 0x88, 0x34, 0x20, 0x8c, 0x71, 0x20, 0x90, 0x72, 0x26,
	0xb5, 0xb2, 0xf9, 0x7e, 0x6b, 0x81, 0x3c, 0xa1, 0x60, 0xc0, 0xb7, 0x06, 0x54, 0xd1, 0xa2, 0xb8,
	0x0f, 0x20, 0x5d, 0x53, 0x09, 0x64, 0x2d, 0x0c, 0x70, 0x6e, 0x80, 0x5c, 0xc4, 0x81, 0x04, 0x02,
	0x85, 0x29, 0xd1, 0xfe, 0xe9, 0x22, 0xf7, 0x86, 0xc8, 0x07, 0x8c, 0x91, 0xcb, 0xc8, 0x9a, 0x7a,
	0x4e, 0xcb, 0xe9, 0x56, 0x42, 0xf5, 0x8c, 0x53, 0x74, 0x4e, 0x84, 0x88, 0x28, 0x4b, 0x52, 0x46,
	0xa3, 0x15, 0x80, 0x88, 0x72, 0xfa, 0xbd, 0xa0, 0x12, 0xbc, 0x83, 0x96, 0xd3, 0x3d, 0xed, 0x7d,
	0xf0, 0xff, 0x7e, 0x37, 0xff, 0x4a, 0x88, 0x91, 0x52, 0x4e, 0x00, 0x44, 0xa8, 0x75, 0x93, 0xbd,
	0xf0, 0x8c, 0xbc, 0x90, 0xc7, 0x7d, 0x54, 0x93, 0xf1, 0x8a, 0x2e, 0x8b, 0x8c, 0x46, 0x65, 0xf3,
	0x9e, 0xab, 0x0a, 0x34, 0x6d, 0x01, 0x7b, 0x33, 0xff, 0xc6, 0xde, 0x2c, 0xac, 0x5a, 0x41, 0x99,
	0xc2, 0x97, 0xe8, 0x34, 0xce, 0x29, 0x01, 0x23, 0x3f, 0xdc, 0x29, 0x47, 0x1a, 0x57, 0xe2, 0x0e,
	0xaa, 0x2f, 0x53, 0x29, 0x08, 0xc4, 0xab, 0x28, 0xe6, 0x05, 0x03, 0xef, 0xa8, 0xe5, 0x74, 0x0f,
	0xc3, 0x9a, 0xcd, 0x0e, 0xcb, 0x64, 0x89, 0xe5, 0x54, 0x0a, 0xce, 0x24, 0x35, 0xd8, 0xb1, 0xc6,
	0x6c, 0x56, 0x63, 0x13, 0x54, 0xbb, 0x4f, 0x73, 0x09, 0x11, 0x01, 0xa0, 0x6b, 0x01, 0xde, 0x89,
	0x6a, 0xe6, 0xdd, 0xd6, 0x97, 0xa5, 0xd1, 0xb0, 0xaa, 0x94, 0x26, 0xc2, 0x63, 0x54, 0xcd, 0xc8,
	0x6f, 0x46, 0x95, 0x7f, 0x37, 0x3a, 0x2d, 0x85, 0xd6, 0xe7, 0x13, 0x72, 0x1f, 0x53, 0xfa, 0xe4,
	0xa1, 0x96, 0xd3, 0xad, 0xf7, 0x3a, 0xdb, 0xf4, 0xe5, 0x32, 0xf8, 0x77, 0x29, 0x7d, 0x0a, 0x95,
	0xa4, 0xfd, 0x11, 0xb9, 0x65, 0x84, 0xcf, 0x50, 0xe3, 0x6e, 0x3a, 0xfa, 0x1a, 0xdd, 0x7e, 0x99,
	0xcf, 0x46, 0xc3, 0xe9, 0x78, 0x3a, 0xfa, 0xdc, 0xd8, 0xc3, 0x15, 0x74, 0x38, 0xb8, 0x9a, 0x4f,
	0x87, 0x0d, 0x07, 0x9f, 0x20, 0x77, 0x7c, 0x7b, 0x7d, 0xdd, 0xd8, 0x1f, 0xd4, 0x51, 0x55, 0x90,
	0xe7, 0x8c, 0x93, 0x65, 0x04, 0xcf, 0x82, 0xb6, 0x7f, 0xec, 0xa3, 0x63, 0xdb, 0xc9, 0x1f, 0x73,
	0x76, 0xfe, 0x73, 0xce, 0x7d, 0xb4, 0x19, 0x8a, 0x36, 0xd8, 0xdf, 0x6d, 0x60, 0x05, 0xd6, 0x60,
	0x33, 0x44, 0x65, 0x70, 0xb0, 0xdb, 0xc0, 0x0a, 0xcc, 0xa6, 0xbd, 0xda, 0x18, 0xe8, 0x6f, 0xc9,
	0x2c, 0x2b, 0xb6, 0x16, 0xb9, 0x88, 0xfd, 0xb9, 0x3a, 0x09, 0x37, 0x0b, 0xa3, 0xe3, 0x01, 0x43,
	0x6f, 0x62, 0xbe, 0xde, 0x32, 0x80, 0x41, 0xa5, 0x9c, 0xc0, 0xac, 0x6c, 0x62, 0xe6, 0x7c, 0xeb,
	0x1b, 0x30, 0xe1, 0x19, 0x61, 0x89, 0xcf, 0xf3, 0x24, 0x48, 0x28, 0x53, 0x2d, 0x06, 0xfa, 0x88,
	0x88, 0x54, 0xbe, 0xf4, 0x83, 0xb8, 0x54, 0xd1, 0xe2, 0x48, 0xb1, 0x17, 0xbf, 0x02, 0x00, 0x00,
	0xff, 0xff, 0xa7, 0x0f, 0xef, 0xad, 0xac, 0x04, 0x00, 0x00,
}
