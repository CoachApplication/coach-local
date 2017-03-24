package local

import (
	"testing"

	base_configprovider "github.com/james-nesbitt/coach-base/handler/configprovider"
)

// Message JSON bytes
var MessageBytes = []byte(`Name: Alice
Body: Hello
Time: 1294706395881547000`)
var MessageStruct = Message{
	Name: "Alice",
	Body: "Hello",
	Time: 1294706395881547000,
}

// Message struct
type Message struct {
	Name string `yaml:"Name"`
	Body string `yaml:"Body"`
	Time int64  `yaml:"Time"`
}

func TestYamlConnectorConfig_Get(t *testing.T) {
	conn := base_configprovider.NewBufferedConnector("key", "scope", MessageBytes)
	c := NewYamlConnectorConfig("key", "scope", conn)

	var m Message
	res := c.Get(&m)
	<-res.Finished()

	if !res.Success() {
		t.Error("YamlConnectorConfig Config reported failure in Get() : ", res.Errors())
	} else {

		if m.Name != MessageStruct.Name {
			t.Error("YamlConnectorConfig provided incorrect data ==> Name : ", m.Name)
		}
		if m.Body != MessageStruct.Body {
			t.Error("YamlConnectorConfig provided incorrect data ==> Body : ", m.Body)
		}
		if m.Time != MessageStruct.Time {
			t.Error("YamlConnectorConfig provided incorrect data ==> Time : ", m.Time)
		}

	}

}

func TestYamlConnectorConfig_Set(t *testing.T) {
	conn := base_configprovider.NewBufferedConnector("key", "scope", []byte{})
	c := NewYamlConnectorConfig("key", "scope", conn)

	res := c.Set(MessageStruct)
	<-res.Finished()

	if !res.Success() {
		t.Error("YamlConnectorConfig Config reported failure in Set()", res.Errors())
	} else {
		var m Message
		res = c.Get(&m)
		<-res.Finished()

		if !res.Success() {
			t.Error("YamlConnectorConfig Config reported failure in Get() : ", res.Errors(), m)
		} else {

			if m.Name != MessageStruct.Name {
				t.Error("YamlConnectorConfig provided incorrect data ==> Name : ", m.Name)
			}
			if m.Body != MessageStruct.Body {
				t.Error("YamlConnectorConfig provided incorrect data ==> Body : ", m.Body)
			}
			if m.Time != MessageStruct.Time {
				t.Error("YamlConnectorConfig provided incorrect data ==> Time : ", m.Time)
			}

		}
	}
}
