package lexer

import (
	"fmt"
	"regexp"
)

const (
	maxErrorShowingLength = 35
)

type regexHandler func(lex *lexer, regex *regexp.Regexp)

type regexPattern struct {
	regex        *regexp.Regexp
	handler      regexHandler
	outSideOfPHP bool
}
type lexer struct {
	patterns  []regexPattern
	Tokens    []Token
	source    string
	pos       int
	insidePHP bool
}

func (l *lexer) advance(n int) {
	l.pos += n
}

func (l *lexer) at() byte {
	return l.source[l.pos]
}

func (l *lexer) remainder() string {
	return l.source[l.pos:]
}

func (l *lexer) atEof() bool {
	return l.pos >= len(l.source)
}

func (l *lexer) push(token Token) {
	l.Tokens = append(l.Tokens, token)
}

func (l *lexer) line() any {
	return 3 //TODO
}

func Tokenize(source string) []Token {
	lex := createLexer(source)

	for !lex.atEof() {
		matched := false

		for _, pattern := range lex.patterns {
			if lex.insidePHP && pattern.outSideOfPHP {
				continue
			}
			loc := pattern.regex.FindStringIndex(lex.remainder())

			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)
				matched = true
				break
			}
		}

		if !matched {
			remainder := lex.remainder()
			if len(remainder) > maxErrorShowingLength {
				remainder = remainder[:maxErrorShowingLength] + "..."
			}

			panic(fmt.Sprintf("lexer error. unexpected token identifier near %s at line %s\n", remainder, lex.line()))
		}

	}

	lex.push(NewToken(TEOF, "EOF"))
	return lex.Tokens
}

func defaultHandler(kind Kind, value string) regexHandler {
	return func(lex *lexer, regex *regexp.Regexp) {
		lex.advance(len(value))
		lex.push(NewToken(kind, value))
	}
}

func htmlHandler(lex *lexer, regex *regexp.Regexp) {
	if lex.insidePHP {
		return
	}

	html := regex.FindStringSubmatch(lex.remainder())
	if len(html) != 2 {
		return
	}

	lex.push(NewToken(TInlineHtml, html[1]))
	lex.advance(len(html[1]))
	lex.insidePHP = true
}
func openTagHandler(lex *lexer, regex *regexp.Regexp) {
	tag := regex.FindString(lex.remainder())
	lex.push(NewToken(TOpenTag, tag))
	lex.advance(len(tag))
}
func closeTagHandler(lex *lexer, regex *regexp.Regexp) {
	tag := regex.FindString(lex.remainder())
	lex.push(NewToken(TCloseTag, tag))
	lex.advance(len(tag))
	lex.insidePHP = false
}

func skipHandler(lex *lexer, regex *regexp.Regexp) {
	tag := regex.FindString(lex.remainder())
	lex.advance(len(tag))
}

