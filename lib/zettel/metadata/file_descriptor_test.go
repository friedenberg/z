package metadata

import "testing"

type descriptorTestCase struct {
	description string
	fd          *FileDescriptor
	expectedTag string
}

func getDescriptorTestCases(t *testing.T) []descriptorTestCase {
	return []descriptorTestCase{
		descriptorTestCase{
			description: "local with extension no kasten",
			fd: &FileDescriptor{
				ZettelId: Id(1),
				Ext:      "test",
			},
			expectedTag: "1.test",
		},
		descriptorTestCase{
			description: "local without extension no kasten",
			fd: &FileDescriptor{
				ZettelId: Id(1),
			},
			expectedTag: "1",
		},
		descriptorTestCase{
			description: "local with extension with kasten",
			fd: &FileDescriptor{
				KastenName: "some_kasten",
				ZettelId:   Id(1),
				Ext:        "test",
			},
			expectedTag: "1.test-some_kasten",
		},
		descriptorTestCase{
			description: "local without extension with kasten",
			fd: &FileDescriptor{
				KastenName: "some_kasten",
				ZettelId:   Id(1),
			},
			expectedTag: "1-some_kasten",
		},
	}
}

func TestDescriptors(t *testing.T) {
	for _, tc := range getDescriptorTestCases(t) {
		t.Run(
			tc.description,
			func(t *testing.T) {
				actualTag := tc.fd.Tag()

				if actualTag != tc.expectedTag {
					t.Errorf("Actual tag was '%s', wanted '%s'", actualTag, tc.expectedTag)
				}

				actualFd := &FileDescriptor{}
				err := actualFd.Set(tc.expectedTag)

				if err != nil {
					t.Errorf("failed to set from tag: %w", err)
				}

				expectedFd := tc.fd

				if *actualFd != *expectedFd {
					t.Errorf("Actual fd was '%s', wanted '%s'", actualFd, expectedFd)
				}
			},
		)
	}
}
