// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package pbehavior

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson950e241aDecodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior(in *jlexer.Lexer, out *Types) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "T":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.T = make(map[string]*Type)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 *Type
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(Type)
						}
						easyjson950e241aDecodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior1(in, v1)
					}
					(out.T)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson950e241aEncodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior(out *jwriter.Writer, in Types) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"T\":"
		out.RawString(prefix[1:])
		if in.T == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v2First := true
			for v2Name, v2Value := range in.T {
				if v2First {
					v2First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v2Name))
				out.RawByte(':')
				if v2Value == nil {
					out.RawString("null")
				} else {
					easyjson950e241aEncodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior1(out, *v2Value)
				}
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Types) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson950e241aEncodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Types) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson950e241aDecodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior(l, v)
}
func easyjson950e241aDecodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior1(in *jlexer.Lexer, out *Type) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "_id":
			out.ID = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "type":
			out.Type = string(in.String())
		case "priority":
			out.Priority = int(in.Int())
		case "icon_name":
			out.IconName = string(in.String())
		case "color":
			out.Color = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson950e241aEncodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior1(out *jwriter.Writer, in Type) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"_id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"priority\":"
		out.RawString(prefix)
		out.Int(int(in.Priority))
	}
	{
		const prefix string = ",\"icon_name\":"
		out.RawString(prefix)
		out.String(string(in.IconName))
	}
	if in.Color != "" {
		const prefix string = ",\"color\":"
		out.RawString(prefix)
		out.String(string(in.Color))
	}
	out.RawByte('}')
}
func easyjson950e241aDecodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior2(in *jlexer.Lexer, out *ComputedPbehavior) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "n":
			out.Name = string(in.String())
		case "r":
			out.Reason = string(in.String())
		case "f":
			out.Filter = string(in.String())
		case "t":
			if in.IsNull() {
				in.Skip()
				out.Types = nil
			} else {
				in.Delim('[')
				if out.Types == nil {
					if !in.IsDelim(']') {
						out.Types = make([]computedType, 0, 0)
					} else {
						out.Types = []computedType{}
					}
				} else {
					out.Types = (out.Types)[:0]
				}
				for !in.IsDelim(']') {
					var v3 computedType
					easyjson950e241aDecodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior3(in, &v3)
					out.Types = append(out.Types, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "c":
			out.Created = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson950e241aEncodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior2(out *jwriter.Writer, in ComputedPbehavior) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"n\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"r\":"
		out.RawString(prefix)
		out.String(string(in.Reason))
	}
	{
		const prefix string = ",\"f\":"
		out.RawString(prefix)
		out.String(string(in.Filter))
	}
	{
		const prefix string = ",\"t\":"
		out.RawString(prefix)
		if in.Types == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v4, v5 := range in.Types {
				if v4 > 0 {
					out.RawByte(',')
				}
				easyjson950e241aEncodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior3(out, v5)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"c\":"
		out.RawString(prefix)
		out.Int64(int64(in.Created))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ComputedPbehavior) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson950e241aEncodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior2(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ComputedPbehavior) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson950e241aDecodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior2(l, v)
}
func easyjson950e241aDecodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior3(in *jlexer.Lexer, out *computedType) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "s":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Span).UnmarshalJSON(data))
			}
		case "t":
			out.ID = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson950e241aEncodeGitCanopsisNetCanopsisCanopsisCommunityCommunityGoEnginesCommunityLibCanopsisPbehavior3(out *jwriter.Writer, in computedType) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"s\":"
		out.RawString(prefix[1:])
		out.Raw((in.Span).MarshalJSON())
	}
	{
		const prefix string = ",\"t\":"
		out.RawString(prefix)
		out.String(string(in.ID))
	}
	out.RawByte('}')
}