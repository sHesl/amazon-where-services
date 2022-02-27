# AWS (amazon-where-services)

Chasing down which AWS services are available in which regions can be a bit of a pain; this repo 
aims to keep an up to date reference for the available endpoints for all AWS services. 

It achieves this by parsing the `botocore/data/endpoints.json` file from the `develop` branch of 
`botocore` ([see here](https://github.com/boto/botocore/blob/develop/botocore/data/endpoints.json)), 
and presenting the endpoints in a variety of different ways (by partition, by service, grouping all 
global services together etc).


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
