package options

import (
	"fmt"
	"github.com/spf13/pflag"
)

type TLSOptions struct {
	UseTLS             bool   `mapstructure:"use-tls" json:"use-tls"`
	InsecureSkipVerify bool   `mapstructure:"insecure-skip-verify" json:"insecure-skip-verify"`
	CaCert             string `mapstructure:"ca-cert" json:"ca-cert"`
	Cert               string `mapstructure:"cert" json:"cert"`
	Key                string `mapstructure:"key" json:"key"`
}

var _ IOptions = (*TLSOptions)(nil)

func NewTLSOptions() *TLSOptions {
	return &TLSOptions{}
}

func (o *TLSOptions) Validate() []error {
	var errs []error

	if !o.UseTLS {
		return errs
	}

	if (o.Cert != "" && o.Key == "") || (o.CaCert == "" && o.Key != "") {
		errs = append(errs, fmt.Errorf("only one of cert and key configure option is setted, you should set both to enable tls"))
	}
	return errs
}

func (o *TLSOptions) AddFlags(fs *pflag.FlagSet, prefixes ...string) {

}

//func (o *TLSOptions) MustTLSConfig() *tls.Config {
//	if !o.UseTLS {
//		return nil
//	}
//
//	tlsConfig := &tls.Config{
//		InsecureSkipVerify: o.InsecureSkipVerify,
//	}
//
//}
//
//func (o *TLSOptions) TLSConfig() (*tls.Config, error) {
//
//}

func (o *TLSOptions) Scheme() string {
	if o.UseTLS {
		return "https"
	}
	return "http"
}
