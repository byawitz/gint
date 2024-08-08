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

func createLexer(source string) *lexer {
	return &lexer{
		Tokens: make([]Token, 0),
		source: source,
		pos:    0,
		patterns: []regexPattern{
			{regexp.MustCompile("([\\S\\s]*?)?<\\?"), htmlHandler, true},
			{regex: regexp.MustCompile("<\\?p?h?p?"), handler: openTagHandler},
			{regex: regexp.MustCompile("<\\?="), handler: staticHandler(TOpenTagWithEcho, "<?=")},
			{regex: regexp.MustCompile("\\?>"), handler: closeTagHandler},
			{regex: regexp.MustCompile("\\s+"), handler: skipHandler},
			{regex: regexp.MustCompile("/\\*+(.|[\\r\\n])*?\\*/"), handler: docCommentHandler},
			{regex: regexp.MustCompile("/\\*[^*](.|[\\r\\n])*?\\*/"), handler: commentHandler},
			{regex: regexp.MustCompile("//.*"), handler: commentHandler},
			{regex: regexp.MustCompile("[0-9]+\\.[0-9]+"), handler: literalHandler(TDNumber)},
			{regex: regexp.MustCompile("[0-9]+"), handler: literalHandler(TLNumber)},
			{regex: regexp.MustCompile("\\$[A-z_]?[A-z_0-9]+"), handler: literalHandler(TVar)},
			{regex: regexp.MustCompile(`;`), handler: staticHandler(TString, ";")},
			{regex: regexp.MustCompile(`\[`), handler: staticHandler(TOpenBracket, "[")},
			{regex: regexp.MustCompile(`]`), handler: staticHandler(TCloseBracket, "]")},
			{regex: regexp.MustCompile(`\{`), handler: staticHandler(TOPENCurly, "{")},
			{regex: regexp.MustCompile(`}`), handler: staticHandler(TCloseCurly, "}")},
			{regex: regexp.MustCompile(`\(`), handler: staticHandler(TOpenParen, "(")},
			{regex: regexp.MustCompile(`\)`), handler: staticHandler(TCloseParen, ")")},
			//{regex: regexp.MustCompile(""),handler: staticHandler(TNoElse, "},
			//{regex: regexp.MustCompile(""),handler: staticHandler(TNameFullyQualified, "},
			//{regex: regexp.MustCompile(""),handler: staticHandler(TNameRelative, "},
			//{regex: regexp.MustCompile(""),handler: staticHandler(TStringVarName, "},
			//{regex: regexp.MustCompile(""),handler: staticHandler(TNumString, "},
			//{regex: regexp.MustCompile(""),handler: staticHandler(TWhitespace, "whitespace},
			//{regex: regexp.MustCompile(""),handler: staticHandler(TStartHeredoc, "heredoc start},
			//{regex: regexp.MustCompile(""),handler: staticHandler(TEndHeredoc, "heredoc end},
			//{regex: regexp.MustCompile(""),handler: staticHandler(TBadCharacter, "invalid character},

			{regex: regexp.MustCompile("[A-z]+\\\\[A-z]+"), handler: staticHandler(TNameQualified, "namespaced name")},
			{regex: regexp.MustCompile("\"[^\"]+\""), handler: stringHandler},
			{regex: regexp.MustCompile("'[^']+'"), handler: stringHandler},
			{regex: regexp.MustCompile("eval"), handler: staticHandler(TEval, "eval")},
			{regex: regexp.MustCompile("new"), handler: staticHandler(TNew, "new")},
			{regex: regexp.MustCompile("exit"), handler: staticHandler(TExit, "exit")},
			{regex: regexp.MustCompile("throw"), handler: staticHandler(TThrow, "throw")},
			{regex: regexp.MustCompile("include_once"), handler: staticHandler(TIncludeOnce, "include_once")},
			{regex: regexp.MustCompile("require_once"), handler: staticHandler(TRequireOnce, "require_once")},
			{regex: regexp.MustCompile("include"), handler: staticHandler(TInclude, "include")},
			{regex: regexp.MustCompile("require"), handler: staticHandler(TRequire, "require")},
			{regex: regexp.MustCompile("or"), handler: staticHandler(TLogicalOr, "or")},
			{regex: regexp.MustCompile("xor"), handler: staticHandler(TLogicalXor, "xor")},
			{regex: regexp.MustCompile("and"), handler: staticHandler(TLogicalAnd, "and")},
			{regex: regexp.MustCompile("print"), handler: staticHandler(TPrint, "print")},
			{regex: regexp.MustCompile("from"), handler: staticHandler(TYieldFrom, "yield from")},
			{regex: regexp.MustCompile("yield"), handler: staticHandler(TYield, "yield")},
			{regex: regexp.MustCompile("=>"), handler: staticHandler(TDoubleArrow, "=>")},
			{regex: regexp.MustCompile("\\+="), handler: staticHandler(TPlusEqual, "+=")},
			{regex: regexp.MustCompile("-="), handler: staticHandler(TMinusEqual, "-=")},
			{regex: regexp.MustCompile("\\*="), handler: staticHandler(TMulEqual, "*=")},
			{regex: regexp.MustCompile("/="), handler: staticHandler(TDivEqual, "/=")},
			{regex: regexp.MustCompile(".="), handler: staticHandler(TConcatEqual, ".=")},
			{regex: regexp.MustCompile("%="), handler: staticHandler(TModEqual, "%=")},
			{regex: regexp.MustCompile("&="), handler: staticHandler(TAndEqual, "&=")},
			{regex: regexp.MustCompile("\\|="), handler: staticHandler(TOrEqual, "|=")},
			{regex: regexp.MustCompile("\\^="), handler: staticHandler(TXorEqual, "^=")},
			{regex: regexp.MustCompile("<<="), handler: staticHandler(TSlEqual, "<<=")},
			{regex: regexp.MustCompile(">>="), handler: staticHandler(TSrEqual, ">>=")},
			{regex: regexp.MustCompile("\\*\\*="), handler: staticHandler(TPowEqual, "**=")},
			{regex: regexp.MustCompile("\\?\\?="), handler: staticHandler(TCoalesceEqual, "??=")},
			{regex: regexp.MustCompile("\\?\\?"), handler: staticHandler(TCoalesce, "??")},
			{regex: regexp.MustCompile("\\|\\|"), handler: staticHandler(TBooleanOr, "||")},
			{regex: regexp.MustCompile("&&"), handler: staticHandler(TBooleanAnd, "&&")},
			{regex: regexp.MustCompile("amp"), handler: staticHandler(TAmpersandNotFollowedByVarOrVararg, "amp")},
			{regex: regexp.MustCompile("&"), handler: staticHandler(TAmpersandFollowedByVarOrVararg, "&")},
			{regex: regexp.MustCompile("=="), handler: staticHandler(TIsEqual, "==")},
			{regex: regexp.MustCompile("!="), handler: staticHandler(TIsNotEqual, "!=")},
			{regex: regexp.MustCompile("!"), handler: staticHandler(TNotEqual, "!")},
			{regex: regexp.MustCompile("==="), handler: staticHandler(TIsIdentical, "===")},
			{regex: regexp.MustCompile("!=="), handler: staticHandler(TIsNotIdentical, "!==")},
			{regex: regexp.MustCompile("<=>"), handler: staticHandler(TSpaceship, "<=>")},
			{regex: regexp.MustCompile("<="), handler: staticHandler(TIsSmallerOrEqual, "<=")},
			{regex: regexp.MustCompile(">="), handler: staticHandler(TIsGreaterOrEqual, ">=")},
			{regex: regexp.MustCompile("<<"), handler: staticHandler(TSl, "<<")},
			{regex: regexp.MustCompile(">>"), handler: staticHandler(TSr, ">>")},
			{regex: regexp.MustCompile(`=`), handler: staticHandler(TAssignment, "=")},
			{regex: regexp.MustCompile(`\+`), handler: staticHandler(TPlus, "+")},
			{regex: regexp.MustCompile(`-`), handler: staticHandler(TDash, "-")},
			{regex: regexp.MustCompile(`/`), handler: staticHandler(TSlash, "/")},
			{regex: regexp.MustCompile(`\*`), handler: staticHandler(TStar, "*")},
			{regex: regexp.MustCompile(`%`), handler: staticHandler(TPercent, "%")},
			{regex: regexp.MustCompile("instanceof"), handler: staticHandler(TInstanceof, "instanceof")},
			{regex: regexp.MustCompile("(int)"), handler: staticHandler(TIntCast, "(int)")},
			{regex: regexp.MustCompile("(double)"), handler: staticHandler(TDoubleCast, "(double)")},
			{regex: regexp.MustCompile("(string)"), handler: staticHandler(TStringCast, "(string)")},
			{regex: regexp.MustCompile("(array)"), handler: staticHandler(TArrayCast, "(array)")},
			{regex: regexp.MustCompile("(object)"), handler: staticHandler(TObjectCast, "(object)")},
			{regex: regexp.MustCompile("(bool)"), handler: staticHandler(TBoolCast, "(bool)")},
			{regex: regexp.MustCompile("(unset)"), handler: staticHandler(TUnsetCast, "(unset)")},
			{regex: regexp.MustCompile("\\*\\*"), handler: staticHandler(TPow, "**")},
			{regex: regexp.MustCompile("clone"), handler: staticHandler(TClone, "clone")},
			{regex: regexp.MustCompile("elseif"), handler: staticHandler(TElseif, "elseif")},
			{regex: regexp.MustCompile("else"), handler: staticHandler(TElse, "else")},
			{regex: regexp.MustCompile("if"), handler: staticHandler(TIf, "if")},
			{regex: regexp.MustCompile("endif"), handler: staticHandler(TEndif, "endif")},
			{regex: regexp.MustCompile("echo"), handler: staticHandler(TEcho, "echo")},
			{regex: regexp.MustCompile("do"), handler: staticHandler(TDo, "do")},
			{regex: regexp.MustCompile("while"), handler: staticHandler(TWhile, "while")},
			{regex: regexp.MustCompile("endwhile"), handler: staticHandler(TEndWhile, "endwhile")},
			{regex: regexp.MustCompile("for"), handler: staticHandler(TFor, "for")},
			{regex: regexp.MustCompile("endfor"), handler: staticHandler(TEndFor, "endfor")},
			{regex: regexp.MustCompile("foreach"), handler: staticHandler(TForeach, "foreach")},
			{regex: regexp.MustCompile("endforeach"), handler: staticHandler(TEndForeach, "endforeach")},
			{regex: regexp.MustCompile("declare"), handler: staticHandler(TDeclare, "declare")},
			{regex: regexp.MustCompile("enddeclare"), handler: staticHandler(TEndDeclare, "enddeclare")},
			{regex: regexp.MustCompile("as"), handler: staticHandler(TAs, "as")},
			{regex: regexp.MustCompile("switch"), handler: staticHandler(TSwitch, "switch")},
			{regex: regexp.MustCompile("endswitch"), handler: staticHandler(TEndSwitch, "endswitch")},
			{regex: regexp.MustCompile("case"), handler: staticHandler(TCase, "case")},
			{regex: regexp.MustCompile("default"), handler: staticHandler(TDefault, "default")},
			{regex: regexp.MustCompile("match"), handler: staticHandler(TMatch, "match")},
			{regex: regexp.MustCompile("break"), handler: staticHandler(TBreak, "break")},
			{regex: regexp.MustCompile("continue"), handler: staticHandler(TContinue, "continue")},
			{regex: regexp.MustCompile("goto"), handler: staticHandler(TGoto, "goto")},
			{regex: regexp.MustCompile("function"), handler: staticHandler(TFunction, "function")},
			{regex: regexp.MustCompile("fn"), handler: staticHandler(TFn, "fn")},
			{regex: regexp.MustCompile("const"), handler: staticHandler(TConst, "const")},
			{regex: regexp.MustCompile("return"), handler: staticHandler(TReturn, "return")},
			{regex: regexp.MustCompile("try"), handler: staticHandler(TTry, "try")},
			{regex: regexp.MustCompile("catch"), handler: staticHandler(TCatch, "catch")},
			{regex: regexp.MustCompile("finally"), handler: staticHandler(TFinally, "finally")},
			{regex: regexp.MustCompile("use"), handler: staticHandler(TUse, "use")},
			{regex: regexp.MustCompile("insteadof"), handler: staticHandler(TInsteadof, "insteadof")},
			{regex: regexp.MustCompile("global"), handler: staticHandler(TGlobal, "global")},
			{regex: regexp.MustCompile("static"), handler: staticHandler(TStatic, "static")},
			{regex: regexp.MustCompile("abstract"), handler: staticHandler(TAbstract, "abstract")},
			{regex: regexp.MustCompile("final"), handler: staticHandler(TFinal, "final")},
			{regex: regexp.MustCompile("private"), handler: staticHandler(TPrivate, "private")},
			{regex: regexp.MustCompile("protected"), handler: staticHandler(TProtected, "protected")},
			{regex: regexp.MustCompile("public"), handler: staticHandler(TPublic, "public")},
			{regex: regexp.MustCompile("readonly"), handler: staticHandler(TReadonly, "readonly")},
			{regex: regexp.MustCompile("var"), handler: staticHandler(TVar, "var")},
			{regex: regexp.MustCompile("unset"), handler: staticHandler(TUnset, "unset")},
			{regex: regexp.MustCompile("isset"), handler: staticHandler(TIsset, "isset")},
			{regex: regexp.MustCompile("empty"), handler: staticHandler(TEmpty, "empty")},
			{regex: regexp.MustCompile("__halt_compiler"), handler: staticHandler(THaltCompiler, "__halt_compiler")},
			{regex: regexp.MustCompile("class"), handler: staticHandler(TClass, "class")},
			{regex: regexp.MustCompile("trait"), handler: staticHandler(TTrait, "trait")},
			{regex: regexp.MustCompile("interface"), handler: staticHandler(TInterface, "interface")},
			{regex: regexp.MustCompile("enum"), handler: staticHandler(TEnum, "enum")},
			{regex: regexp.MustCompile("extends"), handler: staticHandler(TExtends, "extends")},
			{regex: regexp.MustCompile("implements"), handler: staticHandler(TImplements, "implements")},
			{regex: regexp.MustCompile("namespace"), handler: staticHandler(TNamespace, "namespace")},
			{regex: regexp.MustCompile("list"), handler: staticHandler(TList, "list")},
			{regex: regexp.MustCompile("array"), handler: staticHandler(TArray, "array")},
			{regex: regexp.MustCompile("callable"), handler: staticHandler(TCallable, "callable")},
			{regex: regexp.MustCompile("__LINE__"), handler: staticHandler(TLine, "__LINE__")},
			{regex: regexp.MustCompile("__FILE__"), handler: staticHandler(TFile, "__FILE__")},
			{regex: regexp.MustCompile("__DIR__"), handler: staticHandler(TDir, "__DIR__")},
			{regex: regexp.MustCompile("__CLASS__"), handler: staticHandler(TClassC, "__CLASS__")},
			{regex: regexp.MustCompile("__TRAIT__"), handler: staticHandler(TTraitC, "__TRAIT__")},
			{regex: regexp.MustCompile("__METHOD__"), handler: staticHandler(TMethodC, "__METHOD__")},
			{regex: regexp.MustCompile("__FUNCTION__"), handler: staticHandler(TFuncC, "__FUNCTION__")},
			{regex: regexp.MustCompile("__PROPERTY__"), handler: staticHandler(TPropertyC, "__PROPERTY__")},
			{regex: regexp.MustCompile("__NAMESPACE__"), handler: staticHandler(TNsC, "__NAMESPACE__")},
			{regex: regexp.MustCompile("#\\["), handler: staticHandler(TAttribute, "#[")},
			{regex: regexp.MustCompile("\\+\\+"), handler: staticHandler(TInc, "++")},
			{regex: regexp.MustCompile("--"), handler: staticHandler(TDec, "--")},
			{regex: regexp.MustCompile("->"), handler: staticHandler(TObjectOperator, "->")},
			{regex: regexp.MustCompile("\\?->"), handler: staticHandler(TNullSafeObjectOperator, "?->")},
			{regex: regexp.MustCompile("\\${"), handler: staticHandler(TDollarOpenCurlyBraces, "${")},
			{regex: regexp.MustCompile("{\\$"), handler: staticHandler(TCurlyOpen, "{$")},
			{regex: regexp.MustCompile("::"), handler: staticHandler(TPaamayimNekudotayim, "::")},
			{regex: regexp.MustCompile("\\\\"), handler: staticHandler(TNsSeparator, "\\\\")},
			{regex: regexp.MustCompile("\\.\\.\\."), handler: staticHandler(TEllipsis, "...")},
			{regex: regexp.MustCompile("[a-zA-Z_\\x80-\\xff][a-zA-Z0-9_\\x80-\\xff]*"), handler: literalHandler(TString)},
		},
	}
}

func staticHandler(kind Kind, value string) regexHandler {
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
func docCommentHandler(lex *lexer, regex *regexp.Regexp) {

	docComment := regex.FindStringSubmatch(lex.remainder())
	if len(docComment) != 2 {
		return
	}

	lex.push(NewToken(TDocComment, docComment[0]))
	lex.advance(len(docComment[0]))
}
func commentHandler(lex *lexer, regex *regexp.Regexp) {
	comment := regex.FindString(lex.remainder())

	lex.push(NewToken(TComment, comment))
	lex.advance(len(comment))
}
func stringHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	stringLiteral := lex.remainder()[match[0]:match[1]]

	lex.push(NewToken(TString, stringLiteral))
	lex.advance(len(stringLiteral))
}
func literalHandler(kind Kind) regexHandler {
	return func(lex *lexer, regex *regexp.Regexp) {
		match := regex.FindString(lex.remainder())
		lex.push(NewToken(kind, match))
		lex.advance(len(match))
	}
}
