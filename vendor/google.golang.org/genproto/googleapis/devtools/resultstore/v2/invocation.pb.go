// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/devtools/resultstore/v2/invocation.proto

package resultstore

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// An Invocation typically represents the result of running a tool. Each has a
// unique ID, typically generated by the server. Target resources under each
// Invocation contain the bulk of the data.
type Invocation struct {
	// The resource name.  Its format must be:
	// invocations/${INVOCATION_ID}
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The resource ID components that identify the Invocation. They must match
	// the resource name after proper encoding.
	Id *Invocation_Id `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	// The aggregate status of the invocation.
	StatusAttributes *StatusAttributes `protobuf:"bytes,3,opt,name=status_attributes,json=statusAttributes,proto3" json:"status_attributes,omitempty"`
	// When this invocation started and its duration.
	Timing *Timing `protobuf:"bytes,4,opt,name=timing,proto3" json:"timing,omitempty"`
	// Attributes of this invocation.
	InvocationAttributes *InvocationAttributes `protobuf:"bytes,5,opt,name=invocation_attributes,json=invocationAttributes,proto3" json:"invocation_attributes,omitempty"`
	// The workspace the tool was run in.
	WorkspaceInfo *WorkspaceInfo `protobuf:"bytes,6,opt,name=workspace_info,json=workspaceInfo,proto3" json:"workspace_info,omitempty"`
	// Arbitrary name-value pairs.
	// This is implemented as a multi-map. Multiple properties are allowed with
	// the same key. Properties will be returned in lexicographical order by key.
	Properties []*Property `protobuf:"bytes,7,rep,name=properties,proto3" json:"properties,omitempty"`
	// A list of file references for invocation level files.
	// The file IDs must be unique within this list. Duplicate file IDs will
	// result in an error. Files will be returned in lexicographical order by ID.
	// Use this field to specify build logs, and other invocation level logs.
	//
	// Files with the following reserved file IDs cause specific post-processing
	// or have special handling. These files must be immediately available to
	// ResultStore for processing when the reference is uploaded.
	//
	// build.log: The primary log for the Invocation.
	// coverage_report.lcov: Aggregate coverage report for the invocation.
	Files []*File `protobuf:"bytes,8,rep,name=files,proto3" json:"files,omitempty"`
	// Summary of aggregate coverage across all Actions in this Invocation.
	// If missing, this data will be populated by the server from the
	// coverage_report.lcov file or the union of all ActionCoverages under this
	// invocation (in that order).
	CoverageSummaries []*LanguageCoverageSummary `protobuf:"bytes,9,rep,name=coverage_summaries,json=coverageSummaries,proto3" json:"coverage_summaries,omitempty"`
	// Aggregate code coverage for all build and test Actions within this
	// Invocation. If missing, this data will be populated by the server
	// from the coverage_report.lcov file or the union of all ActionCoverages
	// under this invocation (in that order).
	AggregateCoverage *AggregateCoverage `protobuf:"bytes,10,opt,name=aggregate_coverage,json=aggregateCoverage,proto3" json:"aggregate_coverage,omitempty"`
	// NOT IMPLEMENTED.
	// ResultStore will read and parse Files with reserved IDs listed above. Read
	// and parse errors for all these Files are reported here.
	// This is implemented as a map, with one FileProcessingErrors for each file.
	// Typically produced when parsing Files, but may also be provided directly
	// by clients.
	FileProcessingErrors []*FileProcessingErrors `protobuf:"bytes,11,rep,name=file_processing_errors,json=fileProcessingErrors,proto3" json:"file_processing_errors,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *Invocation) Reset()         { *m = Invocation{} }
func (m *Invocation) String() string { return proto.CompactTextString(m) }
func (*Invocation) ProtoMessage()    {}
func (*Invocation) Descriptor() ([]byte, []int) {
	return fileDescriptor_9a57b6cf1112b76d, []int{0}
}

func (m *Invocation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Invocation.Unmarshal(m, b)
}
func (m *Invocation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Invocation.Marshal(b, m, deterministic)
}
func (m *Invocation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Invocation.Merge(m, src)
}
func (m *Invocation) XXX_Size() int {
	return xxx_messageInfo_Invocation.Size(m)
}
func (m *Invocation) XXX_DiscardUnknown() {
	xxx_messageInfo_Invocation.DiscardUnknown(m)
}

var xxx_messageInfo_Invocation proto.InternalMessageInfo

func (m *Invocation) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Invocation) GetId() *Invocation_Id {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Invocation) GetStatusAttributes() *StatusAttributes {
	if m != nil {
		return m.StatusAttributes
	}
	return nil
}

