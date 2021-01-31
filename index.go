package uuid

// IUUID 通知唯一标识符接口
type IUUID interface{
	Generate() string
}