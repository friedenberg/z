package metadata

import "testing"

type descriptorTestCase struct {
	description string
	fd          *LocalFile
	expectedTag string
}

func getDescriptorTestCases(t *testing.T) []descriptorTestCase {
	return []descriptorTestCase{
		descriptorTestCase{
			description: "local with extension no kasten",
			fd: &LocalFile{
				Id:  "1",
				Ext: "test",
			},
			expectedTag: "f-1.test",
		},
		descriptorTestCase{
			description: "local without extension no kasten",
			fd: &LocalFile{
				Id: "1",
			},
			expectedTag: "f-1",
		},
		descriptorTestCase{
			description: "local with extension with kasten",
			fd: &LocalFile{
				Id:  "1",
				Ext: "test",
			},
			expectedTag: "f-1.test",
		},
		descriptorTestCase{
			description: "local without extension with kasten",
			fd: &LocalFile{
				Id: "1",
			},
			expectedTag: "f-1",
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

				actualFd := &LocalFile{}
				err := actualFd.Set(tc.expectedTag)

				if err != nil {
					t.Errorf("failed to set from tag: %w", err)
				}

				expectedFd := tc.fd

				if *actualFd != *expectedFd {
					t.Errorf("Actual fd was '%q', wanted '%q'", actualFd, expectedFd)
				}
			},
		)
	}
}
