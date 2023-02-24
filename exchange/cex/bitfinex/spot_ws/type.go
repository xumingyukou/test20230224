package spot_ws

type APIConf struct {
	ProxyUrl    string
	EndPoint    string
	ReadTimeout int64 //second
	AccessKey   string
	SecretKey   string
	Passphrase  string
	IsTest      bool
	IsPrivate   bool
}

type WSRequest struct {
	Channel string `json:"channel"`
}

type SubMsg struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Symbol  string `json:"symbol"`
}

type BookSubMsg struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Symbol  string `json:"symbol"`
	Prec    string `json:"prec"`
	Freq    string `json:"freq"`
	Len     string `json:"len"`
}

type UnsubMsg struct {
	Event  string `json:"event"`
	ChanID int64  `json:"chanId"`
}

type initialResp struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	ChanID  int64  `json:"chanId"`
	Symbol  string `json:"symbol"`
	Pair    string `json:"pair"`
	Code    int64  `json:"code"`
	Msg     int64  `json:"msg"`
}
type Response []interface{}

type ConfigureMsg struct {
	Event string `json:"event"`
	Flags int64  `json:"flags"`
}
