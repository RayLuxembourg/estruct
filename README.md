# EStruct

estruct traverses a modern javascript project and maps all its dependencies and relationships and outputs a JSON which
is used to build network visualizations of the project.

---

Getting started

```shell
go get github.com/RayLuxembourg/estruct
```

```go
root := "/path/to/project/src"
labels:= make([]estruct.Label,0) // or create real labels
p:= estruct.NewConfig(root, `(.(js|jsx))$`,labels)
relativePath := "./src"
datasets, fileMap, dependenciesMap := p.Init(relativePath) // the output is the json
```

Saving output to a json file

```go
jsonName := "application.json"
os.Remove(jsonName)
b, _ := json.Marshal(datasets)

ioutil.WriteFile(jsonName, b, 0666)
```

Output json structure

```json
[
  {
    "id": "52FDFC07-2182-654F-163F-5F0F9A621D72",
    "path": "/Users/myUser/work/myProject/src/Component.js",
    "dependencies": [
      "__webpack_public_path__",
      "react",
      "react-dom",
      "app-data"
    ],
    "folder": "src",
    "extension": "js",
    "name": "ECMComponent",
    "lines": 106,
    "label": ""
  },
  {
    "id": "9566C74D-1003-7C4D-7BBB-0407D1E2C649",
    "path": "/Users/myUser/work/myProject/src/__webpack_public_path__.js",
    "dependencies": [],
    "folder": "src",
    "extension": "js",
    "name": "__webpack_public_path__",
    "lines": 3,
    "label": ""
  },
  {
    "id": "81855AD8-681D-0D86-D1E9-1E00167939CB",
    "path": "/Users/myUser/work/myProject/src/app-data.jsx",
    "dependencies": [],
    "folder": "src",
    "extension": "jsx",
    "name": "app-data",
    "lines": 4,
    "label": ""
  }
]
```

EStruct will generate unique id for each mapped file to allow relationship creation using those id's

### How is this helpful to your organization ?

Using the output of EStruct, you can visualize your complex large scale JavaScript project using tools like neo4j and
then ask real world questions from the graph database.
![](https://i.imgur.com/uUzJeCL.png)

You can also integrate the EStruct with your devops pipeline to keep track of your project architecture and changes.

Todo list:

1. make EStruct plugable/extendable to support various project structures
2. Generate labels based on user patterns (analyze file name and decide which label to give it)
3. Add support for require syntax.
4. Create CLI tool
5. Publish using NPM
6. Read configuration from project root directory .estruct.json
7. Boost performance
8. Cleaner code
