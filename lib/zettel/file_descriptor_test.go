package zettel

import "testing"

type descriptorTestCase struct {
	description string
	noteTag     Tag
	expectedTag string
}

func getDescriptorTestCases(t *testing.T) []descriptorTestCase {
	return []descriptorTestCase{
		descriptorTestCase{
			description: "local with extension no kasten",
			noteTag: &FileDescriptor{
				ZettelId: Id(1),
				Ext:      "test",
			},
			expectedTag: "1.test",
		},
		descriptorTestCase{
			description: "local without extension no kasten",
			noteTag: &FileDescriptor{
				ZettelId: Id(1),
			},
			expectedTag: "1",
		},
		descriptorTestCase{
			description: "local with extension with kasten",
			noteTag: &FileDescriptor{
				KastenName: "some_kasten",
				ZettelId:   Id(1),
				Ext:        "test",
			},
			expectedTag: "1.test-some_kasten",
		},
		descriptorTestCase{
			description: "local without extension with kasten",
			noteTag: &FileDescriptor{
				KastenName: "some_kasten",
				ZettelId:   Id(1),
			},
			expectedTag: "1-some_kasten",
		},
		descriptorTestCase{
			description: "remote with extension",
			noteTag: &RemoteFileDescriptor{
				FileDescriptor: FileDescriptor{
					KastenName: "some_kasten",
					ZettelId:   Id(1),
					Ext:        "test",
				},
				Version: FileVersion(1),
			},
			expectedTag: "1.test-1-some_kasten",
		},
	}
}

func TestDescriptors(t *testing.T) {
	for _, tc := range getDescriptorTestCases(t) {
		t.Run(
			tc.description,
			func(t *testing.T) {
				actual := tc.noteTag.Tag()

				if actual != tc.expectedTag {
					t.Errorf("Actual tag was '%s', wanted '%s'", tc.expectedTag, actual)
				}
			},
		)
	}
}
