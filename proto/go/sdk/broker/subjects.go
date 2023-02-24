package broker

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

type Writer interface {
	Write(*bytes.Buffer)
	Set(string)
}

type writer string

func (s *writer) Write(buffer *bytes.Buffer) {
	buffer.WriteString(string(*s))
}

func (s *writer) Set(str string) {
	if string(*s) != str {
		panic("Invalid subject extracting, from " + str + " to " + string(*s))
	}
}

// 描述包含var变量的字符串，只支持最多一个变量
type varWriter struct {
	parts       []*string
	varStartPos int
	varEndPos   int
	str         string
}

func (s *varWriter) Write(buffer *bytes.Buffer) {
	for _, e := range s.parts {
		buffer.WriteString(*e)
	}
}

func (s *varWriter) Set(str string) {
	// 根据str反向解析到变量值
	i := s.varStartPos
	if i > 0 {
		tmp := len(str)
		str = strings.TrimPrefix(str, s.str[:i])
		if len(str) == tmp {
			panic("Trim prefix failed: " + str)
		}
		i = 1
	}

	if s.varEndPos < len(s.str) {
		tmp := len(str)
		str = strings.TrimSuffix(str, s.str[s.varEndPos:])
		if tmp == len(str) {
			panic("Trim suffix failed: " + str)
		}
	}

	if i >= 0 {
		*s.parts[i] = str
	}
}

func (s *varWriter) Init(str string, ctx *SubjectContext) {
	s.str = str
	// 解析str内容
	re, _ := regexp.Compile("{{[^{}]+}}")
	locs := re.FindAllStringIndex(str, -1)

	s.parts = make([]*string, 0)
	lastPos := 0
	s.varStartPos = -1
	s.varEndPos = len(str)

	for _, loc := range locs {
		s.varStartPos = loc[0]
		s.varEndPos = loc[1]
		e := str[loc[0]:loc[1]]
		tmp := s.str[lastPos:loc[0]]
		if tmp != "" {
			s.parts = append(s.parts, &tmp)
		}

		if e == "{{ProjectId}}" {
			s.parts = append(s.parts, &ctx.projectId)
		} else if e == "{{Exchange}}" {
			s.parts = append(s.parts, &ctx.exchange)
		} else if e == "{{AccountId}}" {
			s.parts = append(s.parts, &ctx.accountId)
		} else if e == "{{Market}}" {
			s.parts = append(s.parts, &ctx.market)
		} else if e == "{{SymbolType}}" {
			s.parts = append(s.parts, &ctx.symbolType)
		} else if e == "{{Symbol}}" {
			s.parts = append(s.parts, &ctx.symbol)
		} else if e == "{{Token}}" {
			s.parts = append(s.parts, &ctx.token)
		} else {
			panic("Invalid subject variable name: " + e)
		}
		lastPos = loc[1]
	}

	// 每个token最多一个变量
	if len(locs) > 1 {
		panic("At most 1 var defined in a token, there are: " + str)
	}

	if lastPos < len(str) {
		tmp := s.str[lastPos:]
		s.parts = append(s.parts, &tmp)
	}
}

func (s *varWriter) InitConf(str string, ctx *SubjectConfContext) {
	s.str = str
	// 解析str内容
	re, _ := regexp.Compile("{{[^{}]+}}")
	locs := re.FindAllStringIndex(str, -1)

	s.parts = make([]*string, 0)
	lastPos := 0
	s.varStartPos = -1
	s.varEndPos = len(str)

	for _, loc := range locs {
		s.varStartPos = loc[0]
		s.varEndPos = loc[1]
		e := str[loc[0]:loc[1]]
		tmp := s.str[lastPos:loc[0]]
		if tmp != "" {
			s.parts = append(s.parts, &tmp)
		}

		if e == "{{Env}}" {
			s.parts = append(s.parts, &ctx.env)
		} else if e == "{{Service}}" {
			s.parts = append(s.parts, &ctx.service)
		} else {
			panic("Invalid subject variable name: " + e)
		}
		lastPos = loc[1]
	}

	// 每个token最多一个变量
	if len(locs) > 1 {
		panic("At most 1 var defined in a token, there are: " + str)
	}

	if lastPos < len(str) {
		tmp := s.str[lastPos:]
		s.parts = append(s.parts, &tmp)
	}
}

/*
 封装与broker通信相关的辅助函数，主要用于操作subject字符串
*/

// subjectContext用于处理包含变量的subject定义
// 目前只支持projectId,exchange,market,symbolType,symbol,accountId,token6种变量
// 此结构非线程安全
type SubjectContext struct {
	projectId  string
	exchange   string
	market     string
	symbolType string
	symbol     string
	accountId  string
	token      string

	//包含string或者varString
	elems  []Writer
	buffer bytes.Buffer
	tpl    string
}

