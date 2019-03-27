package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO test array of structs with defaults

type Simple struct {
	Name    string
	Age     int64
	OptOut  bool
	Balance float64
}

func TestSimpleStruct(t *testing.T) {
	configString := `
	Name: "stephen"
	Age: 28
	OptOut: true
	Balance: 5.5
	`

	config := Simple{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, "stephen", config.Name)
	require.Equal(t, int64(28), config.Age)
	require.Equal(t, true, config.OptOut)
	require.Equal(t, 5.5, config.Balance)
}

func TestDefaults(t *testing.T) {
	configString := `
	Age: 15
	`

	config := Simple{
		Name:    "stephen",
		Age:     28,
		OptOut:  true,
		Balance: 22.3,
	}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, "stephen", config.Name)
	require.Equal(t, int64(15), config.Age)
	require.Equal(t, true, config.OptOut)
	require.Equal(t, 22.3, config.Balance)
}

func TestStrictWithAllFields(t *testing.T) {
	configString := `
	Name: "stephen"
	Age: 28
	OptOut: true
	Balance: 5.5
	`

	config := Simple{
		Name:    "zero",
		Age:     32,
		OptOut:  false,
		Balance: 22.3}

	err := LoadConfigFromString(configString, &config, true)
	require.NoError(t, err)
	require.Equal(t, "stephen", config.Name)
	require.Equal(t, int64(28), config.Age)
	require.Equal(t, true, config.OptOut)
	require.Equal(t, 5.5, config.Balance)
}

func TestStrictWithMissingFields(t *testing.T) {
	configString := `
	Name: "stephen"
	Age: 28
	OptOut: true
	`

	config := Simple{}

	err := LoadConfigFromString(configString, &config, true)
	require.Error(t, err)
}

func TestStringBadValue(t *testing.T) {
	configString := `
	Name: 23
	`

	config := Simple{}

	err := LoadConfigFromString(configString, &config, false)
	require.Error(t, err)
}

func TestStringNoQuotes(t *testing.T) {
	configString := `
	Name: alpha
	`

	config := Simple{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, "alpha", config.Name)
}

func TestBoolAsString(t *testing.T) {
	configString := `
	OptOut: "true"
	`

	config := Simple{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, true, config.OptOut)
}

func TestBoolBadValue(t *testing.T) {
	configString := `
	OptOut: 32
	`

	config := Simple{}

	err := LoadConfigFromString(configString, &config, false)
	require.Error(t, err)
}

type Ints struct {
	Big    int64
	Medium int32
	Small  int16
	Tiny   int8
	Any    int
}

func TestIntTypes(t *testing.T) {
	configString := `
	Big: 1000000000
	Medium: 100000
	Small: 1000
	Tiny: 100
	Any: 100000
	`

	config := Ints{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, int64(1000000000), config.Big)
	require.Equal(t, int32(100000), config.Medium)
	require.Equal(t, int16(1000), config.Small)
	require.Equal(t, int8(100), config.Tiny)
	require.Equal(t, int(100000), config.Any)
}

func TestIntTypesAsStrings(t *testing.T) {
	configString := `
	Big: "1000000000"
	Medium: "100000"
	Small: "1000"
	Tiny: "100"
	Any: "100000"
	`

	config := Ints{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, int64(1000000000), config.Big)
	require.Equal(t, int32(100000), config.Medium)
	require.Equal(t, int16(1000), config.Small)
	require.Equal(t, int8(100), config.Tiny)
	require.Equal(t, int(100000), config.Any)
}

func TestIntTypesBadValue(t *testing.T) {
	configString := `
	Big: "a"
	Medium: "a"
	Small: "a"
	Tiny: "a"
	Any: "a"
	`

	config := Ints{}

	err := LoadConfigFromString(configString, &config, false)
	require.Error(t, err)
}

type Floats struct {
	Big    float64
	Medium float32
}

func TestFloatTypes(t *testing.T) {
	configString := `
	Big: 1000000.1888
	Medium: 100000.321
	`

	config := Floats{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, float64(1000000.1888), config.Big)
	require.Equal(t, float32(100000.321), config.Medium)
}

func TestFloatTypesAsStrings(t *testing.T) {
	configString := `
	Big: "1000000.1888"
	Medium: "100000.321"
	`

	config := Floats{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, float64(1000000.1888), config.Big)
	require.Equal(t, float32(100000.321), config.Medium)
}

func TestFloatTypesBadValue(t *testing.T) {
	configString := `
	Big: "a"
	Medium: "b"
	`

	config := Floats{}

	err := LoadConfigFromString(configString, &config, false)
	require.Error(t, err)
}

func TestLowerCase(t *testing.T) {
	configString := `
	name: "stephen"
	age: 28
	optout: true
	balance: 5.5
	`

	config := Simple{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, "stephen", config.Name)
	require.Equal(t, int64(28), config.Age)
	require.Equal(t, true, config.OptOut)
	require.Equal(t, 5.5, config.Balance)
}

