package generator

import (
	"testing"

	"github.com/muxinc/protogen/proto3"
)

func TestToProtobufSpec(t *testing.T) {
	type args struct {
		spec *proto3.Spec
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Nil Spec",
			args:    args{nil},
			want:    "",
			wantErr: true,
		},
		{
			name: "Spec with only a Package statement",
			args: args{
				&proto3.Spec{
					Package: "foo",
				},
			},
			want:    "syntax = \"proto3\";\npackage foo;\n",
			wantErr: false,
		},
		{
			name: "Spec with only Package and Imports",
			args: args{
				&proto3.Spec{
					Package: "foo",
					Imports: []proto3.ImportType{"package1"},
				},
			},
			want:    "syntax = \"proto3\";\npackage foo;\nimport \"package1\";\n",
			wantErr: false,
		},
		{
			name: "Spec with only Package and Multiple Imports",
			args: args{
				&proto3.Spec{
					Package: "foo",
					Imports: []proto3.ImportType{"package1", "package2"},
				},
			},
			want:    "syntax = \"proto3\";\npackage foo;\nimport \"package1\";\nimport \"package2\";\n",
			wantErr: false,
		},
		{
			name: "Spec with message using reserved tags",
			args: args{
				&proto3.Spec{
					Package: "foo",
					Messages: []proto3.Message{
						proto3.Message{
							Name: "Beacon",
							Messages: []proto3.Message{
								proto3.Message{
									Name: "Event",
									ReservedValues: []proto3.Reserved{
										proto3.ReservedTagValue{Tag: 1},
										proto3.ReservedTagValue{Tag: 2},
										proto3.ReservedTagValue{Tag: 3},
										proto3.ReservedTagRange{LowerTag: 6, UpperTag: 9},
									},
									Fields: []proto3.Field{
										proto3.CustomField{Name: "Habitat", Typing: "string", Tag: 10, Comment: "What am I?"},
										proto3.ScalarField{Name: "Continent", Typing: proto3.STRING_TYPE, Tag: 11, Comment: "Where am I?"},
										proto3.MapField{Name: "LanguageMap", KeyTyping: proto3.STRING_TYPE, ValueTyping: proto3.STRING_TYPE, Tag: 12, Comment: "Super essential"},
									},
								},
							},
							ReservedValues: []proto3.Reserved{
								proto3.ReservedTagValue{Tag: 1},
								proto3.ReservedTagValue{Tag: 2},
								proto3.ReservedTagValue{Tag: 3},
								proto3.ReservedTagRange{LowerTag: 6, UpperTag: 9},
							},
							Fields: []proto3.Field{
								proto3.CustomField{Name: "Habitat", Typing: "string", Tag: 20, Comment: "What am I?"},
								proto3.ScalarField{Name: "Continent", Typing: proto3.STRING_TYPE, Tag: 21, Comment: "Where am I?"},
								proto3.MapField{Name: "LanguageMap", KeyTyping: proto3.STRING_TYPE, ValueTyping: proto3.STRING_TYPE, Tag: 22, Comment: "Super essential"},
								proto3.CustomMapField{Name: "CustomMap", KeyTyping: proto3.STRING_TYPE, ValueTyping: "Event", Tag: 23},
							},
							Enums: []proto3.Enum{
								proto3.Enum{
									Name: "Country",
									Values: []proto3.EnumValue{
										proto3.EnumValue{Name: "US", Tag: 0},
										proto3.EnumValue{Name: "CA", Tag: 1, Comment: "Canada"},
										proto3.EnumValue{Name: "GB", Tag: 2, Comment: "Great Britain"},
										proto3.EnumValue{Name: "MX", Tag: 3, Comment: "Mexico"},
									},
								},
								proto3.Enum{
									Name:       "PlaybackState",
									AllowAlias: true,
									Values: []proto3.EnumValue{
										proto3.EnumValue{Name: "Waiting", Tag: 0},
										proto3.EnumValue{Name: "Playing", Tag: 1},
										proto3.EnumValue{Name: "Started", Tag: 1},
										proto3.EnumValue{Name: "Stopped", Tag: 2},
									},
								},
							},
						},
					},
				},
			},
			want:    "syntax = \"proto3\";\npackage foo;\nimport \"package1\";\nimport \"package2\";\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := ToProtobufSpec(tt.args.spec)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. ToProtobufSpec() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. ToProtobufSpec() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
