package checkpointcommit

type Handle interface {
	Write([]byte) error
	Close() error
}

type OpenFunc func(partition int) (Handle, error)

type Writer struct {
	Open OpenFunc
}

func (w Writer) WriteAll(payloads [][]byte) error {
	for partition, payload := range payloads {
		handle, err := w.Open(partition)
		if err != nil {
			return err
		}
		defer handle.Close()
		if err := handle.Write(payload); err != nil {
			return err
		}
	}
	return nil
}
