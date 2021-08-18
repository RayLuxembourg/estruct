# EStruct

EStruct traverses javascript projects and maps all the dependencies and relationships to a JSON. 
The output can be used to build network visualizations of the project and document the architecture.

---

### Installations

```shell
go install github.com/RayLuxembourg/estruct@latest
```
---
After installation, you can run estruct from the terminal and by default it will run on your current directory

In order to view list of possible arguments and their default value, use
```shell
estruct -h
```
---
### Running estruct with specific arguments
```shell
 estruct -p /path/to/code -e "ts|tsx"                  
```

### Output json structure example

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

### Todo list

- [x] Parse ES Modules
- [x] Create CLI tool
- [ ] Use abstract syntax tree for better code analysis 
- [ ] Read configuration from project root directory .estructrc.json
- [ ] Support CommonJS
- [ ] Generate labels based on user patterns (analyze file name and decide which label to give it)
- [ ] Publish using NPM
- [ ] Boost performance
- [ ] Cleaner code
- [ ] make EStruct plugable/extendable to support various project structures

