// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package forms

import (
	sql "database/sql"
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

func easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms(in *jlexer.Lexer, out *Vote) {
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
		case "nickname":
			out.Nickname = string(in.String())
		case "voice":
			out.Voice = int(in.Int())
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
func easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms(out *jwriter.Writer, in Vote) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"nickname\":"
		out.RawString(prefix[1:])
		out.String(string(in.Nickname))
	}
	{
		const prefix string = ",\"voice\":"
		out.RawString(prefix)
		out.Int(int(in.Voice))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Vote) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Vote) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Vote) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Vote) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms(l, v)
}
func easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms1(in *jlexer.Lexer, out *ThreadUpdate) {
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
		case "title":
			out.Title = string(in.String())
		case "message":
			out.Message = string(in.String())
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
func easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms1(out *jwriter.Writer, in ThreadUpdate) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Title != "" {
		const prefix string = ",\"title\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Title))
	}
	if in.Message != "" {
		const prefix string = ",\"message\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ThreadUpdate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ThreadUpdate) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ThreadUpdate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ThreadUpdate) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms1(l, v)
}
func easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms2(in *jlexer.Lexer, out *ThreadModel) {
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
		case "Id":
			out.Id = int(in.Int())
		case "Title":
			out.Title = string(in.String())
		case "Author":
			out.Author = string(in.String())
		case "Forum":
			out.Forum = string(in.String())
		case "Message":
			out.Message = string(in.String())
		case "Votes":
			out.Votes = int(in.Int())
		case "Slug":
			easyjson2d00218DecodeDatabaseSql(in, &out.Slug)
		case "Created":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Created).UnmarshalJSON(data))
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
func easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms2(out *jwriter.Writer, in ThreadModel) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Id))
	}
	{
		const prefix string = ",\"Title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"Author\":"
		out.RawString(prefix)
		out.String(string(in.Author))
	}
	{
		const prefix string = ",\"Forum\":"
		out.RawString(prefix)
		out.String(string(in.Forum))
	}
	{
		const prefix string = ",\"Message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"Votes\":"
		out.RawString(prefix)
		out.Int(int(in.Votes))
	}
	{
		const prefix string = ",\"Slug\":"
		out.RawString(prefix)
		easyjson2d00218EncodeDatabaseSql(out, in.Slug)
	}
	{
		const prefix string = ",\"Created\":"
		out.RawString(prefix)
		out.Raw((in.Created).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ThreadModel) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ThreadModel) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ThreadModel) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ThreadModel) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms2(l, v)
}
func easyjson2d00218DecodeDatabaseSql(in *jlexer.Lexer, out *sql.NullString) {
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
		case "String":
			out.String = string(in.String())
		case "Valid":
			out.Valid = bool(in.Bool())
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
func easyjson2d00218EncodeDatabaseSql(out *jwriter.Writer, in sql.NullString) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"String\":"
		out.RawString(prefix[1:])
		out.String(string(in.String))
	}
	{
		const prefix string = ",\"Valid\":"
		out.RawString(prefix)
		out.Bool(bool(in.Valid))
	}
	out.RawByte('}')
}
func easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms3(in *jlexer.Lexer, out *ThreadForm) {
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
		case "id":
			out.Id = int(in.Int())
		case "title":
			out.Title = string(in.String())
		case "author":
			out.Author = string(in.String())
		case "forum":
			out.Forum = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "votes":
			out.Votes = int(in.Int())
		case "slug":
			out.Slug = string(in.String())
		case "created":
			out.Created = string(in.String())
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
func easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms3(out *jwriter.Writer, in ThreadForm) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Id != 0 {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Int(int(in.Id))
	}
	{
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix)
		out.String(string(in.Author))
	}
	if in.Forum != "" {
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.String(string(in.Forum))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	if in.Votes != 0 {
		const prefix string = ",\"votes\":"
		out.RawString(prefix)
		out.Int(int(in.Votes))
	}
	{
		const prefix string = ",\"slug\":"
		out.RawString(prefix)
		out.String(string(in.Slug))
	}
	if in.Created != "" {
		const prefix string = ",\"created\":"
		out.RawString(prefix)
		out.String(string(in.Created))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ThreadForm) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ThreadForm) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ThreadForm) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ThreadForm) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms3(l, v)
}
func easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms4(in *jlexer.Lexer, out *ThreadFilter) {
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
		case "limit":
			out.Limit = int(in.Int())
		case "since":
			out.Since = string(in.String())
		case "desc":
			out.Desc = bool(in.Bool())
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
func easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms4(out *jwriter.Writer, in ThreadFilter) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Limit != 0 {
		const prefix string = ",\"limit\":"
		first = false
		out.RawString(prefix[1:])
		out.Int(int(in.Limit))
	}
	if in.Since != "" {
		const prefix string = ",\"since\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Since))
	}
	if in.Desc {
		const prefix string = ",\"desc\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Desc))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ThreadFilter) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ThreadFilter) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2d00218EncodeGithubComNNKulickovTechnoparkDbmsForumForms4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ThreadFilter) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ThreadFilter) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2d00218DecodeGithubComNNKulickovTechnoparkDbmsForumForms4(l, v)
}
