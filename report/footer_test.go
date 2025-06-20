package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/jung-kurt/gofpdf"
)

func TestFooterFuncLpi(t *testing.T) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	var (
		oldFooterFnc  = "oldFooterFnc"
		bothPages     = "bothPages"
		firstPageOnly = "firstPageOnly"
		lastPageOnly  = "lastPageOnly"
	)

	// This set just for testing, only set SetFooterFuncLpi.
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, oldFooterFnc,
			"", 0, "C", false, 0, "")
	})
	pdf.SetFooterFuncLpi(func(lastPage bool) {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, bothPages, "", 0, "L", false, 0, "")
		if !lastPage {
			pdf.CellFormat(0, 10, firstPageOnly, "", 0, "C", false, 0, "")
		} else {
			pdf.CellFormat(0, 10, lastPageOnly, "", 0, "C", false, 0, "")
		}
	})
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	for j := 1; j <= 40; j++ {
		pdf.CellFormat(0, 10, fmt.Sprintf("Printing line number %d", j),
			"", 1, "", false, 0, "")
	}
	if pdf.Error() != nil {
		t.Fatalf("not expecting error when rendering text")
	}
	w := &bytes.Buffer{}
	if err := pdf.Output(w); err != nil {
		t.Errorf("unexpected err: %s", err)
	}
	b := w.Bytes()
	if bytes.Contains(b, []byte(oldFooterFnc)) {
		t.Errorf("not expecting %s render on pdf when FooterFncLpi is set", oldFooterFnc)
	}
	got := bytes.Count(b, []byte("bothPages"))
	if got != 2 {
		t.Errorf("footer %s should render on two page got:%d", bothPages, got)
	}
	got = bytes.Count(b, []byte(firstPageOnly))
	if got != 1 {
		t.Errorf("footer %s should render only on first page got: %d", firstPageOnly, got)
	}
	got = bytes.Count(b, []byte(lastPageOnly))
	if got != 1 {
		t.Errorf("footer %s should render only on first page got: %d", lastPageOnly, got)
	}
	f := bytes.Index(b, []byte(firstPageOnly))
	l := bytes.Index(b, []byte(lastPageOnly))
	if f > l {
		t.Errorf("index %d (%s) should less than index %d (%s)", f, firstPageOnly, l, lastPageOnly)
	}
}
