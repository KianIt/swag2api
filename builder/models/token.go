package models

import (
	astm "github.com/KianIt/swag2api/builder/ast/models"
)

// Token is a source code token.
//
// Token is used to add a particular piece of
// source code to the AST.
type Token string

const (
	MuxToken           Token = "mux"
	HandlerToken       Token = "http.Handler"
	NewMuxToken        Token = "http.NewServeMux"
	HandleFuncToken    Token = "HandleFunc"
	InitToken          Token = "init"
	WToken             Token = "w"
	WriterToken        Token = "http.ResponseWriter"
	RToken             Token = "r"
	RequestToken       Token = "*http.Request"
	PathValueGetToken  Token = "r.PathValue"
	QueryValueGetToken Token = "r.URL.Query().Get"
	UnmarshalErrToken  Token = "UnmarshalErr"
	BodyToken          Token = "body"
	BodyGetToken       Token = "io.ReadAll(r.Body)"
	BodyReadErrToken   Token = "bodyReadErr"
	BodyTypeToken      Token = "bodyType"
	BodyValueToken     Token = "bodyValue"
	ResultTypeToken    Token = "resultType"
	ResultValueToken   Token = "resultValue"
	ErrToken           Token = "err"
	NilToken           Token = "nil"
	ReturnToken        Token = "return"
	InlineFieldToken   Token = "`json:\",inline\"`"
	CodeToken          Token = "Code"
	MsgToken           Token = "Msg"
	StatusOkToken      Token = "http.StatusOK"
	StatusOkTextToken  Token = "http.StatusText(http.StatusOK)"
)

// AstExpr returns the token as an AST expression.
func (t Token) AstExpr() astm.Expr {
	return astm.GetNameExpr(string(t))
}

// String returns the token as a string.
func (t Token) String() string {
	return string(t)
}
