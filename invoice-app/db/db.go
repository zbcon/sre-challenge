package db

type Client struct {
	invoices []*Invoice
}

type Invoice struct {
	InvoiceId string
	Value     float32
	Currency  string
	IsPaid    bool
}

func InitializeDatabase() *Client {
	invoices := []*Invoice{
		{
			InvoiceId: "I1",
			Value:     12.15,
			Currency:  "EUR",
			IsPaid: false,
		},
		{
			InvoiceId: "I2",
			Value:     10.25,
			Currency:  "GBP",
			IsPaid: false,
		},
		{
			InvoiceId: "I3",
			Value:     66.13,
			Currency:  "DKK",
			IsPaid: false,
		},
	}
	return &Client{
		invoices: invoices,
	}
}

func (c *Client) GetInvoices() []*Invoice {
	return c.invoices
}

func (c *Client) GetUnpaidInvoices() []*Invoice {
	var result []*Invoice
	for _, invoice := range c.invoices {
		if !invoice.IsPaid {
			result = append(result, invoice)
		}
	}
	return result
}

func (c *Client) PayInvoice(invoiceId string) {
	for _, invoice := range c.invoices {
		if invoice.InvoiceId == invoiceId {
			invoice.IsPaid = true
		}
	}
}
