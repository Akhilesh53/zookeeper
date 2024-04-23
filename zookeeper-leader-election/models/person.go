package models

type Person struct {
	id   uint16
	name string
}

func NewPerson(id uint16, name string) *Person {
	return &Person{id: id, name: name}
}

func (p *Person) GetId() uint16 {
	return p.id
}

func (p *Person) GetName() string {
	return p.name
}

func (p *Person) SetId(id uint16) {
	p.id = id
}

func (p *Person) SetName(name string) {
	p.name = name
}
