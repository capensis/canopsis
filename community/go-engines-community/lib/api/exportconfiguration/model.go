package exportconfiguration

type ExportDocuments map[int]map[string]interface{}

type Request struct {
	// Possible export values.
	//   * `configuration` - export configuration collection
	//   * `acl` - export default_rights collection
	//   * `pbehavior` - export pbehavior collection
	//   * `pbehavior_type` - export pbehavior_types collection
	//   * `pbehavior_reason` - export pbehavior_reason collection
	//   * `pbehavior_exception` - export pbehavior_exception collection
	//   * `scenario` - export action_scenario collection
	//   * `metaalarm` - export meta_alarm_rules collection
	//   * `idle_rule` - export idle_rules collection
	//   * `eventfilter` - export eventfilter collection
	//   * `dynamic_infos` - export dynamic_infos collection
	//   * `playlist` - export view_playlist collection
	//   * `state_settings` - export state_settings collection
	//   * `broadcast` - export broadcast_message collection
	//   * `associative_table` - export default_associativetable collection
	//   * `notification` - export notification collection
	//   * `view` - export views collection
	//   * `view_tab` - export viewtabs collection
	//   * `widget` - export widgets collection
	//   * `widget_filter` - export widget_filters collection
	//   * `widget_template` - export widget_templates collection
	//   * `view_group` - export viewgroups collection
	//   * `instruction` - export instruction collection
	//   * `job_config` - export job_config collection
	//   * `job` - export job collection
	//   * `resolve_rule` - export resolve_rule collection
	//   * `flapping_rule` - export flapping rule collection
	//   * `user_preferences` - export userpreferences collection
	//   * `kpi_filter` - export kpi_filter collection
	//   * `pattern` - export pattern collection
	//   * `declare_ticket_rule` - export declare_ticket_rule collection
	//   * `link_rule` - export link_rule collection
	Exports []string `json:"export"`
}
