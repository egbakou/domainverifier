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
