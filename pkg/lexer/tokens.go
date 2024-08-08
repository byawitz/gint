// Package lexer is for parsing PHP syntax
// List is taken from here https://github.com/php/php-src/blob/master/ext/tokenizer/tokenizer_data.c
package lexer

import (
	"fmt"
	"slices"
)

type Kind int

const (
	TEOF Kind = iota

	TLNumber
	TDNumber
	TString
	TNameFullyQualified
	TNameRelative
	TNameQualified
	TVariable
	TInlineHtml
	TEncapsedAndWhitespace
	TConstantEncapsedString
	TStringVarName
	TNumString
	TEval
	TNew
	TExit

	TThrow
	TInclude
	TIncludeOnce
	TRequire
	TRequireOnce
	TLogicalOr
	TLogicalXor
	TLogicalAnd
	TPrint
	TYield
	TDoubleArrow
	TYieldFrom
	TAssignment
	TPlus
	TDash
	TSlash
	TStar
	TPercent
	TPlusEqual
	TMinusEqual
	TMulEqual
	TDivEqual
	TConcatEqual
	TConcat
	TComma
	TModEqual
	TAndEqual
	TOrEqual
	TXorEqual
	TSlEqual
	TSrEqual
	TPowEqual
	TCoalesceEqual
	TCoalesce
	TBooleanOr
	TAt
	TPipe
	TBooleanAnd
	TAmpersandNotFollowedByVarOrVararg
	TAmpersandFollowedByVarOrVararg
	TNotEqual
	TIsEqual
	TIsNotEqual
	TIsIdentical
	TIsNotIdentical
	TSpaceship
	TIsSmallerOrEqual
	TIsGreaterOrEqual
	TIsSmaller
	TIsGreater
	TSl
	TSr
	TInstanceof
	TIntCast
	TDoubleCast
	TStringCast
	TArrayCast
	TObjectCast
	TBoolCast
	TUnsetCast
	TPow
	TClone
	TNoElse
	TElseif
	TElse
	TIf
	TEndif
	TEcho
	TDo
	TWhile
	TEndWhile
	TFor
	TEndFor
	TForeach
	TEndForeach
	TDeclare
	TEndDeclare
	TAs
	TSwitch
	TEndSwitch
	TCase
	TDefault
	TMatch
	TBreak
	TContinue
	TGoto
	TFunction
	TFn
	TConst
	TReturn
	TTry
	TCatch
	TFinally
	TUse
	TInsteadof
	TGlobal
	TStatic
	TAbstract
	TFinal
	TPrivate
	TProtected
	TPublic
	TReadonly
	TVar
	TUnset
	TIsset
	TEmpty
	THaltCompiler
	TClass
	TTrait
	TInterface
	TEnum
	TExtends
	TImplements
	TNamespace
	TList
	TArray
	TCallable
	TLine
	TFile
	TDir
	TClassC
	TTraitC
	TMethodC
	TFuncC
	TPropertyC
	TNsC
	TAttribute
	TInc
	TDec
	TObjectOperator
	TNullSafeObjectOperator
	TComment
	TDocComment
	TOpenTag
	TOpenTagWithEcho
	TCloseTag
	TWhitespace
	TStartHeredoc
	TEndHeredoc
	TDollarOpenCurlyBraces
	TCurlyOpen
	TPaamayimNekudotayim
	TNsSeparator
	TEllipsis
	TBadCharacter
	TError
	TSemiColon
	TColon
	TQuestion
	TOpenBracket
	TCloseBracket
	TOPENCurly
	TCloseCurly
	TOpenParen
	TCloseParen
)

type Token struct {
	Value string
	Kind  Kind
}

func (t Token) isOneOfMany(tokens ...Kind) bool {
	if slices.Contains(tokens, t.Kind) {
		return true
	}
	return false
}
func (t Token) Debug() {
	if t.isOneOfMany(TLNumber, TDNumber, TString) {
		fmt.Printf("%s (%s)\n", TokenKindString(t.Kind), t.Value)
	} else {
		fmt.Printf("%s ()\n", TokenKindString(t.Kind))
	}
}

func NewToken(kind Kind, value string) Token {
	return Token{value, kind}
}

