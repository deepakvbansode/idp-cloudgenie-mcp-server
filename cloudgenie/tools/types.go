package tools

// CreateResourceInput defines the input parameters for creating a resource
type CreateResourceInput struct {
	Name          string                 `json:"name" jsonschema:"Resource name"`
	BlueprintName string                 `json:"blueprintName" jsonschema:"Blueprint name"`
	Description   string                 `json:"description,omitempty" jsonschema:"Resource description"`
	Spec         map[string]interface{} `json:"spec" jsonschema:"Resource properties"`
}

// CreateResourceOutput defines the output returned after creating a resource
type CreateResourceOutput struct {
	Name   string `json:"name" jsonschema:"Resource name"`
	Status map[string]interface{} `json:"status" jsonschema:"Resource status"`
}


type Blueprint struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Version     string                 `json:"version"`
	Category    string                 `json:"category"`
	Parameters  map[string]interface{} `json:"parameters"`
}


// Resource represents a CloudGenie resource returned from the API
type Resource struct {
	Id            string                 `json:"id"`
	Name          string                 `json:"name"`
	BlueprintName string                 `json:"blueprint_name"`
	Description   string                 `json:"description,omitempty"`
	Status        map[string]interface{} `json:"status"`
	CreatedAt     string                 `json:"created_at"`
	UpdatedAt     string                 `json:"updated_at"`
	Spec         map[string]interface{} `json:"spec"`
}

type GetBlueprintInfoInput struct {
	BlueprintName string `json:"blueprint_name" jsonschema:"Blueprint name"`
}

type GetBlueprintInfoOutput struct {
	Blueprint Blueprint `json:"blueprint" jsonschema:"Blueprint details"`
}

type GetBlueprintsInput struct {}

type GetBlueprintsOutput struct {
	Blueprints []Blueprint `json:"blueprints" jsonschema:"List of blueprints"`
}

type GetResourcesInput struct {
	// No input needed for listing all resources
}

type GetResourcesOutput struct {
	Resources []Resource `json:"resources" jsonschema:"List of resources"`
}

type GetResourceByNameInput struct {
	Name string `json:"name" jsonschema:"Resource name"`
}

type GetResourceByNameOutput struct {
	Resource *Resource `json:"resource" jsonschema:"Resource details"`
}
