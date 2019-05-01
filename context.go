package yin

import (
	"encoding/json"
	"io"
	"net/http"
)

type BodyContent interface{
	Get(string)string
}

// yin框架上下文
type Context interface {
	HTML(int, string)
	JSON(int, interface{})
	FILE(int, string)
	SUCCESS(map[string]interface{})
	Body()BodyContent
}

type coreContext struct {
	r *http.Request
	w http.ResponseWriter
	body BodyContent
}

// 创建上下文
func NewContext(r *http.Request, w http.ResponseWriter) Context {
	context := &coreContext{
		r: r,
		w: w,
		body:nil,
	}
	w.Header().Set("Server", "yin v1.0.1")
	context.parseBody()
	return context
}

func (core *coreContext) HTML(stateCode int, content string) {
	core.w.Header().Set("Content-Type", "text/html; charset=utf-8")
	core.w.WriteHeader(stateCode)
	core.w.Write([]byte(content))
}

func (core *coreContext) JSON(stateCode int, content interface{}) {
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		// TODO 500错误处理
		return
	}
	core.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	core.w.WriteHeader(stateCode)
	core.w.Write(jsonBytes)
}

func (core *coreContext) SUCCESS(content map[string]interface{}) {
	if content == nil{
		content = make(map[string]interface{})
	}
	content["status"] = "ok"
	core.JSON(200,content)
}

func (core *coreContext) FILE(stateCode int, fileName string) {
	// TODO 返回文件下载
}

func (core *coreContext) Body()BodyContent{
	return core.body
}

// 解析主体 post from json get
func (core *coreContext) parseBody(){
	// post json解析
	if core.r.Method == "POST" && core.r.Header.Get("Content-Type") == "application/json"{
		err := core.r.ParseForm()
		if err != nil{
			return
		}
		core.body = NewBodyContent(core.r.Body)

	}else if core.r.Method == "GET"{
		core.body = core.r.URL.Query()
	}
}

type jsonValues struct{
	formData map[string]interface{}
}

func NewBodyContent(reader io.Reader)BodyContent{
	c := &jsonValues{
		formData:make(map[string]interface{}),
	}
	json.NewDecoder(reader).Decode(&c.formData)
	return c
}

func (core *jsonValues)Get(key string)string{
	if v,ok := core.formData[key];ok{
		return v.(string)
	}
	return ""
}