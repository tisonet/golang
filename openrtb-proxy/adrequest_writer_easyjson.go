// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package main

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

func easyjson3fb78918DecodeGithubComTisonetGolangOpenrtbProxy(in *jlexer.Lexer, out *Position) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "impression_id":
			out.ImpressionId = string(in.String())
		case "position_id":
			out.PositionId = string(in.String())
		case "width":
			out.Width = int(in.Int())
		case "height":
			out.Height = int(in.Int())
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
func easyjson3fb78918EncodeGithubComTisonetGolangOpenrtbProxy(out *jwriter.Writer, in Position) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"impression_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ImpressionId))
	}
	{
		const prefix string = ",\"position_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.PositionId))
	}
	{
		const prefix string = ",\"width\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Width))
	}
	{
		const prefix string = ",\"height\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Height))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Position) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3fb78918EncodeGithubComTisonetGolangOpenrtbProxy(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Position) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3fb78918EncodeGithubComTisonetGolangOpenrtbProxy(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Position) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3fb78918DecodeGithubComTisonetGolangOpenrtbProxy(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Position) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3fb78918DecodeGithubComTisonetGolangOpenrtbProxy(l, v)
}
func easyjson3fb78918DecodeGithubComTisonetGolangOpenrtbProxy1(in *jlexer.Lexer, out *AdRequestResponseMessage) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "pageview_id":
			out.PageviewId = string(in.String())
		case "ibbid":
			out.Ibbid = string(in.String())
		case "ip":
			out.Ip = string(in.String())
		case "url":
			out.Url = string(in.String())
		case "user_agent":
			out.UserAgent = string(in.String())
		case "positions":
			if in.IsNull() {
				in.Skip()
				out.Positions = nil
			} else {
				in.Delim('[')
				if out.Positions == nil {
					if !in.IsDelim(']') {
						out.Positions = make([]Position, 0, 1)
					} else {
						out.Positions = []Position{}
					}
				} else {
					out.Positions = (out.Positions)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Position
					(v1).UnmarshalEasyJSON(in)
					out.Positions = append(out.Positions, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "prebid":
			out.Prebid = bool(in.Bool())
		case "timestamp":
			out.Timestamp = int64(in.Int64())
		case "status":
			out.Status = int(in.Int())
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
func easyjson3fb78918EncodeGithubComTisonetGolangOpenrtbProxy1(out *jwriter.Writer, in AdRequestResponseMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"pageview_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.PageviewId))
	}
	{
		const prefix string = ",\"ibbid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Ibbid))
	}
	{
		const prefix string = ",\"ip\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Ip))
	}
	{
		const prefix string = ",\"url\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Url))
	}
	{
		const prefix string = ",\"user_agent\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.UserAgent))
	}
	{
		const prefix string = ",\"positions\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Positions == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Positions {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"prebid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Prebid))
	}
	{
		const prefix string = ",\"timestamp\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Timestamp))
	}
	{
		const prefix string = ",\"status\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Status))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AdRequestResponseMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3fb78918EncodeGithubComTisonetGolangOpenrtbProxy1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AdRequestResponseMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3fb78918EncodeGithubComTisonetGolangOpenrtbProxy1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AdRequestResponseMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3fb78918DecodeGithubComTisonetGolangOpenrtbProxy1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AdRequestResponseMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3fb78918DecodeGithubComTisonetGolangOpenrtbProxy1(l, v)
}