func (s *SubjectContext) GetProjectId() string {
	return s.projectId
}

func (s *SubjectContext) SetProjectId(projectId string) {
	s.projectId = projectId
}

func (s *SubjectContext) GetExchange() int {
	if ret, err := strconv.Atoi(s.exchange); err == nil {
		return ret
	} else {
		panic(err)
	}
}

func (s *SubjectContext) SetExchange(exchange int) {
	s.exchange = strconv.Itoa(exchange)
}

func (s *SubjectContext) GetMarket() int {
	if ret, err := strconv.Atoi(s.market); err == nil {
		return ret
	} else {
		panic(err)
	}
}

func (s *SubjectContext) SetMarket(market int) {
	s.market = strconv.Itoa(market)
}

func (s *SubjectContext) GetSymbolType() int {
	if ret, err := strconv.Atoi(s.symbolType); err == nil {
		return ret
	} else {
		panic(err)
	}
}

func (s *SubjectContext) SetSymbolType(symbolType int) {
	s.symbolType = strconv.Itoa(symbolType)
}

func (s *SubjectContext) GetSymbol() string {
	return s.symbol
}

func (s *SubjectContext) SetSymbol(symbol string) {
	s.symbol = symbol
}

func (s *SubjectContext) GetAccountId() string {
	return s.accountId
}

func (s *SubjectContext) SetAccountId(accoountId string) {
	s.accountId = accoountId
}

func (s *SubjectContext) GetToken() string {
	return s.token
}

func (s *SubjectContext) SetToken(token string) {
	s.token = token
}

func (s *SubjectContext) GetTpl() string {
	return s.tpl
}

func (s *SubjectContext) Init(tpl string) {
	s.tpl = tpl
	s.buffer.Reset()
	s.elems = make([]Writer, 0)
	tokens := strings.Split(tpl, ".")

	for _, token := range tokens {
		if strings.ContainsAny(token, "{}") {
			v := varWriter{}
			v.Init(token, s)
			s.elems = append(s.elems, &v)
		} else {
			t := writer(token)
			s.elems = append(s.elems, &t)
		}
	}
}

// 根据tpl和变量值生成subject
func (s *SubjectContext) ToSubject() string {
	s.buffer.Reset()
	next := false
	for _, e := range s.elems {
		if next {
			s.buffer.WriteByte('.')
		} else {
			next = true
		}
		e.Write(&s.buffer)
	}
	return s.buffer.String()
}

// 将subject解析到tpl中包含的变量中，后续可以调用GetXXX函数来获取解析后的值
func (s *SubjectContext) FromSubject(subject string) {
	elems := strings.Split(subject, ".")
	if len(elems) != len(s.elems) {
		panic("subject token number NOT MATCH tpl token number")
	}

	for i := 0; i < len(elems); i = i + 1 {
		s.elems[i].Set(elems[i])
	}
}

// 配置相关subject context，目前只支持Env和Service两种变量
type SubjectConfContext struct {
	env     string
	service string

	//包含string或者varString
	elems  []Writer
	buffer bytes.Buffer
	tpl    string
}

func (s *SubjectConfContext) GetEnv() string {
	return s.env
}

func (s *SubjectConfContext) SetEnv(env string) {
	s.env = env
}

func (s *SubjectConfContext) GetService() string {
	return s.service
}

func (s *SubjectConfContext) SetService(service string) {
	s.service = service
}

func (s *SubjectConfContext) Init(tpl string) {
	s.tpl = tpl
	s.buffer.Reset()
	s.elems = make([]Writer, 0)
	tokens := strings.Split(tpl, ".")

	for _, token := range tokens {
		if strings.ContainsAny(token, "{}") {
			v := varWriter{}
			v.InitConf(token, s)
			s.elems = append(s.elems, &v)
		} else {
			t := writer(token)
			s.elems = append(s.elems, &t)
		}
	}
}

// 根据tpl和变量值生成subject
func (s *SubjectConfContext) ToSubject() string {
	s.buffer.Reset()
	next := false
	for _, e := range s.elems {
		if next {
			s.buffer.WriteByte('.')
		} else {
			next = true
		}
		e.Write(&s.buffer)
	}
	return s.buffer.String()
}

func (s *SubjectConfContext) FromSubject(subject string) {
	elems := strings.Split(subject, ".")
	if len(elems) != len(s.elems) {
		panic("subject token number NOT MATCH tpl token number")
	}

	for i := 0; i < len(elems); i = i + 1 {
		s.elems[i].Set(elems[i])
	}
}
