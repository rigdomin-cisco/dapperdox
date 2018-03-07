Title: Controlling the order of methods
Description: How to control the order that methods are displayed
Keywords: method, ordering, operation, name, specification, swagger, OpenAPI

# Controlling the order of methods

## The default method order

By default, DapperDox will use the `summary` member of method declarations to determine the order
that methods are listed in the navigation and summary pages, effectively listing them in alphabetical
order of `summary`.

This may not be appropriate for your API, for example you may feel that your REST API should have its
methods ordered by [operation name](/docs/glossary-terms#operation-name). DapperDox provides an
extenstion to give you a flexible level of control.

## Changing the method order

DapperDox allows specification writers to set the order that methods are listed. To do this, add
an `x-sortMethodsBy` extension to the top level OpenAPI specification.

This extension takes an `ARRAY` of values. For example:

```json
{
    "swagger": "2.0",
    "x-sortMethodsBy": ["path","operation","summary"],
    "info": {
        "title": "Example API",
        "description" : "Showing how x-sortMethodsBy is used",
        "version": "1.0.0"
    }
}
```
This tells DapperDox to sort methods first by path, then those with the same path will be sorted
by operation, and finaly, those with the same operation will be sorted by summary.

`x-sortMethodsBy` can take the following values:

| value      | description |
|------------|-------------|
| `path`       | The path/URL of the method |
| `method`     | The HTTP method-name, such as `get`, `post` and `put`. |
| `operation`  | The `operation-name` of the method. Contents controlled by [x-operationName](/docs/spec-method-names#methods-and-operations) |
| `summary`    | The `summary` member of the method declaration |
| `navigation` | If [navigating by method name](/docs/spec-method-names#navigating-by-method-name), then this is equivalent to `summary`, otherwise it equivalent to `operation` |

You do not need to supply multiple values to the `x-sortMehtodsBy` array, but where you are sorting
by an entity that contains duplicate values in your API (such as `operation` and there are multiple
`get`s defined), then you may find that the order of these duplicate `operation` values is not consistent. Supplying additional sort criteria allows you to lock this ordering down.


!!!HIGHLIGHT!!!
