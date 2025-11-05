package models

type PluginTransferMetadata struct {
	Name        string `json:"name"`
	Hash        string `json:"hash"`
	TotalSize   int64  `json:"total_size"`
	ChunkSize   int    `json:"chunk_size"`
	TotalChunks int    `json:"total_chunks"`
}

type PluginChunkRequest struct {
	PluginName string `json:"plugin_name"`
	ChunkIndex int    `json:"chunk_index"`
}

type PluginChunkResponse struct {
	PluginName  string `json:"plugin_name"`
	ChunkIndex  int    `json:"chunk_index"`
	TotalChunks int    `json:"total_chunks"`
	Data        string `json:"data"`
	Hash        string `json:"hash"`
}

type PluginTransferStatus struct {
	PluginName      string  `json:"plugin_name"`
	TotalChunks     int     `json:"total_chunks"`
	ReceivedChunks  int     `json:"received_chunks"`
	PercentComplete float64 `json:"percent_complete"`
	Status          string  `json:"status"`
}
