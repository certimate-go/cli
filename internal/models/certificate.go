package models

import "time"

// CertificateSourceType defines where the certificate came from
type CertificateSourceType string

const (
	CertificateSourceRequest CertificateSourceType = "request"
	CertificateSourceUpload  CertificateSourceType = "upload"
)

// Certificate represents an SSL certificate
type Certificate struct {
	ID              string                `json:"id"`
	WorkflowID      string                `json:"workflowId"`
	Source          CertificateSourceType `json:"source"`
	SubjectAltNames StringArray           `json:"subjectAltNames"`
	SerialNumber    string                `json:"serialNumber"`
	Certificate     string                `json:"certificate"`
	PrivateKey      string                `json:"privateKey"`
	Issuer          string                `json:"issuer"`
	NotBefore       CustomTime            `json:"notBefore"`
	NotAfter        CustomTime            `json:"notAfter"`
	Created         CustomTime            `json:"created"`
	Updated         CustomTime            `json:"updated"`
}

// IsExpired returns true if the certificate has expired
func (c *Certificate) IsExpired() bool {
	if c.NotAfter.IsZero() {
		return false
	}
	return time.Now().After(c.NotAfter.Time)
}

// ExpiresWithin returns true if the certificate expires within the given duration
func (c *Certificate) ExpiresWithin(d time.Duration) bool {
	if c.NotAfter.IsZero() {
		return false
	}
	return time.Now().Add(d).After(c.NotAfter.Time)
}

// DaysUntilExpiry returns the number of days until the certificate expires
func (c *Certificate) DaysUntilExpiry() int {
	if c.NotAfter.IsZero() {
		return -1
	}
	return int(time.Until(c.NotAfter.Time).Hours() / 24)
}
