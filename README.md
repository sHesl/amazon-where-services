# AWS (amazon-where-services)

Chasing down which AWS services are available in which regions is a bit of a pain; the AWS docs
aren't very grepable and parsing `botocore/data/endpoints.json` is effort. 

This dumb-bot just keeps an up to date representation of what runs where :)

### Structure
A list of global services, relative to their partition, can be found at:
- `./global-services/$partition/$region`

A list of which services are available, relative to their partition, can be found at:
- `./partitions/$partition`

A list of regional services, relative to their partition, can be found at:
- `./regional-services/$partition/$region`

A list of all services available in a given region, irrespective of partition, can be found at:
- `./regions/$region`

A list of which regions a given service operates in, irrespective of partition, can be found at:
- `./services/$region`

A list of single-region services, relative to their partition, can be found at:
- `./single-region-services/$partition/$region`
