package commands

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestValidateLeaseOptions(t *testing.T) {
	gs := NewWithT(t)

	cases := []struct {
		name              string
		options           []string
		expectedPoolValue string
	}{
		{
			name: "Normal pool name",
			options: []string{
				"cpus=4",
				"memory=16",
				"networks=1",
				"pools=pool1",
			},
			expectedPoolValue: "pool1",
		},
		{
			name: "Pool name with quotes around it",
			options: []string{
				"cpus=4",
				"memory=16",
				"networks=1",
				"pools=\"pool1\"",
			},
			expectedPoolValue: "pool1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			options := getLeaseOptions(tc.options)

			gs.Expect(options.pool).To(Equal(tc.expectedPoolValue))
		})
	}
}
