package elastic

// IndexInfo represents information about an Elasticsearch index
type IndexInfo struct {
	Index     string `json:"index"`
	Status    string `json:"status"`
	Health    string `json:"health"`
	DocsCount string `json:"docs.count"`
	StoreSize string `json:"store.size"`
	PriShards string `json:"pri"`
	RepShards string `json:"rep"`
	UUID      string `json:"uuid"`
}

// ClusterStats represents Elasticsearch cluster statistics
type ClusterStats struct {
	ClusterName string       `json:"cluster_name"`
	ClusterUUID string       `json:"cluster_uuid"`
	Timestamp   int64        `json:"timestamp"`
	Status      string       `json:"status"`
	Indices     IndicesStats `json:"indices"`
	Nodes       NodesStats   `json:"nodes"`
}

// IndicesStats represents statistics about indices
type IndicesStats struct {
	Count int `json:"count"`
	Docs  struct {
		Count   int64 `json:"count"`
		Deleted int64 `json:"deleted"`
	} `json:"docs"`
	Store struct {
		Size string `json:"size_in_bytes"`
	} `json:"store"`
	Fielddata struct {
		Memory    string `json:"memory_size_in_bytes"`
		Evictions int64  `json:"evictions"`
	} `json:"fielddata"`
}

// NodesStats represents statistics about nodes
type NodesStats struct {
	Count struct {
		Total               int `json:"total"`
		CoordinatingOnly    int `json:"coordinating_only"`
		Data                int `json:"data"`
		Ingest              int `json:"ingest"`
		Master              int `json:"master"`
		RemoteClusterClient int `json:"remote_cluster_client"`
	} `json:"count"`
	Versions []string `json:"versions"`
	OS       struct {
		AvailableProcessors int `json:"available_processors"`
		AllocatedProcessors int `json:"allocated_processors"`
	} `json:"os"`
	Process struct {
		CPU struct {
			Percent int `json:"percent"`
		} `json:"cpu"`
		OpenFileDescriptors struct {
			Min int64 `json:"min"`
			Max int64 `json:"max"`
			Avg int64 `json:"avg"`
		} `json:"open_file_descriptors"`
	} `json:"process"`
}

// ClusterHealth represents Elasticsearch cluster health information
type ClusterHealth struct {
	ClusterName                 string                 `json:"cluster_name"`
	Status                      string                 `json:"status"`
	TimedOut                    bool                   `json:"timed_out"`
	NumberOfNodes               int                    `json:"number_of_nodes"`
	NumberOfDataNodes           int                    `json:"number_of_data_nodes"`
	ActivePrimaryShards         int                    `json:"active_primary_shards"`
	ActiveShards                int                    `json:"active_shards"`
	RelocatingShards            int                    `json:"relocating_shards"`
	InitializingShards          int                    `json:"initializing_shards"`
	UnassignedShards            int                    `json:"unassigned_shards"`
	DelayedUnassignedShards     int                    `json:"delayed_unassigned_shards"`
	NumberOfPendingTasks        int                    `json:"number_of_pending_tasks"`
	NumberOfInFlightFetch       int                    `json:"number_of_in_flight_fetch"`
	TaskMaxWaitingInQueueMillis int                    `json:"task_max_waiting_in_queue_millis"`
	ActiveShardsPercentAsNumber float64                `json:"active_shards_percent_as_number"`
	Indices                     map[string]IndexHealth `json:"indices,omitempty"`
}

// IndexHealth represents health information for a specific index
type IndexHealth struct {
	Status              string `json:"status"`
	NumberOfShards      int    `json:"number_of_shards"`
	NumberOfReplicas    int    `json:"number_of_replicas"`
	ActivePrimaryShards int    `json:"active_primary_shards"`
	ActiveShards        int    `json:"active_shards"`
	RelocatingShards    int    `json:"relocating_shards"`
	InitializingShards  int    `json:"initializing_shards"`
	UnassignedShards    int    `json:"unassigned_shards"`
}
