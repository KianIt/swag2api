# swag2api

**Swag2API** uses Go annotations supported by the [swag](https://github.com/swaggo/swag) tool and Go function definitions to automatically build a simple API for HTTP requests. It alllows you to use common Go functions as HTTP handlers by just adding proper annotations to it.

## Getting started

1. Write your functions into a file and provide them with annotations is [swag format](https://github.com/swaggo/swag?tab=readme-ov-file#declarative-comments-format:~:text=Swagger%20Extensions-,Declarative%20Comments%20Format,-General%20API%20Info). Every annotation must contain the `@ID` label that matches the function's name. Also, there must be a `main.go` file with **swag** API root information.

2. Install **swag2api** by using:

```sh
go install github.com/KianIt/swag2api/cmd/swag2api@latest
```

3. Run the API generation in the directory of the file with functions and the `main.go` file:

    a. Add the next line at the beginning of your file with functions:

    ```go
    //go:generate swag2api
    ```

    and use the Golang default generation tool:

    ```sh
    go generate
    ```

    b. Use the swag2api directly:

    ```
    swag2api
    ```

4. Then there is the `generated.go` file with HTTP request handlers for the annotated functions. All the request handlers are avaliable via the `s2aHandler` dispatching requests to the generated APT. API is being setup in the `init` function in the generated file, so you must import the API package or run the file to obtain the non-empty `s2aHandler`.

## swag2api cli

```sh
Usage of ./swag2api:
  -handler string
        Name of the API HTTP handler (default "s2aHandler")
  -main string
        Name of the swag main file (default "main.go")
  -pkg string
        Path to the Golang package (default ".")
  -to string
        Name of the generated API file (default "generated.go")
```

# Examples

## Package

By default, **swag2api** generates API by the files from the current Golang package. You can specify the path to the package by passing the `-pkg` flag to the tool. For example, to generate the API in the package `path/to/package` use:

```
swag2api -pkg path/to/package
```

## Swag main file

By default, **swag2api** uses the `main.go` file as the swag main file. You can specify the name of the swag main file by passing the `-main` flag to the tool. For example, to use the `mymain.go` file as the main swag file use:

```
swag2api -main mymain.go
```

## Generated file

By default, **swag2api** generates the `generated.go` file with HTTP request handlers. You can specify the name of the generaed file by passing the `-to` flag to the tool. For example, to use the `mygenereated.go` file as the generated file use:

```
swag2api -to mygenereated.go
```


## Handler

All the HTTP request handlers are available via a single HTTP multiplexer dispatching requests to the handlers. By default, swag2api generated the `s2aHandler` declaration:

```
var s2aHandler http.Handler
```

You can specify the name of the multiplexer by passing the `-handler` flag to the tool. For example, if you pass `-handler myHandler` there will be generated the declaration:

```
var myHandler http.Handler
```

If you pass a name of an existing handler then there won't be generated the declaration, but the handler will be rewritten with the multiplexer for the generated API. For example, see the [example](example).

## Annotations

The example of the function annotation:

```go
// method1 godoc
// @ID method1
// @Param pathString path string true " "
// @Param pathInt path int true " "
// @Param pathFloat64 path float64 true " "
// @Param pathBool path bool true " "
// @Param pathBytes path []byte true " "
// @Router /path-to-method1 [get]
func method1(pathString string, pathInt int, pathFloat64 float64, pathBool bool, pathBytes []byte) (result string, err error) {
	return "success", nil
}
```

The annotation contains the `@ID` label that matches the function name. So the parser can join the annotation to the function.

Params from the annotation's `@Param` label must fit the function's params.\
If there is a param from the annotation that doesn' t exist in the function feclaration, there will be a warning.\
If there is a param from the function declaration, that doesn't exist in the annotation, there wiill be an error.

The HTTP request handler fot the annotated function:

```go
func _handler_method1(w http.ResponseWriter, r *http.Request) {
	pathStringPath := r.PathValue("pathString")
	pathString, pathStringUnmarshalErr := _unmarshalString[string](pathStringPath)
	if pathStringUnmarshalErr != nil {
		_handleBadRequest(w, pathStringUnmarshalErr)
		return
	}
	pathIntPath := r.PathValue("pathInt")
	pathInt, pathIntUnmarshalErr := _unmarshalString[int](pathIntPath)
	if pathIntUnmarshalErr != nil {
		_handleBadRequest(w, pathIntUnmarshalErr)
		return
	}
	pathFloat64Path := r.PathValue("pathFloat64")
	pathFloat64, pathFloat64UnmarshalErr := _unmarshalString[float64](pathFloat64Path)
	if pathFloat64UnmarshalErr != nil {
		_handleBadRequest(w, pathFloat64UnmarshalErr)
		return
	}
	pathBoolPath := r.PathValue("pathBool")
	pathBool, pathBoolUnmarshalErr := _unmarshalString[bool](pathBoolPath)
	if pathBoolUnmarshalErr != nil {
		_handleBadRequest(w, pathBoolUnmarshalErr)
		return
	}
	pathBytesPath := r.PathValue("pathBytes")
	pathBytes, pathBytesUnmarshalErr := _unmarshalString[[]byte](pathBytesPath)
	if pathBytesUnmarshalErr != nil {
		_handleBadRequest(w, pathBytesUnmarshalErr)
		return
	}
	result, err := method1(pathString, pathInt, pathFloat64, pathBool, pathBytes)
	type _resultType struct {
		Result string `json:"result"`
	}
	_resultValue := _resultType{Result: result}
	_handleResult(w, err, _resultValue)
}

...

func init() {
	_mux := http.NewServeMux()
	_mux.HandleFunc("GET /path-to-method1", _handler_method1)
	...
	s2aHandler = _mux
}
```

## HTTP statuses

You can manipulate with the request handler's HTTP response status code. To do this you need to add an error result to your function and to wrap it with a wrapper from the package

```sh
github.com/KianIt/swag2api/statuses
```

For example, function that returns an error wrapped with `NotFoundError` wrapper:

```go
// method9 godoc
// @ID method9
// @Router /path-to-method9 [get]
func method9() (result string, err error) {
	return "failed", s2aStatuses.NotFoundError(errors.New("test error"))
}
```

This will change the response status code to `404` and return an error response.

Another example, function that returns a `nil` error wrapped with `NotFoundError` wrapper:

```go
// method8 godoc
// @ID method8
// @Router /path-to-method8 [get]
func method8() (result string, err error) {
	return "success", s2aStatuses.NotFoundError(nil)
}
```

This will also change the response status code to `404`, but the result will be successfull.


## More

You can find more examples in the [example](example) directory.
