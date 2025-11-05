package transport

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

const DefaultChunkSize = 256 * 1024

type IPluginDataProvider interface {
	GetPluginMetadata(pluginName string) (*models.PluginTransferMetadata, error)
	GetPluginChunk(pluginName string, chunkIndex int) (*models.PluginChunkResponse, error)
}

type PluginDataProvider struct {
	pluginDir string
	chunkSize int
}

func NewPluginDataProvider(pluginDir string) *PluginDataProvider {
	return &PluginDataProvider{
		pluginDir: pluginDir,
		chunkSize: DefaultChunkSize,
	}
}

func (p *PluginDataProvider) GetPluginMetadata(pluginName string) (*models.PluginTransferMetadata, error) {
	pluginPath := filepath.Join(p.pluginDir, pluginName+".so")
	
	fileInfo, err := os.Stat(pluginPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("plugin not found: %s", pluginName)
		}
		return nil, fmt.Errorf("failed to access plugin: %w", err)
	}

	hash, err := p.calculateFileHash(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}

	totalSize := fileInfo.Size()
	totalChunks := int(totalSize / int64(p.chunkSize))
	if totalSize%int64(p.chunkSize) != 0 {
		totalChunks++
	}

	metadata := &models.PluginTransferMetadata{
		Name:        pluginName,
		Hash:        hash,
		TotalSize:   totalSize,
		ChunkSize:   p.chunkSize,
		TotalChunks: totalChunks,
	}

	logging.Info("Plugin metadata: %s (size: %d bytes, chunks: %d)", pluginName, totalSize, totalChunks)
	return metadata, nil
}

func (p *PluginDataProvider) GetPluginChunk(pluginName string, chunkIndex int) (*models.PluginChunkResponse, error) {
	pluginPath := filepath.Join(p.pluginDir, pluginName+".so")
	
	file, err := os.Open(pluginPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("plugin not found: %s", pluginName)
		}
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat plugin: %w", err)
	}

	totalSize := fileInfo.Size()
	totalChunks := int(totalSize / int64(p.chunkSize))
	if totalSize%int64(p.chunkSize) != 0 {
		totalChunks++
	}

	if chunkIndex < 0 || chunkIndex >= totalChunks {
		return nil, fmt.Errorf("chunk index %d out of range [0, %d)", chunkIndex, totalChunks)
	}

	offset := int64(chunkIndex) * int64(p.chunkSize)
	remainingBytes := totalSize - offset
	chunkSize := p.chunkSize
	if remainingBytes < int64(p.chunkSize) {
		chunkSize = int(remainingBytes)
	}

	chunk := make([]byte, chunkSize)
	_, err = file.ReadAt(chunk, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to read chunk: %w", err)
	}

	chunkHash := fmt.Sprintf("%x", sha256.Sum256(chunk))
	encodedChunk := base64.StdEncoding.EncodeToString(chunk)

	response := &models.PluginChunkResponse{
		PluginName:  pluginName,
		ChunkIndex:  chunkIndex,
		TotalChunks: totalChunks,
		Data:        encodedChunk,
		Hash:        chunkHash,
	}

	logging.Debug("Chunk %d/%d of plugin %s (%d bytes)", chunkIndex+1, totalChunks, pluginName, chunkSize)
	return response, nil
}

func (p *PluginDataProvider) calculateFileHash(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}