func (m *Invocation) GetTiming() *Timing {
	if m != nil {
		return m.Timing
	}
	return nil
}

func (m *Invocation) GetInvocationAttributes() *InvocationAttributes {
	if m != nil {
		return m.InvocationAttributes
	}
	return nil
}

func (m *Invocation) GetWorkspaceInfo() *WorkspaceInfo {
	if m != nil {
		return m.WorkspaceInfo
	}
	return nil
}

func (m *Invocation) GetProperties() []*Property {
	if m != nil {
		return m.Properties
	}
	return nil
}

func (m *Invocation) GetFiles() []*File {
	if m != nil {
		return m.Files
	}
	return nil
}

func (m *Invocation) GetCoverageSummaries() []*LanguageCoverageSummary {
	if m != nil {
		return m.CoverageSummaries
	}
	return nil
}

func (m *Invocation) GetAggregateCoverage() *AggregateCoverage {
	if m != nil {
		return m.AggregateCoverage
	}
	return nil
}

func (m *Invocation) GetFileProcessingErrors() []*FileProcessingErrors {
	if m != nil {
		return m.FileProcessingErrors
	}
	return nil
}

// The resource ID components that identify the Invocation.
type Invocation_Id struct {
	// The Invocation ID.
	InvocationId         string   `protobuf:"bytes,1,opt,name=invocation_id,json=invocationId,proto3" json:"invocation_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Invocation_Id) Reset()         { *m = Invocation_Id{} }
func (m *Invocation_Id) String() string { return proto.CompactTextString(m) }
func (*Invocation_Id) ProtoMessage()    {}
func (*Invocation_Id) Descriptor() ([]byte, []int) {
	return fileDescriptor_9a57b6cf1112b76d, []int{0, 0}
}

func (m *Invocation_Id) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Invocation_Id.Unmarshal(m, b)
}
func (m *Invocation_Id) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Invocation_Id.Marshal(b, m, deterministic)
}
func (m *Invocation_Id) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Invocation_Id.Merge(m, src)
}
func (m *Invocation_Id) XXX_Size() int {
	return xxx_messageInfo_Invocation_Id.Size(m)
}
func (m *Invocation_Id) XXX_DiscardUnknown() {
	xxx_messageInfo_Invocation_Id.DiscardUnknown(m)
}

var xxx_messageInfo_Invocation_Id proto.InternalMessageInfo

func (m *Invocation_Id) GetInvocationId() string {
	if m != nil {
		return m.InvocationId
	}
	return ""
}

// If known, represents the state of the user/build-system workspace.
type WorkspaceContext struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WorkspaceContext) Reset()         { *m = WorkspaceContext{} }
func (m *WorkspaceContext) String() string { return proto.CompactTextString(m) }
func (*WorkspaceContext) ProtoMessage()    {}
func (*WorkspaceContext) Descriptor() ([]byte, []int) {
	return fileDescriptor_9a57b6cf1112b76d, []int{1}
}

func (m *WorkspaceContext) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WorkspaceContext.Unmarshal(m, b)
}
func (m *WorkspaceContext) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WorkspaceContext.Marshal(b, m, deterministic)
}
func (m *WorkspaceContext) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WorkspaceContext.Merge(m, src)
}
func (m *WorkspaceContext) XXX_Size() int {
	return xxx_messageInfo_WorkspaceContext.Size(m)
}
func (m *WorkspaceContext) XXX_DiscardUnknown() {
	xxx_messageInfo_WorkspaceContext.DiscardUnknown(m)
}

var xxx_messageInfo_WorkspaceContext proto.InternalMessageInfo

// Describes the workspace under which the tool was invoked, this includes
// information that was fed into the command, the source code referenced, and
// the tool itself.
type WorkspaceInfo struct {
	// Data about the workspace that might be useful for debugging.
	WorkspaceContext *WorkspaceContext `protobuf:"bytes,1,opt,name=workspace_context,json=workspaceContext,proto3" json:"workspace_context,omitempty"`
	// Where the tool was invoked
	Hostname string `protobuf:"bytes,3,opt,name=hostname,proto3" json:"hostname,omitempty"`
	// The client's working directory where the build/test was run from.
	WorkingDirectory string `protobuf:"bytes,4,opt,name=working_directory,json=workingDirectory,proto3" json:"working_directory,omitempty"`
	// Tools should set tool_tag to the name of the tool or use case.
	ToolTag string `protobuf:"bytes,5,opt,name=tool_tag,json=toolTag,proto3" json:"tool_tag,omitempty"`
	// The command lines invoked. The first command line is the one typed by the
	// user, then each one after that should be an expansion of the previous
	// command line.
	CommandLines         []*CommandLine `protobuf:"bytes,7,rep,name=command_lines,json=commandLines,proto3" json:"command_lines,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *WorkspaceInfo) Reset()         { *m = WorkspaceInfo{} }
func (m *WorkspaceInfo) String() string { return proto.CompactTextString(m) }
func (*WorkspaceInfo) ProtoMessage()    {}
func (*WorkspaceInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_9a57b6cf1112b76d, []int{2}
}

func (m *WorkspaceInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WorkspaceInfo.Unmarshal(m, b)
}
func (m *WorkspaceInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WorkspaceInfo.Marshal(b, m, deterministic)
}
func (m *WorkspaceInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WorkspaceInfo.Merge(m, src)
}
func (m *WorkspaceInfo) XXX_Size() int {
	return xxx_messageInfo_WorkspaceInfo.Size(m)
}
func (m *WorkspaceInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_WorkspaceInfo.DiscardUnknown(m)
}

var xxx_messageInfo_WorkspaceInfo proto.InternalMessageInfo

func (m *WorkspaceInfo) GetWorkspaceContext() *WorkspaceContext {
	if m != nil {
		return m.WorkspaceContext
	}
	return nil
}

func (m *WorkspaceInfo) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *WorkspaceInfo) GetWorkingDirectory() string {
	if m != nil {
		return m.WorkingDirectory
	}
	return ""
}

