package defs

#outputdef: {
    alias: string // alias used to store the output value
    path:  string // jq expression to extract the value from the task result
}

#TFTask: {
    @task(opentf.TF)
    config: _
    dep: _
    out: _
    outputs?: [...#outputdef]
    inputs?: _
}

#providers: {
    provider: {
        "aws": {}
    }
    terraform: {
        required_providers: {
            aws: {
                source:  "hashicorp/aws"
                version: ">= 4.67.0"
            }
        }
    }
}

inputs: _

common: {
    project_name: "static-website"
    vpc_cidr:     "10.0.0.0/16"
    db_parameter_group_name: "default.postgres14"
    db_subnet_group_name: "default"
    db_password: "edu"
    db_username: "edu"
}

locals: {
    azs: ["us-west-1a", "us-west-1c"] 
    private_subnets: ["10.0.1.0/24", "10.0.2.0/24"]
    public_subnets:  ["10.0.101.0/24", "10.0.102.0/24"]
}