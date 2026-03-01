package models

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken_AstExpr(t *testing.T) {
	token := MuxToken
	astExpr := token.AstExpr()

	assert.IsType(t, &ast.Ident{}, astExpr)
	assert.Equal(t, token.String(), astExpr.(*ast.Ident).Name)
}

func TestToken_String(t *testing.T) {
	token := MuxToken
	str := token.String()

	assert.IsType(t, "", str)
	assert.Equal(t, string(MuxToken), str)
}
