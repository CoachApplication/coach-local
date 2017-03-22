package configprovider

import "io"

type FileConnector struct {
}

func (fc *FileConnector) Scopes() []string {
	return []string{}
}
func (fc *FileConnector) Keys() []string {
	return []string{}
}
func (fc *FileConnector) Get(key, scope string) (io.ReadCloser, error) {
	return nil, nil
}
func (fc *FileConnector) Set(key, scope string, source io.ReadCloser) error {
	return nil
}
