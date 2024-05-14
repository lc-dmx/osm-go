# The File xx.proto is from https://github.com/openstreetmap/OSM-binary/tree/master/osmpbf
# Change：
#   1. change the type of field 'StringTable.s' from []byte to string
pb_gen:
	protoc -I=. osmpbf/model_pb/*.proto --go_out=./osmpbf