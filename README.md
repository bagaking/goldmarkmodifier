# goldmarkmodifier

**goldmarkmodifier** is a subproject of kbs (WIP), which is responsible for providing the **"Markdown AST editing capabilities"** needed to implement the **"markdown project"** in kbs.

The goldmarkmodifier is a tool that allows you **to modify the ast of markdown**, extending the functionality of [goldmark](https://github.com/yuin/goldmark) to allow you to edit the ast of markdown generated by goldmark.

With goldmarkmodifier you can easily edit the markdown asts parsed by goldmark as follows:

- Insert arbitrary text
- Parse and insert another Markdown file

In addition, goldmarkmodifier also provides some quick ways to manipulate ast by providing a data structure called **Mapping**, which allows you to set up matching rules and **add**, **delete**, **replace**, etc. to the matched structure in ast.

## Getting started

Execute `go run ./example` to see the example