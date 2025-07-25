// SPDX-FileCopyrightText: The go-mail Authors
//
// SPDX-License-Identifier: MIT

package mail

import (
	"testing"
)

var (
	genHeaderTests = []struct {
		name   string
		header Header
		want   string
	}{
		{"Header: Content-Description", HeaderContentDescription, "Content-Description"},
		{"Header: Content-Disposition", HeaderContentDisposition, "Content-Disposition"},
		{"Header: Content-ID", HeaderContentID, "Content-ID"},
		{"Header: Content-Language", HeaderContentLang, "Content-Language"},
		{"Header: Content-Location", HeaderContentLocation, "Content-Location"},
		{"Header: Content-Transfer-Encoding", HeaderContentTransferEnc, "Content-Transfer-Encoding"},
		{"Header: Content-Type", HeaderContentType, "Content-Type"},
		{"Header: Date", HeaderDate, "Date"},
		{
			"Header: Disposition-Notification-To", HeaderDispositionNotificationTo,
			"Disposition-Notification-To",
		},
		{"Header: Importance", HeaderImportance, "Importance"},
		{"Header: In-Reply-To", HeaderInReplyTo, "In-Reply-To"},
		{"Header: List-Unsubscribe", HeaderListUnsubscribe, "List-Unsubscribe"},
		{"Header: List-Unsubscribe-Post", HeaderListUnsubscribePost, "List-Unsubscribe-Post"},
		{"Header: Message-ID", HeaderMessageID, "Message-ID"},
		{"Header: MIME-Version", HeaderMIMEVersion, "MIME-Version"},
		{"Header: Organization", HeaderOrganization, "Organization"},
		{"Header: Precedence", HeaderPrecedence, "Precedence"},
		{"Header: Priority", HeaderPriority, "Priority"},
		{"Header: References", HeaderReferences, "References"},
		{"Header: Subject", HeaderSubject, "Subject"},
		{"Header: User-Agent", HeaderUserAgent, "User-Agent"},
		{"Header: X-Auto-Response-Suppress", HeaderXAutoResponseSuppress, "X-Auto-Response-Suppress"},
		{"Header: X-Mailer", HeaderXMailer, "X-Mailer"},
		{"Header: X-MSMail-Priority", HeaderXMSMailPriority, "X-MSMail-Priority"},
		{"Header: X-Priority", HeaderXPriority, "X-Priority"},
	}
	addrHeaderTests = []struct {
		name   string
		header AddrHeader
		want   string
	}{
		{"EnvelopeFrom", HeaderEnvelopeFrom, "EnvelopeFrom"},
		{"From", HeaderFrom, "From"},
		{"To", HeaderTo, "To"},
		{"Cc", HeaderCc, "Cc"},
		{"Bcc", HeaderBcc, "Bcc"},
		{"Reply-To", HeaderReplyTo, "Reply-To"},
	}
)

func TestImportance_Stringer(t *testing.T) {
	tests := []struct {
		name    string
		imp     Importance
		wantnum string
		xprio   string
		want    string
	}{
		{"Non-Urgent", ImportanceNonUrgent, "0", "5", "non-urgent"},
		{"Low", ImportanceLow, "0", "5", "low"},
		{"Normal", ImportanceNormal, "", "", ""},
		{"High", ImportanceHigh, "1", "1", "high"},
		{"Urgent", ImportanceUrgent, "1", "1", "urgent"},
		{"Unknown", 9, "", "", ""},
	}
	t.Run("String", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.imp.String() != tt.want {
					t.Errorf("wrong string for Importance returned. Expected: %s, got: %s", tt.want, tt.imp.String())
				}
			})
		}
	})
	t.Run("NumString", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.imp.NumString() != tt.wantnum {
					t.Errorf("wrong number string for Importance returned. Expected: %s, got: %s", tt.wantnum,
						tt.imp.NumString())
				}
			})
		}
	})
	t.Run("XPrioString", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.imp.XPrioString() != tt.xprio {
					t.Errorf("wrong x-prio string for Importance returned. Expected: %s, got: %s", tt.xprio,
						tt.imp.XPrioString())
				}
			})
		}
	})
}

func TestAddrHeader_Stringer(t *testing.T) {
	for _, tt := range addrHeaderTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.header.String() != tt.want {
				t.Errorf("wrong string for AddrHeader returned. Expected: %s, got: %s",
					tt.want, tt.header.String())
			}
		})
	}
}

func TestHeader_Stringer(t *testing.T) {
	for _, tt := range genHeaderTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.header.String() != tt.want {
				t.Errorf("wrong string for Header returned. Expected: %s, got: %s",
					tt.want, tt.header.String())
			}
		})
	}
}

func TestIsAddrHeader(t *testing.T) {
	t.Run("validate all address headers", func(t *testing.T) {
		for _, tt := range addrHeaderTests {
			t.Run(tt.name+" is a valid address header", func(t *testing.T) {
				if !IsAddrHeader(tt.header.String()) {
					t.Errorf("expected %s to be a valid address header", tt.header.String())
				}
			})
		}
	})
	t.Run("validate all non-address headers", func(t *testing.T) {
		for _, tt := range genHeaderTests {
			t.Run(tt.name+" is not an address header", func(t *testing.T) {
				if IsAddrHeader(tt.header.String()) {
					t.Errorf("expected %s to not be an address header", tt.header.String())
				}
			})
		}
	})
	t.Run("validate mixed headers from map", func(t *testing.T) {
		someHeaders := map[string]string{
			"To":         "toni.tester@example.com",
			"From":       "tina.tester@example.com",
			"Subject":    "this is the subject",
			"Message-ID": "test.mail.1234567@localhost.local",
		}
		message := NewMsg()
		for k, v := range someHeaders {
			if IsAddrHeader(k) {
				if err := message.SetAddrHeader(AddrHeader(k), v); err != nil {
					t.Errorf("failed to set address header: %s", err)
				}
				continue
			}
			message.SetGenHeader(Header(k), v)
		}
		checkAddrHeader(t, message, HeaderTo, "IsAddrHeader", 0, 1, "toni.tester@example.com", "")
		checkAddrHeader(t, message, HeaderFrom, "IsAddrHeader", 0, 1, "tina.tester@example.com", "")
		checkGenHeader(t, message, HeaderSubject, "IsAddrHeader", 0, 1, "this is the subject")
		checkGenHeader(t, message, HeaderMessageID, "IsAddrHeader", 0, 1, "test.mail.1234567@localhost.local")
	})
}
