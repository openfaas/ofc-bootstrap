# How to make a godoc
Go can make package and documentise automatically.

```
// what you want to specify
package calc

// what is this func doing
func Sum(a int, b int) int {
	return a + b
}
```

you can add your comments upside package and keyword function.
do `godoc` in your `$GOPATH`

```
~$ cd $GOPATH
~/hello_project$ godoc calc
PACKAGE DOCUMENTATION

package calc
    import "calc"

    calculation package

FUNCTIONS
func Sum(a int, b int) int
    add to integers
```

You can find a result. You can also do `godoc package_name function_name` to print the specific information.

```
~$ cd $GOPATH
~/hello_project$ godoc calc Sum
func Sum(a int, b int) int
    add to integers
```

and also you can do a website version of your own godoc.

```
~$ cd $GOPATH
~/hello_project$ godoc -http=:6060

```

`godoc -http=:<port number>` can show your own godoc by your own local.

visit your `http://127.0.0.1:6060/pkg/calc` and find your own godoc.
