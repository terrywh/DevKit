package client

type StreamFile struct {
	Path string
}

func (sf *StreamFile) Request() stream.Request {
	return &stream.StreamFileReq{}
}

type FetchFile struct{}

func (ff *FetchFile) ServeStream() {}
