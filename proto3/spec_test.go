package proto3

import "testing"

func TestScalarField_Validate(t *testing.T) {
	type fields struct {
		Name    NameType
		Tag     TagType
		Rule    FieldRule
		Comment string
		Typing  FieldType
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Valid Scalar field",
			fields:  fields{Name: "MyMap", Tag: 1, Typing: STRING_TYPE},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s := &ScalarField{
			Name:    tt.fields.Name,
			Tag:     tt.fields.Tag,
			Rule:    tt.fields.Rule,
			Comment: tt.fields.Comment,
			Typing:  tt.fields.Typing,
		}
		if err := s.Validate(); (err != nil) != tt.wantErr {
			t.Errorf("%q. ScalarField.Validate() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
