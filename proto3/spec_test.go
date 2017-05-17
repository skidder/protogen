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

func TestSpec_Write(t *testing.T) {
	type fields struct {
		Package  string
		Imports  []ImportType
		Messages []Message
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "Invalid spec with zero messages",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "Spec with message using reserved tags",
			fields: fields{
				Package: "foo",
				Messages: []Message{
					Message{
						Name: "Beacon",
						Messages: []Message{
							Message{
								Name: "Event",
								ReservedValues: []Reserved{
									ReservedTagValue{Tag: 1},
									ReservedTagValue{Tag: 2},
									ReservedTagValue{Tag: 3},
									ReservedTagRange{LowerTag: 6, UpperTag: 9},
								},
								Fields: []Field{
									CustomField{Name: "Habitat", Typing: "string", Tag: 10, Rule: REPEATED, Comment: "What am I?"},
									ScalarField{Name: "Continent", Typing: STRING_TYPE, Tag: 11, Comment: "Where am I?"},
									MapField{Name: "LanguageMap", KeyTyping: STRING_TYPE, ValueTyping: STRING_TYPE, Tag: 12, Comment: "Super essential"},
								},
							},
						},
						ReservedValues: []Reserved{
							ReservedTagValue{Tag: 1},
							ReservedTagValue{Tag: 2},
							ReservedTagValue{Tag: 3},
							ReservedTagRange{LowerTag: 6, UpperTag: 9},
						},
						Fields: []Field{
							CustomField{Name: "Habitat", Typing: "string", Tag: 20, Rule: REQUIRED, Comment: "What am I?"},
							ScalarField{Name: "Continent", Typing: STRING_TYPE, Tag: 21, Rule: OPTIONAL, Comment: "Where am I?"},
							MapField{Name: "LanguageMap", KeyTyping: STRING_TYPE, ValueTyping: STRING_TYPE, Tag: 22, Comment: "Super essential"},
							CustomMapField{Name: "CustomMap", KeyTyping: STRING_TYPE, ValueTyping: "Event", Tag: 23},
						},
						Enums: []Enum{
							Enum{
								Name: "Country",
								Values: []EnumValue{
									EnumValue{Name: "US", Tag: 0},
									EnumValue{Name: "CA", Tag: 1, Comment: "Canada"},
									EnumValue{Name: "GB", Tag: 2, Comment: "Great Britain"},
									EnumValue{Name: "MX", Tag: 3, Comment: "Mexico"},
								},
							},
							Enum{
								Name:       "PlaybackState",
								AllowAlias: true,
								Values: []EnumValue{
									EnumValue{Name: "Waiting", Tag: 0},
									EnumValue{Name: "Playing", Tag: 1},
									EnumValue{Name: "Started", Tag: 1},
									EnumValue{Name: "Stopped", Tag: 2},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s := &Spec{
			Package:  tt.fields.Package,
			Imports:  tt.fields.Imports,
			Messages: tt.fields.Messages,
		}
		got, err := s.Write()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Spec.Write() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		t.Logf("Generated protobuf spec:\n%s\n", got)
		// if got != tt.want {
		// 	t.Errorf("%q. Spec.Write() = %v, want %v", tt.name, got, tt.want)
		// }
	}
}
