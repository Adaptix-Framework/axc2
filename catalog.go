package adaptix

type AgentCatalogItem struct {
	Name      string   `json:"name"`
	Listeners []string `json:"listeners"`
	AXS       string   `json:"axs,omitempty"`
}

type ListenerCatalogItem struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Type     string `json:"type"`
	AXS      string `json:"axs,omitempty"`
}

type ServiceCatalogCommand struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Example     string `json:"example,omitempty"`
	Message     string `json:"message,omitempty"`
	Destructive bool   `json:"destructive,omitempty"`
}

type ServiceCatalogItem struct {
	Name        string                  `json:"name"`
	Title       string                  `json:"title,omitempty"`
	Description string                  `json:"description,omitempty"`
	RPC         bool                    `json:"rpc"`
	Commands    []ServiceCatalogCommand `json:"commands"`
}
