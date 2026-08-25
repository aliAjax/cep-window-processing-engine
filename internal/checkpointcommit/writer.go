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
		if err := w.writeOne(partition, payload); err != nil {
			return err
		}
	}
	return nil
}

func (w Writer) writeOne(partition int, payload []byte) (err error) {
	handle, err := w.Open(partition)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := handle.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	return handle.Write(payload)
}
