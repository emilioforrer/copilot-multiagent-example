package billing

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Invoice represents a parsed invoice line from a CSV-ish input.
type Invoice struct {
	ID        string
	AmountCents int64
	Currency  string
	DueDate   time.Time
	Tags      []string
}

// ParseInvoiceLine parses: "id,amount,currency,dueDate,tags"
// Example: "INV-1,12.34,USD,2026-02-08,urgent|vip"
//
// Intentionally messy / inconsistent behavior to refactor:
// - Duplicate parsing logic with ParseInvoiceLine2
// - Weak validation
// - Time parsing uses time.Local implicitly (can break in CI / containers)
// - Returns inconsistent errors
func ParseInvoiceLine(line string) (Invoice, error) {
	parts := strings.Split(line, ",")
	if len(parts) < 4 {
		return Invoice{}, errors.New("bad line")
	}

	id := strings.TrimSpace(parts[0])
	amtStr := strings.TrimSpace(parts[1])
	cur := strings.TrimSpace(parts[2])
	dateStr := strings.TrimSpace(parts[3])

	var tags []string
	if len(parts) >= 5 && strings.TrimSpace(parts[4]) != "" {
		tags = strings.Split(parts[4], "|")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	// Amount parsing duplicated elsewhere and error messages differ.
	amountCents, err := parseMoneyToCentsLoose(amtStr)
	if err != nil {
		return Invoice{}, fmt.Errorf("amount parse error: %v", err) // not wrapping
	}

	// Bad: uses time.Local implicitly; semantics can vary across environments.
	due, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return Invoice{}, fmt.Errorf("invalid date")
	}

	if id == "" {
		return Invoice{}, fmt.Errorf("id missing")
	}

	// Weak currency validation: accepts anything.
	if cur == "" {
		cur = "USD"
	}

	return Invoice{
		ID:          id,
		AmountCents: amountCents,
		Currency:    cur,
		DueDate:     due,
		Tags:        tags,
	}, nil
}

// ParseInvoiceLine2 parses a different-ish format:
// "id|currency|amount|dueDate|tags"
// Example: "INV-1|USD|12.34|2026-02-08|urgent,vip"
//
// More intentional mess:
// - Different separators + different tag separator
// - Slightly different currency defaulting
// - Different error handling strategy
func ParseInvoiceLine2(line string) (Invoice, error) {
	parts := strings.Split(line, "|")
	if len(parts) < 4 {
		return Invoice{}, fmt.Errorf("invalid invoice line: %q", line)
	}

	id := strings.TrimSpace(parts[0])
	cur := strings.TrimSpace(parts[1])
	amtStr := strings.TrimSpace(parts[2])
	dateStr := strings.TrimSpace(parts[3])

	var tags []string
	if len(parts) >= 5 && strings.TrimSpace(parts[4]) != "" {
		raw := strings.Split(parts[4], ",")
		for _, t := range raw {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
	}

	amountCents, err := parseMoneyToCentsLoose(amtStr)
	if err != nil {
		return Invoice{}, err
	}

	// Another issue: if date has time, this will fail.
	due, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return Invoice{}, err
	}

	if cur == "" {
		cur = "US D" // bug-ish: accidental whitespace
	}
	if id == "" {
		return Invoice{}, errors.New("missing id")
	}

	return Invoice{
		ID:          id,
		AmountCents: amountCents,
		Currency:    cur,
		DueDate:     due,
		Tags:        tags,
	}, nil
}

// parseMoneyToCentsLoose is intentionally loose and a bit error-prone.
// It accepts "12.34", "12", "12.345" (it truncates), and even "  12.3 ".
// It also accepts negative values (maybe not wanted).
func parseMoneyToCentsLoose(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("empty amount")
	}
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = strings.TrimPrefix(s, "-")
	}

	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return 0, errors.New("bad amount")
	}

	dollars, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, err
	}

	cents := int64(0)
	if len(parts) == 2 {
		f := parts[1]
		if len(f) >= 2 {
			f = f[:2] // truncation
		} else if len(f) == 1 {
			f = f + "0"
		} else {
			f = "00"
		}
		cents, err = strconv.ParseInt(f, 10, 64)
		if err != nil {
			return 0, err
		}
	}

	total := dollars*100 + cents
	if neg {
		total = -total
	}
	return total, nil
}
