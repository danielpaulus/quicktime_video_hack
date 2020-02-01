package gstadapter_test

import (
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/gstadapter"
	"github.com/stretchr/testify/assert"
)

func TestCustomPipelineParsing(t *testing.T) {

	_, err := gstadapter.NewWithCustomPipeline("daniel")
	assert.Error(t, err)

	_, err = gstadapter.NewWithCustomPipeline("queue name=my_filesrc ! fakesink")
	assert.Error(t, err)

	_, err = gstadapter.NewWithCustomPipeline("queue name=audio_target ! fakesink")
	assert.Error(t, err)

	gsta, err := gstadapter.NewWithCustomPipeline("rtpmux name=mux ! fakesink \n queue name=audio_target ! mux.sink_0 \n queue name=video_target ! mux.sink_1")
	assert.NoError(t, err)
	assert.NotNil(t, gsta)
}
