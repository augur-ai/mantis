package schemas

#ModuleRDS: {
	source:               string
	version:              string
	identifier:           string
	engine:               string
	engine_version:       string
	instance_class:       string
	allocated_storage:    int
	username:             string
	password:             string
	db_subnet_group_name: string | *null
	vpc_security_group_ids: [string] | *null
	parameter_group_name: string
	publicly_accessible:  bool
	skip_final_snapshot:  bool
}
