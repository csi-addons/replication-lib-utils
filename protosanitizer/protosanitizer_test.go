/*
Copyright 2021 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package protosanitizer

import (
	"fmt"
	"testing"

	"github.com/csi-addons/spec/lib/go/replication"
	"github.com/stretchr/testify/assert"
)

func TestStripSecrets(t *testing.T) {

	var (
		stripped    fmt.Stringer
		secretName  = "secret-abc"
		secretValue = "123"
	)
	// Enable Volume replication.
	enableVolumeReplication := &replication.EnableVolumeReplicationRequest{
		VolumeId: "foo",
		Secrets: map[string]string{
			secretName:   secretValue,
			"secret-xyz": "987",
		},
		Parameters: map[string]string{
			"mirroringMode": "snapshot",
		},
	}

	// Disable Volume replication.
	disableVolumeReplication := &replication.DisableVolumeReplicationRequest{
		VolumeId: "foo",
		Secrets: map[string]string{
			secretName:   secretValue,
			"secret-xyz": "987",
		},
		Parameters: map[string]string{
			"force": "false",
		},
	}

	// Promote volume.
	promoteVolume := &replication.PromoteVolumeRequest{
		VolumeId: "foo",
		Secrets: map[string]string{
			secretName:   secretValue,
			"secret-xyz": "987",
		},
		Parameters: map[string]string{
			"force": "false",
		},
	}

	// Demote volume.
	demoteVolume := &replication.DemoteVolumeRequest{
		VolumeId: "foo",
		Secrets: map[string]string{
			secretName:   secretValue,
			"secret-xyz": "987",
		},
		Parameters: map[string]string{
			"force": "false",
		},
	}
	// Demote volume.
	resyncVolume := &replication.ResyncVolumeRequest{
		VolumeId: "foo",
		Secrets: map[string]string{
			secretName:   secretValue,
			"secret-xyz": "987",
		},
		Parameters: map[string]string{
			"force": "false",
		},
	}
	type testcase struct {
		original, stripped interface{}
	}

	cases := []testcase{
		{nil, "null"},
		{1, "1"},
		{"hello world", `"hello world"`},
		{true, "true"},
		{false, "false"},
		{&replication.EnableVolumeReplicationRequest{}, `{}`},
		{enableVolumeReplication, `{"parameters":{"mirroringMode":"snapshot"},"secrets":"***stripped***","volume_id":"foo"}`},
		{disableVolumeReplication, `{"parameters":{"force":"false"},"secrets":"***stripped***","volume_id":"foo"}`},
		{promoteVolume, `{"parameters":{"force":"false"},"secrets":"***stripped***","volume_id":"foo"}`},
		{demoteVolume, `{"parameters":{"force":"false"},"secrets":"***stripped***","volume_id":"foo"}`},
		{resyncVolume, `{"parameters":{"force":"false"},"secrets":"***stripped***","volume_id":"foo"}`},
	}

	for _, c := range cases {
		before := fmt.Sprint(c.original)
		stripped = StripReplicationSecrets(c.original)
		if assert.Equal(t, c.stripped, fmt.Sprintf("%s", stripped), "unexpected result for fmt s of %s", c.original) {
			if assert.Equal(t, c.stripped, fmt.Sprintf("%v", stripped), "unexpected result for fmt v of %s", c.original) {
				assert.Equal(t, c.stripped, fmt.Sprintf("%+v", stripped), "unexpected result for fmt +v of %s", c.original)
			}
		}
		assert.Equal(t, before, fmt.Sprint(c.original), "original value modified")
	}

	// The secret is hidden because StripSecrets is a struct referencing it.
	dump := fmt.Sprintf("%#v", StripReplicationSecrets(enableVolumeReplication))
	assert.NotContains(t, dump, secretName)
	assert.NotContains(t, dump, secretValue)

	dump = fmt.Sprintf("%#v", StripReplicationSecrets(disableVolumeReplication))
	assert.NotContains(t, dump, secretName)
	assert.NotContains(t, dump, secretValue)

	dump = fmt.Sprintf("%#v", StripReplicationSecrets(promoteVolume))
	assert.NotContains(t, dump, secretName)
	assert.NotContains(t, dump, secretValue)

	dump = fmt.Sprintf("%#v", StripReplicationSecrets(demoteVolume))
	assert.NotContains(t, dump, secretName)
	assert.NotContains(t, dump, secretValue)

	dump = fmt.Sprintf("%#v", StripReplicationSecrets(resyncVolume))
	assert.NotContains(t, dump, secretName)
	assert.NotContains(t, dump, secretValue)
}
