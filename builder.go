package local

type Builder struct {
	api_app_builder.SimpleBaseBuilder

	parent api_api.API
	settings Settings
	ops api_operation.Operations
}

func NewBuilder(settings Settings) *Builder {
	return &Builder{
		settings: settings,
	}
}

func (local *Builder) Builder() api_builder.Builder {
	return api_builder.Builder(local)
}

func (local *Builder) SetParent(parent api_api.API) {
	base.parent = parent
}

func (local *Builder) Builder() api_builder.Builder {
	return api_builder.Builder(local)
}

func (local *Builder) Activate([]string, SettingsProvider) {

}

func (local *Builder) Validate() api_result.Result {
	
}

func (local *Builder) Operations() api_operation.Operations {
	return local.ops
}