func createLexer(source string) *lexer {
	return &lexer{
		Tokens: make([]Token, 0),
		source: source,
		pos:    0,
		patterns: []regexPattern{
			{regexp.MustCompile("([\\S\\s]*?)?<\\?"), htmlHandler, true},
			{regexp.MustCompile("<\\?p?h?p?"), openTagHandler, false},
			{regexp.MustCompile("\\?>"), closeTagHandler, false},
			{regexp.MustCompile("\\s+"), skipHandler, false},
			//{regexp.MustCompile(""), defaultHandler(TVariable, ""),false},
			//{regexp.MustCompile(""), defaultHandler(TOpenTagWithEcho, "<?="),false},
			//{regexp.MustCompile(""), defaultHandler(TCloseTag, "?>"),false},
			//{regexp.MustCompile(""), defaultHandler(TNoElse, ""),false},
			//{regexp.MustCompile(""), defaultHandler(TLNumber, ""),false},
			//{regexp.MustCompile(""), defaultHandler(TDNumber, ""),false},
			//{regexp.MustCompile(""), defaultHandler(TNameFullyQualified, ""),false},
			//{regexp.MustCompile(""), defaultHandler(TNameRelative, ""),false},
			//{regexp.MustCompile(""), defaultHandler(TStringVarName, ""),false},
			//{regexp.MustCompile(""), defaultHandler(TNumString, ""),false},
			//{regexp.MustCompile(""), defaultHandler(TComment, "comment"),false},
			//{regexp.MustCompile(""), defaultHandler(TDocComment, "doc comment"),false},
			//{regexp.MustCompile(""), defaultHandler(TWhitespace, "whitespace"),false},
			//{regexp.MustCompile(""), defaultHandler(TStartHeredoc, "heredoc start"),false},
			//{regexp.MustCompile(""), defaultHandler(TEndHeredoc, "heredoc end"),false},
			//{regexp.MustCompile(""), defaultHandler(TBadCharacter, "invalid character"),false},
			//{regexp.MustCompile(""), defaultHandler(TString, ""), false},
			{regexp.MustCompile("[A-z]+\\\\[A-z]+"), defaultHandler(TNameQualified, "namespaced name"), false},
			{regexp.MustCompile("\".+\""), defaultHandler(TConstantEncapsedString, "string content"), false},
			{regexp.MustCompile("eval"), defaultHandler(TEval, "eval"), false},
			{regexp.MustCompile("new"), defaultHandler(TNew, "new"), false},
			{regexp.MustCompile("exit"), defaultHandler(TExit, "exit"), false},
			{regexp.MustCompile("throw"), defaultHandler(TThrow, "throw"), false},
			{regexp.MustCompile("include_once"), defaultHandler(TIncludeOnce, "include_once"), false},
			{regexp.MustCompile("require_once"), defaultHandler(TRequireOnce, "require_once"), false},
			{regexp.MustCompile("include"), defaultHandler(TInclude, "include"), false},
			{regexp.MustCompile("require"), defaultHandler(TRequire, "require"), false},
			{regexp.MustCompile("or"), defaultHandler(TLogicalOr, "or"), false},
			{regexp.MustCompile("xor"), defaultHandler(TLogicalXor, "xor"), false},
			{regexp.MustCompile("and"), defaultHandler(TLogicalAnd, "and"), false},
			{regexp.MustCompile("print"), defaultHandler(TPrint, "print"), false},
			{regexp.MustCompile("from"), defaultHandler(TYieldFrom, "yield from"), false},
			{regexp.MustCompile("yield"), defaultHandler(TYield, "yield"), false},
			{regexp.MustCompile("=>"), defaultHandler(TDoubleArrow, "=>"), false},
			{regexp.MustCompile("\\+="), defaultHandler(TPlusEqual, "+="), false},
			{regexp.MustCompile("-="), defaultHandler(TMinusEqual, "-="), false},
			{regexp.MustCompile("\\*="), defaultHandler(TMulEqual, "*="), false},
			{regexp.MustCompile("/="), defaultHandler(TDivEqual, "/="), false},
			{regexp.MustCompile(".="), defaultHandler(TConcatEqual, ".="), false},
			{regexp.MustCompile("%="), defaultHandler(TModEqual, "%="), false},
			{regexp.MustCompile("&="), defaultHandler(TAndEqual, "&="), false},
			{regexp.MustCompile("|="), defaultHandler(TOrEqual, "|="), false},
			{regexp.MustCompile("^="), defaultHandler(TXorEqual, "^="), false},
			{regexp.MustCompile("<<="), defaultHandler(TSlEqual, "<<="), false},
			{regexp.MustCompile(">>="), defaultHandler(TSrEqual, ">>="), false},
			{regexp.MustCompile("\\*\\*="), defaultHandler(TPowEqual, "**="), false},
			{regexp.MustCompile("\\?\\?="), defaultHandler(TCoalesceEqual, "??="), false},
			{regexp.MustCompile("\\?\\?"), defaultHandler(TCoalesce, "??"), false},
			{regexp.MustCompile("\\|\\|"), defaultHandler(TBooleanOr, "||"), false},
			{regexp.MustCompile("&&"), defaultHandler(TBooleanAnd, "&&"), false},
			{regexp.MustCompile("amp"), defaultHandler(TAmpersandNotFollowedByVarOrVararg, "amp"), false},
			{regexp.MustCompile("&"), defaultHandler(TAmpersandFollowedByVarOrVararg, "&"), false},
			{regexp.MustCompile("=="), defaultHandler(TIsEqual, "=="), false},
			{regexp.MustCompile("!="), defaultHandler(TIsNotEqual, "!="), false},
			{regexp.MustCompile("==="), defaultHandler(TIsIdentical, "==="), false},
			{regexp.MustCompile("!=="), defaultHandler(TIsNotIdentical, "!=="), false},
			{regexp.MustCompile("<=>"), defaultHandler(TSpaceship, "<=>"), false},
			{regexp.MustCompile("<="), defaultHandler(TIsSmallerOrEqual, "<="), false},
			{regexp.MustCompile(">="), defaultHandler(TIsGreaterOrEqual, ">="), false},
			{regexp.MustCompile("<<"), defaultHandler(TSl, "<<"), false},
			{regexp.MustCompile(">>"), defaultHandler(TSr, ">>"), false},
			{regexp.MustCompile("instanceof"), defaultHandler(TInstanceof, "instanceof"), false},
			{regexp.MustCompile("(int)"), defaultHandler(TIntCast, "(int)"), false},
			{regexp.MustCompile("(double)"), defaultHandler(TDoubleCast, "(double)"), false},
			{regexp.MustCompile("(string)"), defaultHandler(TStringCast, "(string)"), false},
			{regexp.MustCompile("(array)"), defaultHandler(TArrayCast, "(array)"), false},
			{regexp.MustCompile("(object)"), defaultHandler(TObjectCast, "(object)"), false},
			{regexp.MustCompile("(bool)"), defaultHandler(TBoolCast, "(bool)"), false},
			{regexp.MustCompile("(unset)"), defaultHandler(TUnsetCast, "(unset)"), false},
			{regexp.MustCompile("\\*\\*"), defaultHandler(TPow, "**"), false},
			{regexp.MustCompile("clone"), defaultHandler(TClone, "clone"), false},
			{regexp.MustCompile("elseif"), defaultHandler(TElseif, "elseif"), false},
			{regexp.MustCompile("else"), defaultHandler(TElse, "else"), false},
			{regexp.MustCompile("if"), defaultHandler(TIf, "if"), false},
			{regexp.MustCompile("endif"), defaultHandler(TEndif, "endif"), false},
			{regexp.MustCompile("echo"), defaultHandler(TEcho, "echo"), false},
			{regexp.MustCompile("do"), defaultHandler(TDo, "do"), false},
			{regexp.MustCompile("while"), defaultHandler(TWhile, "while"), false},
			{regexp.MustCompile("endwhile"), defaultHandler(TEndWhile, "endwhile"), false},
			{regexp.MustCompile("for"), defaultHandler(TFor, "for"), false},
			{regexp.MustCompile("endfor"), defaultHandler(TEndFor, "endfor"), false},
			{regexp.MustCompile("foreach"), defaultHandler(TForeach, "foreach"), false},
			{regexp.MustCompile("endforeach"), defaultHandler(TEndForeach, "endforeach"), false},
			{regexp.MustCompile("declare"), defaultHandler(TDeclare, "declare"), false},
			{regexp.MustCompile("enddeclare"), defaultHandler(TEndDeclare, "enddeclare"), false},
			{regexp.MustCompile("as"), defaultHandler(TAs, "as"), false},
			{regexp.MustCompile("switch"), defaultHandler(TSwitch, "switch"), false},
			{regexp.MustCompile("endswitch"), defaultHandler(TEndSwitch, "endswitch"), false},
			{regexp.MustCompile("case"), defaultHandler(TCase, "case"), false},
			{regexp.MustCompile("default"), defaultHandler(TDefault, "default"), false},
			{regexp.MustCompile("match"), defaultHandler(TMatch, "match"), false},
			{regexp.MustCompile("break"), defaultHandler(TBreak, "break"), false},
			{regexp.MustCompile("continue"), defaultHandler(TContinue, "continue"), false},
			{regexp.MustCompile("goto"), defaultHandler(TGoto, "goto"), false},
			{regexp.MustCompile("function"), defaultHandler(TFunction, "function"), false},
			{regexp.MustCompile("fn"), defaultHandler(TFn, "fn"), false},
			{regexp.MustCompile("const"), defaultHandler(TConst, "const"), false},
			{regexp.MustCompile("return"), defaultHandler(TReturn, "return"), false},
			{regexp.MustCompile("try"), defaultHandler(TTry, "try"), false},
			{regexp.MustCompile("catch"), defaultHandler(TCatch, "catch"), false},
			{regexp.MustCompile("finally"), defaultHandler(TFinally, "finally"), false},
			{regexp.MustCompile("use"), defaultHandler(TUse, "use"), false},
			{regexp.MustCompile("insteadof"), defaultHandler(TInsteadof, "insteadof"), false},
			{regexp.MustCompile("global"), defaultHandler(TGlobal, "global"), false},
			{regexp.MustCompile("static"), defaultHandler(TStatic, "static"), false},
			{regexp.MustCompile("abstract"), defaultHandler(TAbstract, "abstract"), false},
			{regexp.MustCompile("final"), defaultHandler(TFinal, "final"), false},
			{regexp.MustCompile("private"), defaultHandler(TPrivate, "private"), false},
			{regexp.MustCompile("protected"), defaultHandler(TProtected, "protected"), false},
			{regexp.MustCompile("public"), defaultHandler(TPublic, "public"), false},
			{regexp.MustCompile("readonly"), defaultHandler(TReadonly, "readonly"), false},
			{regexp.MustCompile("var"), defaultHandler(TVar, "var"), false},
			{regexp.MustCompile("unset"), defaultHandler(TUnset, "unset"), false},
			{regexp.MustCompile("isset"), defaultHandler(TIsset, "isset"), false},
			{regexp.MustCompile("empty"), defaultHandler(TEmpty, "empty"), false},
			{regexp.MustCompile("__halt_compiler"), defaultHandler(THaltCompiler, "__halt_compiler"), false},
			{regexp.MustCompile("class"), defaultHandler(TClass, "class"), false},
			{regexp.MustCompile("trait"), defaultHandler(TTrait, "trait"), false},
			{regexp.MustCompile("interface"), defaultHandler(TInterface, "interface"), false},
			{regexp.MustCompile("enum"), defaultHandler(TEnum, "enum"), false},
			{regexp.MustCompile("extends"), defaultHandler(TExtends, "extends"), false},
			{regexp.MustCompile("implements"), defaultHandler(TImplements, "implements"), false},
			{regexp.MustCompile("namespace"), defaultHandler(TNamespace, "namespace"), false},
			{regexp.MustCompile("list"), defaultHandler(TList, "list"), false},
			{regexp.MustCompile("array"), defaultHandler(TArray, "array"), false},
			{regexp.MustCompile("callable"), defaultHandler(TCallable, "callable"), false},
			{regexp.MustCompile("__LINE__"), defaultHandler(TLine, "__LINE__"), false},
			{regexp.MustCompile("__FILE__"), defaultHandler(TFile, "__FILE__"), false},
			{regexp.MustCompile("__DIR__"), defaultHandler(TDir, "__DIR__"), false},
			{regexp.MustCompile("__CLASS__"), defaultHandler(TClassC, "__CLASS__"), false},
			{regexp.MustCompile("__TRAIT__"), defaultHandler(TTraitC, "__TRAIT__"), false},
			{regexp.MustCompile("__METHOD__"), defaultHandler(TMethodC, "__METHOD__"), false},
			{regexp.MustCompile("__FUNCTION__"), defaultHandler(TFuncC, "__FUNCTION__"), false},
			{regexp.MustCompile("__PROPERTY__"), defaultHandler(TPropertyC, "__PROPERTY__"), false},
			{regexp.MustCompile("__NAMESPACE__"), defaultHandler(TNsC, "__NAMESPACE__"), false},
			{regexp.MustCompile("#\\["), defaultHandler(TAttribute, "#["), false},
			{regexp.MustCompile("\\+\\+"), defaultHandler(TInc, "++"), false},
			{regexp.MustCompile("--"), defaultHandler(TDec, "--"), false},
			{regexp.MustCompile("->"), defaultHandler(TObjectOperator, "->"), false},
			{regexp.MustCompile("\\?->"), defaultHandler(TNullSafeObjectOperator, "?->"), false},
			{regexp.MustCompile("\\${"), defaultHandler(TDollarOpenCurlyBraces, "${"), false},
			{regexp.MustCompile("{\\$"), defaultHandler(TCurlyOpen, "{$"), false},
			{regexp.MustCompile("::"), defaultHandler(TPaamayimNekudotayim, "::"), false},
			{regexp.MustCompile("\\\\"), defaultHandler(TNsSeparator, "\\\\"), false},
			{regexp.MustCompile("..."), defaultHandler(TEllipsis, "..."), false},
		},
	}
}
