package parser

type Record struct {
	ID      string
	Payload []byte
}

type Decoder struct{ buffer []byte }

func (d *Decoder) Parse(id, payload string) Record {
	d.buffer = append(d.buffer[:0], payload...)
	return Record{ID: id, Payload: d.buffer}
}

func Clone(record Record) Record { return record }
