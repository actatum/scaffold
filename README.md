# Scaffold

Scaffold is a CLI tool to generate code and folder structure to maintain consistency across the organization.

## Supported Languages
- Go (golang)

## Installation
To install, clone the repo and cd into the root of the project then run the following command:

```go install```

## Usage
1. Edit the .scaffold.json file change the value
of the key "root" to the path where you cloned the scaffold project.

2. Move the .scaffold.json file to your home directory

After that you can verify your installation by running
```scaffold version``` in your terminal.

To scaffold out your Go REST api run 

```scaffold rest -s PROJECTNAME```

alternatively run: ```scaffold r -s PROJECTNAME``` to accomplish the same thing

Use the -s flag to specify your project name 
and in turn the top-level directory that scaffold will create.

That's it once you've scaffolded out your project make sure to replace
all templated code with your own methods and naming