func TestEqualSignAndSpace(t *testing.T) {
	configString := `
	name: "stephen"
	age = 28
	optout true
	balance: 5.5
	`

	config := Simple{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, "stephen", config.Name)
	require.Equal(t, int64(28), config.Age)
	require.Equal(t, true, config.OptOut)
	require.Equal(t, 5.5, config.Balance)
}

type PrimitiveArrays struct {
	Ints     []int
	Int8s    []int8
	Int16s   []int16
	Int32s   []int32
	Int64s   []int64
	Float32s []float32
	Float64s []float64
	Strings  []string
}

func TestArrays(t *testing.T) {
	configString := `
	strings: [
		"mister",
		"zero"
	],
	ints: [10, 15, -1],
	int8s: [10, 15, -1],
	int16s: [10, 15, -1],
	int32s: [10, 15, -1],
	int64s: [10, 15, -1],
	float32s: [1.1, 2.2, 3.3],
	float64s: [1.1, 2.2, 3.3],
	`

	config := PrimitiveArrays{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"mister", "zero"}, config.Strings)
	require.ElementsMatch(t, []int{10, 15, -1}, config.Ints)
	require.ElementsMatch(t, []int8{10, 15, -1}, config.Int8s)
	require.ElementsMatch(t, []int16{10, 15, -1}, config.Int16s)
	require.ElementsMatch(t, []int32{10, 15, -1}, config.Int32s)
	require.ElementsMatch(t, []int64{10, 15, -1}, config.Int64s)
	require.ElementsMatch(t, []float32{1.1, 2.2, 3.3}, config.Float32s)
	require.ElementsMatch(t, []float64{1.1, 2.2, 3.3}, config.Float64s)
}

func TestArraysSingleton(t *testing.T) {
	configString := `
	strings: "mister",
	ints: 10,
	int8s: 10,
	int16s: 10,
	int32s: 10,
	int64s: 10,
	float32s: 1.2,
	float64s: 1.3,
	`

	config := PrimitiveArrays{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"mister"}, config.Strings)
	require.ElementsMatch(t, []int{10}, config.Ints)
	require.ElementsMatch(t, []int8{10}, config.Int8s)
	require.ElementsMatch(t, []int16{10}, config.Int16s)
	require.ElementsMatch(t, []int32{10}, config.Int32s)
	require.ElementsMatch(t, []int64{10}, config.Int64s)
	require.ElementsMatch(t, []float32{1.2}, config.Float32s)
	require.ElementsMatch(t, []float64{1.3}, config.Float64s)
}

func TestArraysSingletonBadValue(t *testing.T) {
	configString := `
	ints: hello,
	`

	config := PrimitiveArrays{}
	err := LoadConfigFromString(configString, &config, false)
	require.Error(t, err)
}

func TestArraysBadValue(t *testing.T) {
	configString := `
	strings: 43a
	`

	config := PrimitiveArrays{}

	err := LoadConfigFromString(configString, &config, false)
	require.Error(t, err)

	configString = `
	ints: 43a
	`
	config = PrimitiveArrays{}
	err = LoadConfigFromString(configString, &config, false)
	require.Error(t, err)
}

type Person struct {
	Name string
}

type Parent struct {
	Name  string
	Child Person
}

type GrandParent struct {
	Name     string
	Child    Parent
	Children []Parent
}

func TestStructField(t *testing.T) {
	configString := `
	Name: "stephen"
	Child: {
		Name: "mister",
		Child: {
			Name: "zero",
		},
	}
	`

	config := GrandParent{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, "stephen", config.Name)
	require.Equal(t, "mister", config.Child.Name)
	require.Equal(t, "zero", config.Child.Child.Name)
}

func TestStructFieldBadValue(t *testing.T) {
	configString := `
	Name: "stephen"
	Child: "string"
	`

	config := GrandParent{}

	err := LoadConfigFromString(configString, &config, false)
	require.Error(t, err)
}

func TestStructFieldWithDefault(t *testing.T) {
	configString := `
	Name: "stephen"
	Child: {
		Child: {
			Name: "zero",
		},
	}
	`

	config := GrandParent{
		Child: Parent{
			Name: "mister",
		},
	}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, "stephen", config.Name)
	require.Equal(t, "mister", config.Child.Name)
	require.Equal(t, "zero", config.Child.Child.Name)
}

func TestArrayStructs(t *testing.T) {
	configString := `
	Name: "stephen"
	Children: [{Name: "mister", Child:{Name: "zero"}}, {Name:"alpha", Child:{Name: "omega"}}]
	`

	config := GrandParent{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, "stephen", config.Name)
	require.Equal(t, 2, len(config.Children))
	require.Equal(t, "mister", config.Children[0].Name)
	require.Equal(t, "zero", config.Children[0].Child.Name)
	require.Equal(t, "alpha", config.Children[1].Name)
	require.Equal(t, "omega", config.Children[1].Child.Name)
}

