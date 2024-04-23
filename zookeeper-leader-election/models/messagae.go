package models

type Message struct {
	msg    string
	person []Person
}

func NewMessage(msg string, person []Person) *Message {
	return &Message{msg: msg, person: person}
}

func (m *Message) GetMsg() string {
	return m.msg
}

func (m *Message) GetPerson() []Person {
	return m.person
}

func (m *Message) SetMsg(msg string) {
	m.msg = msg
}

func (m *Message) SetPerson(person []Person) {
	m.person = person
}
