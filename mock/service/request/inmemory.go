package request

import "context"

type inmemory struct {
	errs map[string]error
	data []byte
}

func NewInmemory(errs map[string]error, data []byte) *inmemory {
	return &inmemory{
		errs: errs,
		data: data,
	}
}

func (c *inmemory) Send(ctx context.Context, url string) ([]byte, error) {
	err, exists := c.errs[url]
	if exists {
		return nil, err
	}

	return c.data, nil
}
