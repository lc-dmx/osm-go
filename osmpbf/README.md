# osm-go/osmpbf

Package osmpbf provides a writer for encoding large [OSM PBF](https://wiki.openstreetmap.org/wiki/PBF_Format) files.

## Example

```go
outputFile, err := os.Create("output.osm.pbf")
if err != nil {
    panic(err)
}
defer outputFile.Close()

w, err := osmpbf.NewWriter(context.Background(), outputFile)
if err = w.WriteEntity(entity.NewNode(1)); err != nil {
	panic(err)
}
if err = w.WriteEntity(entity.NewWay(2)); err != nil {
	panic(err)
}
if err = w.WriteEntity(entity.NewRelation(3)); err != nil {
	panic(err)
}

if err = w.Close(); err != nil {
	panic(err)
}
```

## Pay Attention to

1. When calculating the space occupied by each fileblock, we adopted the method of estimating based on field types, 
because it is impossible to marshal and get the actual size every time we write an entity. At the same time, 
because of the use of protobuf variable-length encoding, there will be some deviations in the estimate. 
There are two optimization methods that can be thought of here: 
one is to customize the space estimation for variable-length types; 
the other is to limit the entity quantity for each fileblock (For example, the OSM official Java library limits 8000 
in version 0.38). In actual use, we haven't added any optimization, but the result is still acceptable.

2. Theoretically, within each fileblock, if entity IDs are in order, the optimal compression effect will be achieved. 
However, because writing is streaming, this cannot be done in the SDK. If you want to sort, you can only do it 
previously. During the test, we found that if it is disordered, it will occupy 40%-60% more space. If it is 
divided into two parts and each part is ordered internally, it will occupy an additional 10% of space. Notice: this 
result has not been verified for many rounds.

## TODO

1. Writing benchmark tests.
2. If the performance does not meet the requirements, consider to develop parallel write.
