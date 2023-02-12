package domainverifier

import "testing"

func TestCheckHtmlMetaTag(t *testing.T) {
	type args struct {
		domain         string
		metaTagName    string
		metaTagContent string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Successful html meta verification",
			args: args{
				domain:         "fr.lioncoding.com",
				metaTagName:    "google-site-verification",
				metaTagContent: "sfRybH_Mn50-a_lGoRf21hf28qx1iOucU8CsBe_hEVM",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Failed html meta verification",
			args: args{
				domain:         "fr.lioncoding.com",
				metaTagName:    "myapp-site-verification",
				metaTagContent: "1234567891",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Invalid domain",
			args: args{
				domain:         "invalid domain",
				metaTagName:    "myapp-site-verification",
				metaTagContent: "1234567891",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckHtmlMetaTag(tt.args.domain, tt.args.metaTagName, tt.args.metaTagContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckHtmlMetaTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckHtmlMetaTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type ownershipVerification struct {
	Code string `xml:"code" json:"myapp_site_verification"`
}

func TestCheckJsonFile(t *testing.T) {
	type args struct {
		domain        string
		fileName      string
		expectedValue interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Successful json verification",
			args: args{
				domain:        "domainverify.lioncoding.workers.dev",
				fileName:      "myapp-site-verification.json",
				expectedValue: ownershipVerification{Code: "dcf56hgvghy674fc"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Failed json verification",
			args: args{
				domain:        "domainverify.lioncoding.workers.dev",
				fileName:      "myapp-site-verification.json",
				expectedValue: ownershipVerification{Code: "1234567891"},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Invalid domain",
			args: args{
				domain:        "invalid domain",
				fileName:      "myapp-site-verification.json",
				expectedValue: ownershipVerification{Code: "1234567891"},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckJsonFile(tt.args.domain, tt.args.fileName, tt.args.expectedValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckJsonFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckJsonFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckXmlFile(t *testing.T) {
	type args struct {
		domain        string
		fileName      string
		expectedValue interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Successful xml verification",
			args: args{
				domain:        "domainverify.lioncoding.workers.dev",
				fileName:      "myappSiteAuth.xml",
				expectedValue: ownershipVerification{Code: "dcf56hgvghy674fc"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Failed xml verification",
			args: args{
				domain:        "domainverify.lioncoding.workers.dev",
				fileName:      "myappSiteAuth.xml",
				expectedValue: ownershipVerification{Code: "1234567891"},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Invalid domain",
			args: args{
				domain:        "invalid domain",
				fileName:      "myappSiteAuth.xml",
				expectedValue: ownershipVerification{Code: "1234567891"},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckXmlFile(tt.args.domain, tt.args.fileName, tt.args.expectedValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckXmlFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckXmlFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
