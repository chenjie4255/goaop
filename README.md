## GOAOP

Go does not support Aspect-Oriented Programming natively, so it's not easy to do something like logging for method of interface, it's not beatiful at least.

Insprited by Go Generate, I found out that it's possible to write an AOP Framework to help us achieve this goal in an interesting way.

For example, if you want to measure the elaplsed time for your every method of interface, you may need to write something like this:

```
type DB interface {
    GetUserCount() (int, error)
}

func (d *db) GetUserCount(int, error) {
    t := time.Now()
    defer func() {
        fmt.Println("elapsed time:", time.Now().Sub(t).String())
    }()

    // db handle
    ...
}

```

the interface DB above just have one method for now, what happen if this interface have more than 10 method? you may need to write lots of extremely repeated code, it's not graceful and easy to go wrong.

## How GoAOP Works?

GoAOP could help you generate those repeatable codes and provide you an entry to handle every method call in a union form.

### install

```
go get github.com/chenjie4255/goaop
go install github.com/chenjie4255/goaop
```

### generate aop codes

At first you need to add one line ```//go:generate goaop -f=$GOFILE``` into your .go file so the  ```go generate``` could recognize it, then add one comment line ```// @ifmeasure```above the interface you wanted to generate aop codes, for example:

```
// @ifmeasure
type DB interface{
    GetConnectionCount() (int, error)
}
```

after finished above steps, you could go into the folder where this file placed, and execute the below command
```
go genearte
```

niced code will be generated automaticlly.

## More Example

see example/example.go

## Notice

this project is still in developing, any breaking-changes could be made in the feature, please use go vendor to maintain versions  