func TestArrayStructsSingleton(t *testing.T) {
	configString := `
	Name: "stephen"
	Children: {Name: "mister", Child:{Name: "zero"}}
	`

	config := GrandParent{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, "stephen", config.Name)
	require.Equal(t, 1, len(config.Children))
	require.Equal(t, "mister", config.Children[0].Name)
	require.Equal(t, "zero", config.Children[0].Child.Name)
}

func TestArrayStructsBadArray(t *testing.T) {
	configString := `
	Name: "stephen"
	Children: "string"
	`

	config := GrandParent{}

	err := LoadConfigFromString(configString, &config, false)
	require.Error(t, err)
}

type Server struct {
	Name      string
	Interface HostPort
}

func TestHostPortField(t *testing.T) {
	configString := `
	Name: "stephen"
	Interface: 8080
	`

	config := Server{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, "stephen", config.Name)
	require.Equal(t, 8080, config.Interface.Port)
	require.Equal(t, "", config.Interface.Host)

	configString = `
	Name: "stephen"
	Interface: localhost:8080
	`

	config = Server{}

	err = LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, "stephen", config.Name)
	require.Equal(t, 8080, config.Interface.Port)
	require.Equal(t, "localhost", config.Interface.Host)
}

func TestBadHostPort(t *testing.T) {
	configString := `
	Interface: localhost
	`

	config := Server{}

	err := LoadConfigFromString(configString, &config, false)
	require.Error(t, err)

	configString = `
	Interface: localhost:hello
	`

	config = Server{}

	err = LoadConfigFromString(configString, &config, false)
	require.Error(t, err)
}

type Service struct {
	Name       string
	Interfaces []HostPort
}

func TestHostPortArray(t *testing.T) {
	configString := `
	Name: "stephen"
	Interfaces: [8080, "localhost:9090", "192.168.1.1:20"]
	`

	config := Service{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, "stephen", config.Name)
	require.Equal(t, 3, len(config.Interfaces))
	require.Equal(t, 8080, config.Interfaces[0].Port)
	require.Equal(t, "", config.Interfaces[0].Host)
	require.Equal(t, 9090, config.Interfaces[1].Port)
	require.Equal(t, "localhost", config.Interfaces[1].Host)
	require.Equal(t, 20, config.Interfaces[2].Port)
	require.Equal(t, "192.168.1.1", config.Interfaces[2].Host)
}

func TestHostPortArraySingleton(t *testing.T) {
	configString := `
	Name: "stephen"
	Interfaces: "localhost:9090"
	`

	config := Service{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, "stephen", config.Name)
	require.Equal(t, 1, len(config.Interfaces))
	require.Equal(t, 9090, config.Interfaces[0].Port)
	require.Equal(t, "localhost", config.Interfaces[0].Host)
}

func TestBadArray(t *testing.T) {
	configString := `
	Name: "stephen"
	Interfaces: false
	`
	config := Service{}
	err := LoadConfigFromString(configString, &config, false)
	require.Error(t, err)
}

type ComplexHolder struct {
	Name string
	Comp complex64
}

func TestComplexFailsOnStrict(t *testing.T) {
	configString := `
	Name: "stephen"
	Comp: 34
	`

	config := ComplexHolder{}

	err := LoadConfigFromString(configString, &config, true)
	require.Error(t, err)
}

type BoolArray struct {
	Name     string
	TheArray []bool
}

func TestBoolArrayNotSupported(t *testing.T) {
	configString := `
	Name: "stephen"
	TheArray: [true, false, true]
	`

	config := BoolArray{}

	err := LoadConfigFromString(configString, &config, true)
	require.Error(t, err)
}

func TestEmptyConfig(t *testing.T) {
	configString := ``

	config := Simple{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)

	err = LoadConfigFromString(configString, &config, true)
	require.Error(t, err)
}

type HasPrivate struct {
	name string
}

func TestPrivateField(t *testing.T) {
	configString := `
	name: "stephen"
	`

	config := HasPrivate{
		name: "alberto",
	}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)

	err = LoadConfigFromString(configString, &config, true)
	require.Error(t, err)
}

type SimpleMapConf struct {
	One map[string]interface{}
}

type StringStringMapConf struct {
	One map[string]string
}

func TestMap(t *testing.T) {
	configString := `
	One: {
	Name: "stephen"
	}
	`

	config := SimpleMapConf{}

	err := LoadConfigFromString(configString, &config, false)
	require.NoError(t, err)
	require.Equal(t, config.One["Name"], "stephen")
}

func TestBadMapField(t *testing.T) {
	configString := `
	One: "test"
	`

	config := SimpleMapConf{}

	err := LoadConfigFromString(configString, &config, false)
	require.Error(t, err)
}

func TestBadMap(t *testing.T) {
	configString := `
	One: {
	Name: "stephen"
	Age: 28
	OptOut: true
	Balance: 5.5
	}
	`

	config := StringStringMapConf{}

	err := LoadConfigFromString(configString, &config, false)
	require.Error(t, err)
}
