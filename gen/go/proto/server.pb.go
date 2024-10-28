// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: proto/server.proto

package choicesv1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Enumeration of possible user choices.
type Winner int32

const (
	// Default unspecified value.
	Winner_UNSPECIFIED Winner = 0
	// User didn't like either joke.
	Winner_NONE Winner = 1
	// User chose the left joke.
	Winner_LEFT Winner = 2
	// User chose the right joke.
	Winner_RIGHT Winner = 3
	// User liked both jokes equally.
	Winner_BOTH Winner = 4
)

// Enum value maps for Winner.
var (
	Winner_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "NONE",
		2: "LEFT",
		3: "RIGHT",
		4: "BOTH",
	}
	Winner_value = map[string]int32{
		"UNSPECIFIED": 0,
		"NONE":        1,
		"LEFT":        2,
		"RIGHT":       3,
		"BOTH":        4,
	}
)

func (x Winner) Enum() *Winner {
	p := new(Winner)
	*p = x
	return p
}

func (x Winner) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Winner) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_server_proto_enumTypes[0].Descriptor()
}

func (Winner) Type() protoreflect.EnumType {
	return &file_proto_server_proto_enumTypes[0]
}

func (x Winner) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Winner.Descriptor instead.
func (Winner) EnumDescriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{0}
}

type GetChoicesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// TODO(rbtz): allow passing theme for jokes
	// A UUID for the user's session.
	SessionId string `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
}

func (x *GetChoicesRequest) Reset() {
	*x = GetChoicesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChoicesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChoicesRequest) ProtoMessage() {}

func (x *GetChoicesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChoicesRequest.ProtoReflect.Descriptor instead.
func (*GetChoicesRequest) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{0}
}

func (x *GetChoicesRequest) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

type GetChoicesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier for the joke pair.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Theme for both jokes.
	Theme string `protobuf:"bytes,2,opt,name=theme,proto3" json:"theme,omitempty"`
	// Text of the left joke.
	LeftJoke string `protobuf:"bytes,3,opt,name=left_joke,json=leftJoke,proto3" json:"left_joke,omitempty"`
	// Text of the right joke.
	RightJoke string `protobuf:"bytes,4,opt,name=right_joke,json=rightJoke,proto3" json:"right_joke,omitempty"`
}

func (x *GetChoicesResponse) Reset() {
	*x = GetChoicesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChoicesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChoicesResponse) ProtoMessage() {}

func (x *GetChoicesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChoicesResponse.ProtoReflect.Descriptor instead.
func (*GetChoicesResponse) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{1}
}

func (x *GetChoicesResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GetChoicesResponse) GetTheme() string {
	if x != nil {
		return x.Theme
	}
	return ""
}

func (x *GetChoicesResponse) GetLeftJoke() string {
	if x != nil {
		return x.LeftJoke
	}
	return ""
}

func (x *GetChoicesResponse) GetRightJoke() string {
	if x != nil {
		return x.RightJoke
	}
	return ""
}

// Request to rate the presented jokes.
type RateChoicesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Identifier of the joke pair being rated.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The user's selection.
	Winner Winner `protobuf:"varint,2,opt,name=winner,proto3,enum=choices.v1.Winner" json:"winner,omitempty"`
	// Known jokes.
	Known Winner `protobuf:"varint,3,opt,name=known,proto3,enum=choices.v1.Winner" json:"known,omitempty"`
}

func (x *RateChoicesRequest) Reset() {
	*x = RateChoicesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateChoicesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateChoicesRequest) ProtoMessage() {}

func (x *RateChoicesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateChoicesRequest.ProtoReflect.Descriptor instead.
func (*RateChoicesRequest) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{2}
}

func (x *RateChoicesRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RateChoicesRequest) GetWinner() Winner {
	if x != nil {
		return x.Winner
	}
	return Winner_UNSPECIFIED
}

func (x *RateChoicesRequest) GetKnown() Winner {
	if x != nil {
		return x.Known
	}
	return Winner_UNSPECIFIED
}

type RateChoicesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RateChoicesResponse) Reset() {
	*x = RateChoicesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateChoicesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateChoicesResponse) ProtoMessage() {}

func (x *RateChoicesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateChoicesResponse.ProtoReflect.Descriptor instead.
func (*RateChoicesResponse) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{3}
}

// Request to get the Arena leaderboard.
type GetLeaderboardRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetLeaderboardRequest) Reset() {
	*x = GetLeaderboardRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLeaderboardRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLeaderboardRequest) ProtoMessage() {}

func (x *GetLeaderboardRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLeaderboardRequest.ProtoReflect.Descriptor instead.
func (*GetLeaderboardRequest) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{4}
}

// LeaderboardEntry contains the model name and its Bradley-Terry rating.
type LeaderboardEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Public model name.
	Model string `protobuf:"bytes,1,opt,name=model,proto3" json:"model,omitempty"`
	// Newman Score of the model.
	NewmanScore float64 `protobuf:"fixed64,2,opt,name=newmanScore,proto3" json:"newmanScore,omitempty"`
	// Total votes for the model.
	Votes uint64 `protobuf:"varint,3,opt,name=votes,proto3" json:"votes,omitempty"`
}

func (x *LeaderboardEntry) Reset() {
	*x = LeaderboardEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeaderboardEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeaderboardEntry) ProtoMessage() {}

func (x *LeaderboardEntry) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeaderboardEntry.ProtoReflect.Descriptor instead.
func (*LeaderboardEntry) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{5}
}

func (x *LeaderboardEntry) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *LeaderboardEntry) GetNewmanScore() float64 {
	if x != nil {
		return x.NewmanScore
	}
	return 0
}

func (x *LeaderboardEntry) GetVotes() uint64 {
	if x != nil {
		return x.Votes
	}
	return 0
}

// GetLeaderboardResponse contains the leaderboard of joke models.
type GetLeaderboardResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entries []*LeaderboardEntry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *GetLeaderboardResponse) Reset() {
	*x = GetLeaderboardResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLeaderboardResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLeaderboardResponse) ProtoMessage() {}

func (x *GetLeaderboardResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLeaderboardResponse.ProtoReflect.Descriptor instead.
func (*GetLeaderboardResponse) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{6}
}

func (x *GetLeaderboardResponse) GetEntries() []*LeaderboardEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

var File_proto_server_proto protoreflect.FileDescriptor

var file_proto_server_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x76, 0x31,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x32, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x49, 0x64, 0x22, 0x76, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x68, 0x65,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x12,
	0x1b, 0x0a, 0x09, 0x6c, 0x65, 0x66, 0x74, 0x5f, 0x6a, 0x6f, 0x6b, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6c, 0x65, 0x66, 0x74, 0x4a, 0x6f, 0x6b, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x72, 0x69, 0x67, 0x68, 0x74, 0x5f, 0x6a, 0x6f, 0x6b, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x72, 0x69, 0x67, 0x68, 0x74, 0x4a, 0x6f, 0x6b, 0x65, 0x22, 0x7f, 0x0a, 0x12, 0x52,
	0x61, 0x74, 0x65, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x13, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0,
	0x41, 0x02, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2a, 0x0a, 0x06, 0x77, 0x69, 0x6e, 0x6e, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x57, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x52, 0x06, 0x77, 0x69, 0x6e, 0x6e,
	0x65, 0x72, 0x12, 0x28, 0x0a, 0x05, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x12, 0x2e, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x57,
	0x69, 0x6e, 0x6e, 0x65, 0x72, 0x52, 0x05, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x22, 0x15, 0x0a, 0x13,
	0x52, 0x61, 0x74, 0x65, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x17, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x62, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x60, 0x0a, 0x10,
	0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x65, 0x77, 0x6d, 0x61, 0x6e,
	0x53, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x6e, 0x65, 0x77,
	0x6d, 0x61, 0x6e, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x6f, 0x74, 0x65,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x6f, 0x74, 0x65, 0x73, 0x22, 0x50,
	0x0a, 0x16, 0x47, 0x65, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72,
	0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x68, 0x6f, 0x69,
	0x63, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61,
	0x72, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x2a, 0x42, 0x0a, 0x06, 0x57, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x4e,
	0x4f, 0x4e, 0x45, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x4c, 0x45, 0x46, 0x54, 0x10, 0x02, 0x12,
	0x09, 0x0a, 0x05, 0x52, 0x49, 0x47, 0x48, 0x54, 0x10, 0x03, 0x12, 0x08, 0x0a, 0x04, 0x42, 0x4f,
	0x54, 0x48, 0x10, 0x04, 0x32, 0xcb, 0x02, 0x0a, 0x05, 0x41, 0x72, 0x65, 0x6e, 0x61, 0x12, 0x5f,
	0x0a, 0x0a, 0x47, 0x65, 0x74, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x12, 0x1d, 0x2e, 0x63,
	0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x68, 0x6f,
	0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x63, 0x68,
	0x6f, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x68, 0x6f, 0x69,
	0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x12, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x12,
	0x6f, 0x0a, 0x0b, 0x52, 0x61, 0x74, 0x65, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x12, 0x1e,
	0x2e, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61, 0x74, 0x65,
	0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f,
	0x2e, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61, 0x74, 0x65,
	0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x1f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x3a, 0x01, 0x2a, 0x22, 0x14, 0x2f, 0x76, 0x31, 0x2f,
	0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x2f, 0x72, 0x61, 0x74, 0x65,
	0x12, 0x70, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61,
	0x72, 0x64, 0x12, 0x21, 0x2e, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x11, 0x12, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61,
	0x72, 0x64, 0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x53, 0x61, 0x76, 0x65, 0x54, 0x68, 0x65, 0x52, 0x62, 0x74, 0x7a, 0x2f, 0x68, 0x75, 0x6d,
	0x6f, 0x72, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65,
	0x73, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x76, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_server_proto_rawDescOnce sync.Once
	file_proto_server_proto_rawDescData = file_proto_server_proto_rawDesc
)

func file_proto_server_proto_rawDescGZIP() []byte {
	file_proto_server_proto_rawDescOnce.Do(func() {
		file_proto_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_server_proto_rawDescData)
	})
	return file_proto_server_proto_rawDescData
}

var file_proto_server_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_server_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_proto_server_proto_goTypes = []any{
	(Winner)(0),                    // 0: choices.v1.Winner
	(*GetChoicesRequest)(nil),      // 1: choices.v1.GetChoicesRequest
	(*GetChoicesResponse)(nil),     // 2: choices.v1.GetChoicesResponse
	(*RateChoicesRequest)(nil),     // 3: choices.v1.RateChoicesRequest
	(*RateChoicesResponse)(nil),    // 4: choices.v1.RateChoicesResponse
	(*GetLeaderboardRequest)(nil),  // 5: choices.v1.GetLeaderboardRequest
	(*LeaderboardEntry)(nil),       // 6: choices.v1.LeaderboardEntry
	(*GetLeaderboardResponse)(nil), // 7: choices.v1.GetLeaderboardResponse
}
var file_proto_server_proto_depIdxs = []int32{
	0, // 0: choices.v1.RateChoicesRequest.winner:type_name -> choices.v1.Winner
	0, // 1: choices.v1.RateChoicesRequest.known:type_name -> choices.v1.Winner
	6, // 2: choices.v1.GetLeaderboardResponse.entries:type_name -> choices.v1.LeaderboardEntry
	1, // 3: choices.v1.Arena.GetChoices:input_type -> choices.v1.GetChoicesRequest
	3, // 4: choices.v1.Arena.RateChoices:input_type -> choices.v1.RateChoicesRequest
	5, // 5: choices.v1.Arena.GetLeaderboard:input_type -> choices.v1.GetLeaderboardRequest
	2, // 6: choices.v1.Arena.GetChoices:output_type -> choices.v1.GetChoicesResponse
	4, // 7: choices.v1.Arena.RateChoices:output_type -> choices.v1.RateChoicesResponse
	7, // 8: choices.v1.Arena.GetLeaderboard:output_type -> choices.v1.GetLeaderboardResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_server_proto_init() }
func file_proto_server_proto_init() {
	if File_proto_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_server_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GetChoicesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_server_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*GetChoicesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_server_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*RateChoicesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_server_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*RateChoicesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_server_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*GetLeaderboardRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_server_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*LeaderboardEntry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_server_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*GetLeaderboardResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_server_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_server_proto_goTypes,
		DependencyIndexes: file_proto_server_proto_depIdxs,
		EnumInfos:         file_proto_server_proto_enumTypes,
		MessageInfos:      file_proto_server_proto_msgTypes,
	}.Build()
	File_proto_server_proto = out.File
	file_proto_server_proto_rawDesc = nil
	file_proto_server_proto_goTypes = nil
	file_proto_server_proto_depIdxs = nil
}
