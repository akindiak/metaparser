Simple parser for internal usage to parse dbt yml files with models and get data from meta fields to prepare json files with defined structure for later use.

First do make build.

Then call go run ./... -path= -system= arg1 arg2 ...
arguments:
    path - dbt project path
    system - name of the folder where to store json metafiles

    all other args should be model/table names