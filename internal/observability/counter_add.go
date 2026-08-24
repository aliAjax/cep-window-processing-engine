package observability

func (c *Counter) Add(n uint64) {
	current := c.Value
	c.Value = current + n
}
