package config

import "testing"

func TestHmlMetaTagGenerator_Validate(t *testing.T) {
	type fields struct {
		TagName string
		Code    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty tag name",
			fields: fields{
				TagName: "",
			},
			wantErr: true,
		},
		{
			name: "empty code",
			fields: fields{
				TagName: "test",
				Code:    "",
			},
			wantErr: true,
		},
		{
			name: "valid",
			fields: fields{
				TagName: "test",
				Code:    "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HmlMetaTagGenerator{
				TagName: tt.fields.TagName,
				Code:    tt.fields.Code,
			}
			if err := h.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJsonGenerator_Validate(t *testing.T) {
	type fields struct {
		FileName  string
		Attribute string
		Code      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty file name",
			fields: fields{
				FileName: "",
			},
			wantErr: true,
		},
		{
			name: "empty attribute",
			fields: fields{
				FileName:  "test.json",
				Attribute: "",
			},
			wantErr: true,
		},
		{
			name: "empty code",
			fields: fields{
				FileName:  "test.json",
				Attribute: "test",
				Code:      "",
			},
			wantErr: true,
		},
		{
			name: "valid",
			fields: fields{
				FileName:  "test.json",
				Attribute: "test",
				Code:      "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JsonGenerator{
				FileName:  tt.fields.FileName,
				Attribute: tt.fields.Attribute,
				Code:      tt.fields.Code,
			}
			if err := j.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestXmlGenerator_Validate(t *testing.T) {
	type fields struct {
		FileName string
		RootName string
		Code     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty file name",
			fields: fields{
				FileName: "",
			},
			wantErr: true,
		},
		{
			name: "empty root name",
			fields: fields{
				FileName: "test.xml",
				RootName: "",
			},
			wantErr: true,
		},
		{
			name: "empty code",
			fields: fields{
				FileName: "test.xml",
				RootName: "test",
				Code:     "",
			},
			wantErr: true,
		},
		{
			name: "valid",
			fields: fields{
				FileName: "test.xml",
				RootName: "test",
				Code:     "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &XmlGenerator{
				FileName: tt.fields.FileName,
				RootName: tt.fields.RootName,
				Code:     tt.fields.Code,
			}
			if err := x.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestXmlGenerator_ToXml(t *testing.T) {
	type fields struct {
		FileName string
		RootName string
		Code     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "valid",
			fields: fields{
				FileName: "test.xml",
				RootName: "test",
				Code:     "test",
			},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<test><code>test</code></test>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &XmlGenerator{
				FileName: tt.fields.FileName,
				RootName: tt.fields.RootName,
				Code:     tt.fields.Code,
			}
			if got := x.ToXml(); got != tt.want {
				t.Errorf("ToXml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTxtRecordGenerator_Validate(t1 *testing.T) {
	type fields struct {
		HostName             string
		RecordAttribute      string
		RecordAttributeValue string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty host name",
			fields: fields{
				HostName: "",
			},
			wantErr: true,
		},
		{
			name: "empty record attribute",
			fields: fields{
				HostName:        "test",
				RecordAttribute: "",
			},
			wantErr: true,
		},
		{
			name: "empty record attribute value",
			fields: fields{
				HostName:             "test",
				RecordAttribute:      "test",
				RecordAttributeValue: "",
			},
			wantErr: true,
		},
		{
			name: "valid",
			fields: fields{
				HostName:             "test",
				RecordAttribute:      "test",
				RecordAttributeValue: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TxtRecordGenerator{
				HostName:             tt.fields.HostName,
				RecordAttribute:      tt.fields.RecordAttribute,
				RecordAttributeValue: tt.fields.RecordAttributeValue,
			}
			if err := t.Validate(); (err != nil) != tt.wantErr {
				t1.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCnameRecordGenerator_Validate(t *testing.T) {
	type fields struct {
		RecordName   string
		RecordTarget string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty record name",
			fields: fields{
				RecordName: "",
			},
			wantErr: true,
		},
		{
			name: "empty record target",
			fields: fields{
				RecordName:   "test",
				RecordTarget: "",
			},
			wantErr: true,
		},
		{
			name: "valid",
			fields: fields{
				RecordName:   "test",
				RecordTarget: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CnameRecordGenerator{
				RecordName:   tt.fields.RecordName,
				RecordTarget: tt.fields.RecordTarget,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
