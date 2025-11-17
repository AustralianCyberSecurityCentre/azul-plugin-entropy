// Test / Example HTTP Worker in GO
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AustralianCyberSecurityCentre/azul-bedrock/v9/gosrc/events"
	"github.com/AustralianCyberSecurityCentre/azul-bedrock/v9/gosrc/plugin"
	"github.com/AustralianCyberSecurityCentre/azul-entropy.git/entropy"
)

// Maximum number of blocks for calculating entropy, no reason why 800 could be changed in the future.
const Blocks = 800

// 10MB
var maxBufferSize = uint64(10 * 1024 * 1024)

type EntropyPlugin struct {
}

func (ep *EntropyPlugin) GetName() string {
	return "Entropy"
}

func (ep *EntropyPlugin) GetVersion() string {
	return "2025.05.16"
}

func (ep *EntropyPlugin) GetDescription() string {
	return "Calculates entropy over blocks of supplied files."
}

func (ep *EntropyPlugin) GetFeatures() []events.PluginEntityFeature {
	return []events.PluginEntityFeature{
		{Name: "entropy", Type: "float", Description: "Overall entropy calculated for the binary"},
	}
}

func (ep *EntropyPlugin) GetDefaultSettings() *plugin.PluginSettings {
	return plugin.NewDefaultPluginSettings()
}

func (ep *EntropyPlugin) Execute(context context.Context, job *plugin.Job, inputUtils *plugin.PluginInputUtils) *plugin.PluginError {
	bufferedEntropy := entropy.NewBuffered(job.GetSourceEvent().Entity.Size, Blocks)

	endOfFile := false
	var rawChunk []byte
	var err error
	var pluginErr *plugin.PluginError
	startChunk := uint64(0)
	// Calculate entropy
	for !endOfFile {
		rawChunk, endOfFile, pluginErr = job.GetContentChunk(startChunk, startChunk+maxBufferSize)
		if pluginErr != nil {
			return pluginErr
		}
		bufferedEntropy.AppendAndCalculateBufferedValues(rawChunk)
		startChunk += uint64(len(rawChunk))
	}

	entChunks, entSize, entCount := bufferedEntropy.GetChunkEntropySizeAndCount()
	_, err = bufferedEntropy.TotalValue()
	if err != nil {
		return plugin.NewPluginError(plugin.ErrorException, "TotalValue error", "TotalValue error").WithCausalError(err)
	}

	entropyInfo := EventInfoEntropy{
		Overall:    1,
		Blocks:     entChunks,
		BlockSize:  entSize,
		BlockCount: entCount,
	}
	encodedEntropyInfo, err := json.Marshal(&map[string]any{"entropy": entropyInfo})
	if err != nil {
		return plugin.NewPluginError(plugin.ErrorException, "Failed to marshal info", fmt.Sprintf("could not marshal produced entropy info %v", entropyInfo)).WithCausalError(err)
	}
	job.AddInfo(encodedEntropyInfo)
	pluginErr = job.AddFeature("entropy", 1)
	if pluginErr != nil {
		return pluginErr
	}
	return nil
}

func main() {
	pr := plugin.NewPluginRunner(&EntropyPlugin{})
	pr.Run()
}
