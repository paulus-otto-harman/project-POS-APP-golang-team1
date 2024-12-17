# Project POS App Team 1

Project POS App is chapter 50-55 Team Project

## Installation
clone this project and install dependencies
```bash
go mod tidy
```

## Environment Variables
copy .env.example to .env
```bash
copy .env.example .env
```

## Database Migration and Seeding
by default, migration and seeding are **disabled**  
to enable them, set ```DB_MIGRATION``` and ```DB_SEEDING``` in ```.env``` to ```true```  
to disable, set ```DB_MIGRATION``` and ```DB_SEEDING``` in ```.env``` to ```false```  
to overrides ```.env```, use flags ```-m``` and ```-s```

```bash
cd cmd
go run . -m -s
```

The flag ```-m``` overrides .env ```DB_MIGRATE```  
and flag ```-s``` overrides .env ```DB_SEEDING```

.env | flag | result  |
--- |------|---------|
:white_check_mark:	 | :white_check_mark:  | flag    |
:white_check_mark:	 | :x:  | .env    |
:x: | :white_check_mark:  | flag    |
:x: | :x:  | default |


## Cron Jobs
to run cron jobs, go to folder ```cmd/cron```
```bash
go run .
```
