// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package azuremonitorreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azuremonitorreceiver"

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/assert"
)

func Test_newMetricsQueryFilterFromDimensions(t *testing.T) {
	assert.Nil(t, newMetricsQueryFilterFromDimensions(nil))
	assert.Nil(t, newMetricsQueryFilterFromDimensions([]string{}))
	assert.Equal(t, to.Ptr("hello eq '*'"), newMetricsQueryFilterFromDimensions([]string{"hello"}))
	assert.Equal(t, to.Ptr("hello eq '*' and hi eq '*'"), newMetricsQueryFilterFromDimensions([]string{"hello", "hi"}))
}
