package billing

import "testing"

func TestParseInvoiceLine_CommaFormat(t *testing.T) {
	inv, err := ParseInvoiceLine("INV-1,12.34,USD,2026-02-08,urgent|vip")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inv.ID != "INV-1" {
		t.Fatalf("id: got %q", inv.ID)
	}
	if inv.AmountCents != 1234 {
		t.Fatalf("amount: got %d", inv.AmountCents)
	}
	if inv.Currency != "USD" {
		t.Fatalf("currency: got %q", inv.Currency)
	}
	if len(inv.Tags) != 2 || inv.Tags[0] != "urgent" || inv.Tags[1] != "vip" {
		t.Fatalf("tags: %#v", inv.Tags)
	}
}

func TestParseInvoiceLine_PartialCents_Truncates(t *testing.T) {
	inv, err := ParseInvoiceLine("INV-2,12.345,USD,2026-02-08,")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// current behavior truncates to 12.34
	if inv.AmountCents != 1234 {
		t.Fatalf("amount: got %d", inv.AmountCents)
	}
}

func TestParseInvoiceLine2_PipeFormat(t *testing.T) {
	inv, err := ParseInvoiceLine2("INV-3|USD|12.34|2026-02-08|urgent,vip")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inv.AmountCents != 1234 {
		t.Fatalf("amount: got %d", inv.AmountCents)
	}
	// This test documents current behavior and will catch accidental changes.
	// Note: currency default bug exists but this case uses USD.
	if inv.Currency != "USD" {
		t.Fatalf("currency: got %q", inv.Currency)
	}
	if len(inv.Tags) != 2 {
		t.Fatalf("tags: %#v", inv.Tags)
	}
}

func TestParseInvoiceLine2_DefaultCurrency_BugDocumented(t *testing.T) {
	inv, err := ParseInvoiceLine2("INV-4||12.34|2026-02-08|")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Current buggy behavior: "US D" with a space.
	if inv.Currency != "US D" {
		t.Fatalf("currency: got %q", inv.Currency)
	}
}