func TokenKindString(token Kind) string {
	switch token {
	case TLNumber:
		return "T_LNUMBER"
	case TDNumber:
		return "T_DNUMBER"
	case TString:
		return "T_STRING"
	case TNameFullyQualified:
		return "T_NAME_FULLY_QUALIFIED"
	case TNameRelative:
		return "T_NAME_RELATIVE"
	case TNameQualified:
		return "T_NAME_QUALIFIED"
	case TVariable:
		return "T_VARIABLE"
	case TInlineHtml:
		return "T_INLINE_HTML"
	case TEncapsedAndWhitespace:
		return "T_ENCAPSED_AND_WHITESPACE"
	case TConstantEncapsedString:
		return "T_CONSTANT_ENCAPSED_STRING"
	case TStringVarName:
		return "T_STRING_VARNAME"
	case TNumString:
		return "T_NUM_STRING"
	case TInclude:
		return "T_INCLUDE"
	case TIncludeOnce:
		return "T_INCLUDE_ONCE"
	case TEval:
		return "T_EVAL"
	case TRequire:
		return "T_REQUIRE"
	case TRequireOnce:
		return "T_REQUIRE_ONCE"
	case TLogicalOr:
		return "T_LOGICAL_OR"
	case TLogicalXor:
		return "T_LOGICAL_XOR"
	case TLogicalAnd:
		return "T_LOGICAL_AND"
	case TPrint:
		return "T_PRINT"
	case TYield:
		return "T_YIELD"
	case TYieldFrom:
		return "T_YIELD_FROM"
	case TInstanceof:
		return "T_INSTANCEOF"
	case TNew:
		return "T_NEW"
	case TClone:
		return "T_CLONE"
	case TExit:
		return "T_EXIT"
	case TIf:
		return "T_IF"
	case TElseif:
		return "T_ELSEIF"
	case TElse:
		return "T_ELSE"
	case TEndif:
		return "T_ENDIF"
	case TEcho:
		return "T_ECHO"
	case TDo:
		return "T_DO"
	case TWhile:
		return "T_WHILE"
	case TEndWhile:
		return "T_ENDWHILE"
	case TFor:
		return "T_FOR"
	case TEndFor:
		return "T_ENDFOR"
	case TForeach:
		return "T_FOREACH"
	case TEndForeach:
		return "T_ENDFOREACH"
	case TDeclare:
		return "T_DECLARE"
	case TEndDeclare:
		return "T_ENDDECLARE"
	case TAs:
		return "T_AS"
	case TSwitch:
		return "T_SWITCH"
	case TEndSwitch:
		return "T_ENDSWITCH"
	case TCase:
		return "T_CASE"
	case TDefault:
		return "T_DEFAULT"
	case TMatch:
		return "T_MATCH"
	case TBreak:
		return "T_BREAK"
	case TContinue:
		return "T_CONTINUE"
	case TGoto:
		return "T_GOTO"
	case TFunction:
		return "T_FUNCTION"
	case TFn:
		return "T_FN"
	case TConst:
		return "T_CONST"
	case TReturn:
		return "T_RETURN"
	case TTry:
		return "T_TRY"
	case TCatch:
		return "T_CATCH"
	case TFinally:
		return "T_FINALLY"
	case TThrow:
		return "T_THROW"
	case TUse:
		return "T_USE"
	case TInsteadof:
		return "T_INSTEADOF"
	case TGlobal:
		return "T_GLOBAL"
	case TStatic:
		return "T_STATIC"
	case TAbstract:
		return "T_ABSTRACT"
	case TFinal:
		return "T_FINAL"
	case TPrivate:
		return "T_PRIVATE"
	case TProtected:
		return "T_PROTECTED"
	case TPublic:
		return "T_PUBLIC"
	case TReadonly:
		return "T_READONLY"
	case TVar:
		return "T_VAR"
	case TUnset:
		return "T_UNSET"
	case TIsset:
		return "T_ISSET"
	case TEmpty:
		return "T_EMPTY"
	case THaltCompiler:
		return "T_HALT_COMPILER"
	case TClass:
		return "T_CLASS"
	case TTrait:
		return "T_TRAIT"
	case TInterface:
		return "T_INTERFACE"
	case TEnum:
		return "T_ENUM"
	case TExtends:
		return "T_EXTENDS"
	case TImplements:
		return "T_IMPLEMENTS"
	case TNamespace:
		return "T_NAMESPACE"
	case TList:
		return "T_LIST"
	case TArray:
		return "T_ARRAY"
	case TCallable:
		return "T_CALLABLE"
	case TLine:
		return "T_LINE"
	case TFile:
		return "T_FILE"
	case TDir:
		return "T_DIR"
	case TClassC:
		return "T_CLASS_C"
	case TTraitC:
		return "T_TRAIT_C"
	case TMethodC:
		return "T_METHOD_C"
	case TFuncC:
		return "T_FUNC_C"
	case TPropertyC:
		return "T_PROPERTY_C"
	case TNsC:
		return "T_NS_C"
	case TAttribute:
		return "T_ATTRIBUTE"
	case TAssignment:
		return "T_ASSIGNMENT"
	case TPlus:
		return "T_PLUS"
	case TDash:
		return "T_DASH"
	case TSlash:
		return "T_SLASH"
	case TStar:
		return "T_STAR"
	case TPercent:
		return "T_PERCENT"
	case TPlusEqual:
		return "T_PLUS_EQUAL"
	case TMinusEqual:
		return "T_MINUS_EQUAL"
	case TMulEqual:
		return "T_MUL_EQUAL"
	case TDivEqual:
		return "T_DIV_EQUAL"
	case TConcatEqual:
		return "T_CONCAT_EQUAL"
	case TConcat:
		return "T_CONCAT"
	case TComma:
		return "T_COMMA"
	case TModEqual:
		return "T_MOD_EQUAL"
	case TAndEqual:
		return "T_AND_EQUAL"
	case TOrEqual:
		return "T_OR_EQUAL"
	case TXorEqual:
		return "T_XOR_EQUAL"
	case TSlEqual:
		return "T_SL_EQUAL"
	case TSrEqual:
		return "T_SR_EQUAL"
	case TCoalesceEqual:
		return "T_COALESCE_EQUAL"
	case TBooleanOr:
		return "T_BOOLEAN_OR"
	case TAt:
		return "T_AT"
	case TPipe:
		return "T_PIPE"
	case TBooleanAnd:
		return "T_BOOLEAN_AND"
	case TIsEqual:
		return "T_IS_EQUAL"
	case TIsNotEqual:
		return "T_IS_NOT_EQUAL"
	case TNotEqual:
		return "T_NOT_EQUAL"
	case TIsIdentical:
		return "T_IS_IDENTICAL"
	case TIsNotIdentical:
		return "T_IS_NOT_IDENTICAL"
	case TIsSmallerOrEqual:
		return "T_IS_SMALLER_OR_EQUAL"
	case TIsGreaterOrEqual:
		return "T_IS_GREATER_OR_EQUAL"
	case TIsSmaller:
		return "T_IS_SMALLER"
	case TIsGreater:
		return "T_IS_GREATER"
	case TSpaceship:
		return "T_SPACESHIP"
	case TSl:
		return "T_SL"
	case TSr:
		return "T_SR"
	case TInc:
		return "T_INC"
	case TDec:
		return "T_DEC"
	case TIntCast:
		return "T_INT_CAST"
	case TDoubleCast:
		return "T_DOUBLE_CAST"
	case TStringCast:
		return "T_STRING_CAST"
	case TArrayCast:
		return "T_ARRAY_CAST"
	case TObjectCast:
		return "T_OBJECT_CAST"
	case TBoolCast:
		return "T_BOOL_CAST"
	case TUnsetCast:
		return "T_UNSET_CAST"
	case TObjectOperator:
		return "T_OBJECT_OPERATOR"
	case TNullSafeObjectOperator:
		return "T_NULLSAFE_OBJECT_OPERATOR"
	case TDoubleArrow:
		return "T_DOUBLE_ARROW"
	case TComment:
		return "T_COMMENT"
	case TDocComment:
		return "T_DOC_COMMENT"
	case TOpenTag:
		return "T_OPEN_TAG"
	case TOpenTagWithEcho:
		return "T_OPEN_TAG_WITH_ECHO"
	case TCloseTag:
		return "T_CLOSE_TAG"
	case TWhitespace:
		return "T_WHITESPACE"
	case TStartHeredoc:
		return "T_START_HEREDOC"
	case TEndHeredoc:
		return "T_END_HEREDOC"
	case TDollarOpenCurlyBraces:
		return "T_DOLLAR_OPEN_CURLY_BRACES"
	case TCurlyOpen:
		return "T_CURLY_OPEN"
	case TPaamayimNekudotayim:
		return "T_DOUBLE_COLON"
	case TNsSeparator:
		return "T_NS_SEPARATOR"
	case TEllipsis:
		return "T_ELLIPSIS"
	case TCoalesce:
		return "T_COALESCE"
	case TPow:
		return "T_POW"
	case TPowEqual:
		return "T_POW_EQUAL"
	case TAmpersandFollowedByVarOrVararg:
		return "T_AMPERSAND_FOLLOWED_BY_VAR_OR_VARARG"
	case TAmpersandNotFollowedByVarOrVararg:
		return "T_AMPERSAND_NOT_FOLLOWED_BY_VAR_OR_VARARG"
	case TBadCharacter:
		return "T_BAD_CHARACTER"
	case TSemiColon:
		return "T_SEMI_COLON"
	case TColon:
		return "T_COLON"
	case TQuestion:
		return "T_question"
	case TOpenBracket:
		return "T_OPEN_BRACKET"
	case TCloseBracket:
		return "T_CLOSE_BRACKET"
	case TOPENCurly:
		return "T_OPEN_CURLY"
	case TCloseCurly:
		return "T_CLOSE_CURLY"
	case TOpenParen:
		return "T_OPEN_PAREN"
	case TCloseParen:
		return "T_CLOSE_PAREN"
	default:
		return "T_BAD_CHARACTER"
	}
}
