package yin

// yin框架上下文
type Context interface {
	HTML(int, string)
	JSON(int, interface{})
	FILE(int, string)
}