func (m *WorkspaceInfo) GetToolTag() string {
	if m != nil {
		return m.ToolTag
	}
	return ""
}

func (m *WorkspaceInfo) GetCommandLines() []*CommandLine {
	if m != nil {
		return m.CommandLines
	}
	return nil
}

// The command and arguments that produced this Invocation.
type CommandLine struct {
	// A label describing this command line.
	Label string `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty"`
	// The command-line tool that is run: argv[0].
	Tool string `protobuf:"bytes,2,opt,name=tool,proto3" json:"tool,omitempty"`
	// The arguments to the above tool: argv[1]...argv[N].
	Args []string `protobuf:"bytes,3,rep,name=args,proto3" json:"args,omitempty"`
	// The actual command that was run with the tool.  (e.g. "build", or "test")
	// Omit if the tool doesn't accept a command.
	// This is a duplicate of one of the fields in args.
	Command              string   `protobuf:"bytes,4,opt,name=command,proto3" json:"command,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandLine) Reset()         { *m = CommandLine{} }
func (m *CommandLine) String() string { return proto.CompactTextString(m) }
func (*CommandLine) ProtoMessage()    {}
func (*CommandLine) Descriptor() ([]byte, []int) {
	return fileDescriptor_9a57b6cf1112b76d, []int{3}
}

func (m *CommandLine) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandLine.Unmarshal(m, b)
}
func (m *CommandLine) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandLine.Marshal(b, m, deterministic)
}
func (m *CommandLine) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandLine.Merge(m, src)
}
func (m *CommandLine) XXX_Size() int {
	return xxx_messageInfo_CommandLine.Size(m)
}
func (m *CommandLine) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandLine.DiscardUnknown(m)
}

var xxx_messageInfo_CommandLine proto.InternalMessageInfo

func (m *CommandLine) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *CommandLine) GetTool() string {
	if m != nil {
		return m.Tool
	}
	return ""
}

func (m *CommandLine) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *CommandLine) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

// Attributes that apply to all invocations.
type InvocationAttributes struct {
	// Immutable.
	// The Cloud Project that owns this invocation (this is different than the
	// Consumer Cloud Project that calls this API).
	// This must be set in the CreateInvocation call, and can't be changed.
	ProjectId string `protobuf:"bytes,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	// The list of users in the command chain.  The first user in this sequence
	// is the one who instigated the first command in the chain.
	Users []string `protobuf:"bytes,2,rep,name=users,proto3" json:"users,omitempty"`
	// Labels to categorize this invocation.
	// This is implemented as a set. All labels will be unique. Any duplicate
	// labels added will be ignored. Labels will be returned in lexicographical
	// order. Labels should be a list of words describing the Invocation. Labels
	// should be short, easy to read, and you shouldn't have more than a handful.
	// Labels should not be used for unique properties such as unique IDs. Use
	// properties in cases that don't meet these conditions.
	Labels []string `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty"`
	// This field describes the overall context or purpose of this invocation.
	// It will be used in the UI to give users more information about
	// how or why this invocation was run.
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	// If this Invocation was run in the context of a larger Continuous
	// Integration build or other automated system, this field may contain more
	// information about the greater context.
	InvocationContexts   []*InvocationContext `protobuf:"bytes,6,rep,name=invocation_contexts,json=invocationContexts,proto3" json:"invocation_contexts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *InvocationAttributes) Reset()         { *m = InvocationAttributes{} }
func (m *InvocationAttributes) String() string { return proto.CompactTextString(m) }
func (*InvocationAttributes) ProtoMessage()    {}
func (*InvocationAttributes) Descriptor() ([]byte, []int) {
	return fileDescriptor_9a57b6cf1112b76d, []int{4}
}

func (m *InvocationAttributes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InvocationAttributes.Unmarshal(m, b)
}
func (m *InvocationAttributes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InvocationAttributes.Marshal(b, m, deterministic)
}
func (m *InvocationAttributes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InvocationAttributes.Merge(m, src)
}
func (m *InvocationAttributes) XXX_Size() int {
	return xxx_messageInfo_InvocationAttributes.Size(m)
}
func (m *InvocationAttributes) XXX_DiscardUnknown() {
	xxx_messageInfo_InvocationAttributes.DiscardUnknown(m)
}

var xxx_messageInfo_InvocationAttributes proto.InternalMessageInfo

func (m *InvocationAttributes) GetProjectId() string {
	if m != nil {
		return m.ProjectId
	}
	return ""
}

func (m *InvocationAttributes) GetUsers() []string {
	if m != nil {
		return m.Users
	}
	return nil
}

func (m *InvocationAttributes) GetLabels() []string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *InvocationAttributes) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *InvocationAttributes) GetInvocationContexts() []*InvocationContext {
	if m != nil {
		return m.InvocationContexts
	}
	return nil
}

// Describes the invocation context which includes a display name and URL.
type InvocationContext struct {
	// A human readable name for the context under which this Invocation was run.
	DisplayName string `protobuf:"bytes,1,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	// A URL pointing to a UI containing more information
	Url                  string   `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InvocationContext) Reset()         { *m = InvocationContext{} }
func (m *InvocationContext) String() string { return proto.CompactTextString(m) }
func (*InvocationContext) ProtoMessage()    {}
func (*InvocationContext) Descriptor() ([]byte, []int) {
	return fileDescriptor_9a57b6cf1112b76d, []int{5}
}

func (m *InvocationContext) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InvocationContext.Unmarshal(m, b)
}
func (m *InvocationContext) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InvocationContext.Marshal(b, m, deterministic)
}
func (m *InvocationContext) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InvocationContext.Merge(m, src)
}
func (m *InvocationContext) XXX_Size() int {
	return xxx_messageInfo_InvocationContext.Size(m)
}
func (m *InvocationContext) XXX_DiscardUnknown() {
	xxx_messageInfo_InvocationContext.DiscardUnknown(m)
}

var xxx_messageInfo_InvocationContext proto.InternalMessageInfo

func (m *InvocationContext) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *InvocationContext) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func init() {
	proto.RegisterType((*Invocation)(nil), "google.devtools.resultstore.v2.Invocation")
	proto.RegisterType((*Invocation_Id)(nil), "google.devtools.resultstore.v2.Invocation.Id")
	proto.RegisterType((*WorkspaceContext)(nil), "google.devtools.resultstore.v2.WorkspaceContext")
	proto.RegisterType((*WorkspaceInfo)(nil), "google.devtools.resultstore.v2.WorkspaceInfo")
	proto.RegisterType((*CommandLine)(nil), "google.devtools.resultstore.v2.CommandLine")
	proto.RegisterType((*InvocationAttributes)(nil), "google.devtools.resultstore.v2.InvocationAttributes")
	proto.RegisterType((*InvocationContext)(nil), "google.devtools.resultstore.v2.InvocationContext")
}

func init() {
	proto.RegisterFile("google/devtools/resultstore/v2/invocation.proto", fileDescriptor_9a57b6cf1112b76d)
}

var fileDescriptor_9a57b6cf1112b76d = []byte{
	// 773 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x56, 0x71, 0x6b, 0x1a, 0x49,
	0x14, 0x47, 0x4d, 0x34, 0x3e, 0xe3, 0xa1, 0x73, 0x5e, 0xd8, 0x13, 0xee, 0xf0, 0xbc, 0xe3, 0x30,
	0x84, 0xac, 0x8d, 0x4d, 0x29, 0xa4, 0xb4, 0x90, 0xa6, 0x2d, 0x11, 0x42, 0x91, 0x4d, 0xa0, 0x50,
	0x28, 0xdb, 0x71, 0x77, 0x9c, 0x4e, 0xba, 0xee, 0x6c, 0x67, 0x46, 0x53, 0xbf, 0x4f, 0xe9, 0x77,
	0xeb, 0xb7, 0x28, 0x33, 0xbb, 0xab, 0x1b, 0x6b, 0xbb, 0xfe, 0x37, 0xef, 0xb7, 0xef, 0xf7, 0x7b,
	0x6f, 0x9e, 0xef, 0xbd, 0x11, 0xfa, 0x94, 0x73, 0x1a, 0x90, 0xbe, 0x4f, 0xe6, 0x8a, 0xf3, 0x40,
	0xf6, 0x05, 0x91, 0xb3, 0x40, 0x49, 0xc5, 0x05, 0xe9, 0xcf, 0x07, 0x7d, 0x16, 0xce, 0xb9, 0x87,
	0x15, 0xe3, 0xa1, 0x1d, 0x09, 0xae, 0x38, 0xfa, 0x3b, 0x26, 0xd8, 0x29, 0xc1, 0xce, 0x10, 0xec,
	0xf9, 0xa0, 0x7d, 0x94, 0x23, 0xe8, 0xf1, 0xe9, 0x34, 0x15, 0x6b, 0x1f, 0xe7, 0x3a, 0xcf, 0x89,
	0xc0, 0x94, 0x24, 0xee, 0x8f, 0xb6, 0x74, 0x77, 0xe5, 0x6c, 0x3a, 0xc5, 0x62, 0x91, 0xd0, 0x0e,
	0x73, 0x68, 0x13, 0x16, 0xa4, 0x11, 0xce, 0xb6, 0x70, 0x75, 0x23, 0xc1, 0x3d, 0x22, 0x25, 0x0b,
	0xa9, 0x4b, 0x84, 0xe0, 0x22, 0xe6, 0x76, 0xbf, 0x56, 0x00, 0x86, 0xcb, 0x72, 0x21, 0x04, 0x3b,
	0x21, 0x9e, 0x12, 0xab, 0xd0, 0x29, 0xf4, 0xaa, 0x8e, 0x39, 0xa3, 0xa7, 0x50, 0x64, 0xbe, 0x55,
	0xec, 0x14, 0x7a, 0xb5, 0xc1, 0xb1, 0xfd, 0xeb, 0x4a, 0xda, 0x2b, 0x2d, 0x7b, 0xe8, 0x3b, 0x45,
	0xe6, 0xa3, 0x77, 0xd0, 0x94, 0x0a, 0xab, 0x99, 0x74, 0xb1, 0x52, 0x82, 0x8d, 0x67, 0x8a, 0x48,
	0xab, 0x64, 0xd4, 0x1e, 0xe4, 0xa9, 0x5d, 0x1b, 0xe2, 0xf9, 0x92, 0xe7, 0x34, 0xe4, 0x1a, 0x82,
	0x9e, 0x41, 0x59, 0xb1, 0x29, 0x0b, 0xa9, 0xb5, 0x63, 0x34, 0xff, 0xcf, 0xd3, 0xbc, 0x31, 0xde,
	0x4e, 0xc2, 0x42, 0x0c, 0xfe, 0x58, 0xb5, 0x4b, 0x36, 0xc5, 0x5d, 0x23, 0x77, 0xba, 0xfd, 0x85,
	0x33, 0x69, 0xb6, 0xd8, 0x06, 0x14, 0xdd, 0xc0, 0x6f, 0x77, 0x5c, 0x7c, 0x94, 0x11, 0xf6, 0x88,
	0xcb, 0xc2, 0x09, 0xb7, 0xca, 0xdb, 0x15, 0xf5, 0x4d, 0xca, 0x1a, 0x86, 0x13, 0xee, 0xd4, 0xef,
	0xb2, 0x26, 0xba, 0x04, 0x88, 0x04, 0x8f, 0x88, 0x50, 0x8c, 0x48, 0xab, 0xd2, 0x29, 0xf5, 0x6a,
	0x83, 0x5e, 0x9e, 0xe2, 0x28, 0x66, 0x2c, 0x9c, 0x0c, 0x17, 0x9d, 0xc1, 0xae, 0x6e, 0x15, 0x69,
	0xed, 0x19, 0x91, 0xff, 0xf2, 0x44, 0x5e, 0xb1, 0x80, 0x38, 0x31, 0x05, 0x4d, 0x00, 0xad, 0x35,
	0xb2, 0xce, 0xa6, 0x6a, 0x84, 0x1e, 0xe7, 0x09, 0x5d, 0xe1, 0x90, 0xce, 0x30, 0x25, 0x17, 0x89,
	0xc2, 0x75, 0x3c, 0x09, 0x4e, 0xd3, 0xbb, 0x07, 0xe8, 0x1c, 0xdf, 0x03, 0xc2, 0x94, 0x0a, 0x42,
	0xb1, 0x22, 0x6e, 0xfa, 0xd9, 0x02, 0x53, 0xc7, 0x93, 0xbc, 0x38, 0xe7, 0x29, 0x33, 0x0d, 0xe4,
	0x34, 0xf1, 0x3a, 0x84, 0x6e, 0xe1, 0x60, 0xe3, 0xc0, 0x48, 0xab, 0x66, 0x6e, 0x73, 0xba, 0x4d,
	0x59, 0x46, 0x4b, 0xf2, 0x4b, 0xc3, 0x75, 0x5a, 0x93, 0x0d, 0x68, 0xfb, 0x10, 0x8a, 0x43, 0x1f,
	0xfd, 0x0b, 0xf5, 0x4c, 0x0b, 0x32, 0x3f, 0x99, 0xbe, 0xfd, 0x15, 0x38, 0xf4, 0xbb, 0x08, 0x1a,
	0xcb, 0x36, 0xb8, 0xe0, 0xa1, 0x22, 0x9f, 0x55, 0xf7, 0x4b, 0x11, 0xea, 0xf7, 0x7a, 0x43, 0x0f,
	0xdb, 0xaa, 0xc5, 0xbc, 0xd8, 0xcd, 0xc8, 0x6d, 0x31, 0x6c, 0xeb, 0xf2, 0x4e, 0xe3, 0x6e, 0x0d,
	0x41, 0x6d, 0xd8, 0xfb, 0xc0, 0xa5, 0x32, 0x2b, 0xa2, 0x64, 0x92, 0x5c, 0xda, 0xe8, 0x28, 0x0e,
	0xad, 0xeb, 0xe5, 0x33, 0x41, 0x3c, 0xc5, 0xc5, 0xc2, 0xcc, 0x64, 0x35, 0x16, 0x62, 0x21, 0x7d,
	0x91, 0xe2, 0xe8, 0x4f, 0xd8, 0xd3, 0x39, 0xb8, 0x0a, 0x53, 0x33, 0x68, 0x55, 0xa7, 0xa2, 0xed,
	0x1b, 0x4c, 0xd1, 0x08, 0xea, 0x7a, 0xdd, 0xe2, 0xd0, 0x77, 0x03, 0x16, 0x2e, 0x5b, 0xfa, 0x28,
	0x2f, 0xfd, 0x8b, 0x98, 0x74, 0xc5, 0x42, 0xe2, 0xec, 0x7b, 0x2b, 0x43, 0x76, 0x09, 0xd4, 0x32,
	0x1f, 0x51, 0x0b, 0x76, 0x03, 0x3c, 0x26, 0x41, 0x52, 0xe6, 0xd8, 0xd0, 0x9b, 0x4f, 0xcb, 0x9a,
	0x3d, 0x57, 0x75, 0xcc, 0x59, 0x63, 0x58, 0x50, 0xbd, 0xad, 0x4a, 0x1a, 0xd3, 0x67, 0x64, 0x41,
	0x25, 0x11, 0x4f, 0x2e, 0x97, 0x9a, 0xdd, 0x6f, 0x05, 0x68, 0x6d, 0xda, 0x06, 0xe8, 0x2f, 0x33,
	0xa1, 0xb7, 0xc4, 0x53, 0xab, 0x1f, 0xb7, 0x9a, 0x20, 0x43, 0x5f, 0xe7, 0x33, 0x93, 0x44, 0x48,
	0xab, 0x68, 0xc2, 0xc4, 0x06, 0x3a, 0x80, 0xb2, 0x49, 0x2c, 0x8d, 0x9e, 0x58, 0xa8, 0x03, 0x35,
	0x9f, 0x48, 0x4f, 0xb0, 0x48, 0x47, 0x49, 0x72, 0xc8, 0x42, 0x68, 0x0c, 0xbf, 0x67, 0xda, 0x29,
	0x69, 0x02, 0x69, 0x95, 0x4d, 0x19, 0x4f, 0xb6, 0xdf, 0x67, 0x69, 0x1b, 0x20, 0xb6, 0x0e, 0xc9,
	0xee, 0x25, 0x34, 0x7f, 0x70, 0x44, 0xff, 0xc0, 0xbe, 0xcf, 0x64, 0x14, 0xe0, 0x85, 0x9b, 0x79,
	0x44, 0x6a, 0x09, 0xf6, 0x5a, 0x37, 0x49, 0x03, 0x4a, 0x33, 0x91, 0x16, 0x59, 0x1f, 0x9f, 0x7f,
	0x82, 0xae, 0xc7, 0xa7, 0x39, 0x59, 0x8d, 0x0a, 0x6f, 0x87, 0x89, 0x07, 0xe5, 0x01, 0x0e, 0xa9,
	0xcd, 0x05, 0xed, 0x53, 0x12, 0x9a, 0x47, 0x2c, 0xf9, 0x3b, 0x80, 0x23, 0x26, 0x7f, 0xf6, 0x06,
	0x3e, 0xc9, 0x98, 0xe3, 0xb2, 0x61, 0x3d, 0xfc, 0x1e, 0x00, 0x00, 0xff, 0xff, 0x10, 0xbe, 0x01,
	0x00, 0x47, 0x08, 0x00, 0x00,
}