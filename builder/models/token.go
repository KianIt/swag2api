package models

import (
	astm "github.com/KianIt/swag2api/builder/ast/models"
)

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

func (t Token) AstExpr() astm.Expr {
	return astm.GetNameExpr(string(t))
}

func (t Token) String() string {
	return string(t)
}
