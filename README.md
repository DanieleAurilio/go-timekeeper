## TimeKeeper

timekeeper is a Golang side project, developed for learning purposes and automate scripts. So it can be full of mistake. Not a Golang puritan.

### How to configure
 
 - Rename the file 'config.default.json' to 'config.json'.

 - Follow these suggestion for configure the **config.json**:

```json
{
    "poolms": 60000, // To be declared in ms, this is the frequency which jobs are pulled
    "outputDir": "", // The directory output, where executable are located 
    "jobs": [
        {
            "id": "", // The ID of the job
            "enable": true, // Used for enable or disable job
            "filename": "", // The name of the file .go that should be builded
            "params": { // Here you can declare some params used by the script
                "": "" // Be aware, all params must be declared as string
            },
            "schedule": { // The time when the script should be executed
                "hours": 0, // Hours is 24h format. i.e. 23 or 10
                "minutes": 0,   
                "seconds": 0,
                "month": "", // i.e. May, December, March OR *. The wildcard enable this job all months of the year
                "weekDay": "" // i.e. Monday,Tuesday,Sunday OR *. The wildcard enable this job all days of the week
            }
        }
    ]
}
```

- This must be the folder tree. The most important is the **jobs** folder with subfolders and files:

```bash

├── README.md
├── config.json
├── go.mod
├── jobs
│   ├── test1
│   │   └── test1.go
│   └── test2
│       └── test2.go
├── main.go
├── models
│   ├── models.go
│   └── queue
│       └── queue.go
├── output
│   └── test1
├── timekeeper
│   └── timekeeper.go
└── utils
    └── utils.go

```
- You can retrive the **params** declared for a certain job in this way:

```go

package main

func main() {
    var params map[string]string
    json.Unmarshal([]byte(strings.Join(os.Args[2:], "")), &params)
    fmt.Printf("%v", params)
}

```

### Execution

Now you can run these commands:

```shell
    go run main.go
```
OR

```shell
    go build -o ${main}
    ./${main}
```
In this case what you need this tree folders structure:

```bash
├── config.json
├── jobs
│   ├── test1
│   │   └── test1.go
│   └── test2
│       └── test2.go
├── output
│   └── test1
├── main
```

### Note 
- You may need to set permissions to the output folder.
- Not tested on windows environment.

