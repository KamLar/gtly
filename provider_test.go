package gtly_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/gtly"
	"github.com/viant/gtly/codec/json"
	"log"
	"reflect"
	"testing"
	"time"
)

//ExampleProvider_NewObject new object example
func ExampleProvider_NewObject() {

	fooProvider := gtly.NewProvider("foo",
		gtly.NewField("id", gtly.FieldTypeInt),
		gtly.NewField("firsName", gtly.FieldTypeString),
		gtly.NewField("description", gtly.FieldTypeString, gtly.OmitEmptyOpt(true)),
		gtly.NewField("updated", gtly.FieldTypeTime, gtly.DateLayoutOpt("2006-01-02T15:04:05Z07:00")),
		gtly.NewField("numbers", gtly.FieldTypeArray, gtly.ComponentTypeOpt(gtly.FieldTypeInt)),
	)
	foo1 := fooProvider.NewObject()
	foo1.SetValue("id", 1)
	foo1.SetValue("firsName", "Adam")
	foo1.SetValue("updated", time.Now())
	foo1.SetValue("numbers", []int{1, 2, 3})

	JSON, err := json.Marshal(foo1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", JSON)
	fooProvider.OutputCaseFormat(gtly.CaseLowerCamel, gtly.CaseUpperUnderscore)
	JSON, _ = json.Marshal(foo1)
	fmt.Printf("%s\n", JSON)

	fooProvider.OutputCaseFormat(gtly.CaseLowerCamel, gtly.CaseLowerUnderscore)
	foo1.SetValue("active", true)
	foo1.SetValue("description", "some description")

	JSON, _ = json.Marshal(foo1)
	fmt.Printf("%s\n", JSON)
}

//TestProvider_NewArray new array example
func ExampleProvider_NewArray() {
	fooProvider := gtly.NewProvider("foo",
		gtly.NewField("id", gtly.FieldTypeInt),
		gtly.NewField("firsName", gtly.FieldTypeString),
		gtly.NewField("income", gtly.FieldTypeFloat),
		gtly.NewField("description", gtly.FieldTypeString, gtly.OmitEmptyOpt(true)),
		gtly.NewField("updated", gtly.FieldTypeTime, gtly.DateLayoutOpt(time.RFC3339)),
		gtly.NewField("numbers", gtly.FieldTypeArray, gtly.ComponentTypeOpt(gtly.FieldTypeInt)),
	)
	fooArray1 := fooProvider.NewArray()

	for i := 0; i < 10; i++ {
		foo1 := fooProvider.NewObject()
		foo1.SetValue("id", 1)
		foo1.SetValue("income", 64000.0*float64(1+(10/(i+1))))
		foo1.SetValue("firsName", "Adam")
		fooArray1.AddObject(foo1)
	}

	now := time.Now()
	fooArray1.Add(map[string]interface{}{
		"id":       100,
		"firsName": "Tom",
		"updated":  now,
	})

	totalIncome := 0.0
	incomeField := fooProvider.Field("income")
	//Iterating collection
	err := fooArray1.Objects(func(object *gtly.Object) (bool, error) {
		fmt.Printf("id: %v\n", object.Value("id"))
		fmt.Printf("name: %v\n", object.Value("name"))
		value, _ := object.FloatAt(incomeField.Index)
		totalIncome += value
		return true, nil
	})
	fmt.Printf("income total: %v\n", totalIncome)
	if err != nil {
		log.Fatal(err)
	}
	JSON, err := json.Marshal(fooArray1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", JSON)
}

//ExampleProvider_NewMap new map example
func ExampleProvider_NewMap() {

	fooProvider := gtly.NewProvider("foo",
		gtly.NewField("id", gtly.FieldTypeInt),
		gtly.NewField("firsName", gtly.FieldTypeString),
		gtly.NewField("description", gtly.FieldTypeString, gtly.OmitEmptyOpt(true)),
		gtly.NewField("updated", gtly.FieldTypeTime, gtly.DateLayoutOpt(time.RFC3339)),
		gtly.NewField("numbers", gtly.FieldTypeArray, gtly.ComponentTypeOpt(gtly.FieldTypeInt)),
	)

	//Creates a map keyed by id field
	aMap := fooProvider.NewMap(gtly.NewKeyProvider("id"))
	for i := 0; i < 10; i++ {
		foo := fooProvider.NewObject()
		foo.SetValue("id", i)
		foo.SetValue("firsName", fmt.Sprintf("Name %v", i))
		aMap.AddObject(foo)
	}

	//Accessing map
	foo1 := aMap.Object("1")
	JSON, err := json.Marshal(foo1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", JSON)

	//Iterating map
	err = aMap.Pairs(func(key interface{}, item *gtly.Object) (bool, error) {
		fmt.Printf("id: %v\n", item.Value("id"))
		fmt.Printf("name: %v\n", item.Value("name"))
		return true, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	JSON, err = json.Marshal(aMap)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", JSON)
	//[{"id":1,"firsName":"Name 1","updated":null,"numbers":null},{"id":3,"firsName":"Name 3","updated":null,"numbers":null},{"id":4,"firsName":"Name 4","updated":null,"numbers":null},{"id":6,"firsName":"Name 6","updated":null,"numbers":null},{"id":8,"firsName":"Name 8","updated":null,"numbers":null},{"id":9,"firsName":"Name 9","updated":null,"numbers":null},{"id":0,"firsName":"Name 0","updated":null,"numbers":null},{"id":2,"firsName":"Name 2","updated":null,"numbers":null},{"id":5,"firsName":"Name 5","updated":null,"numbers":null},{"id":7,"firsName":"Name 7","updated":null,"numbers":null}]

}

//ExampleProvider_NewMulti new multi example
func ExampleProvider_NewMultimap() {
	fooProvider := gtly.NewProvider("foo",
		gtly.NewField("id", gtly.FieldTypeInt),
		gtly.NewField("firsName", gtly.FieldTypeString),
		gtly.NewField("city", gtly.FieldTypeString, gtly.OmitEmptyOpt(true)),
		gtly.NewField("updated", gtly.FieldTypeTime, gtly.DateLayoutOpt(time.RFC3339)),
		gtly.NewField("numbers", gtly.FieldTypeArray, gtly.ComponentTypeOpt(gtly.FieldTypeInt)),
	)

	//Creates a multi map keyed by id field
	aMap := fooProvider.NewMultimap(gtly.NewKeyProvider("city"))
	for i := 0; i < 10; i++ {
		foo := fooProvider.NewObject()
		foo.SetValue("id", i)
		foo.SetValue("firsName", fmt.Sprintf("Name %v", i))
		if i%2 == 0 {
			foo.SetValue("city", "Cracow")
		} else {
			foo.SetValue("city", "Warsaw")
		}
		aMap.AddObject(foo)
	}

	//Accessing map
	fooInWarsawSlice := aMap.Slice("Warsaw")
	JSON, err := json.Marshal(fooInWarsawSlice)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", JSON)
	//Prints [{"id":1,"firsName":"Name 1","city":"Warsaw","updated":null,"numbers":null},{"id":3,"firsName":"Name 3","city":"Warsaw","updated":null,"numbers":null},{"id":5,"firsName":"Name 5","city":"Warsaw","updated":null,"numbers":null},{"id":7,"firsName":"Name 7","city":"Warsaw","updated":null,"numbers":null},{"id":9,"firsName":"Name 9","city":"Warsaw","updated":null,"numbers":null}]

	//Iterating multi map
	err = aMap.Slices(func(key interface{}, value *gtly.Array) (bool, error) {
		fmt.Printf("%v -> %v\n", key, value.Size())
		return true, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	JSON, err = json.Marshal(aMap)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", JSON)
	//[{"id":0,"firsName":"Name 0","city":"Cracow","updated":null,"numbers":null},{"id":2,"firsName":"Name 2","city":"Cracow","updated":null,"numbers":null},{"id":4,"firsName":"Name 4","city":"Cracow","updated":null,"numbers":null},{"id":6,"firsName":"Name 6","city":"Cracow","updated":null,"numbers":null},{"id":8,"firsName":"Name 8","city":"Cracow","updated":null,"numbers":null},{"id":1,"firsName":"Name 1","city":"Warsaw","updated":null,"numbers":null},{"id":3,"firsName":"Name 3","city":"Warsaw","updated":null,"numbers":null},{"id":5,"firsName":"Name 5","city":"Warsaw","updated":null,"numbers":null},{"id":7,"firsName":"Name 7","city":"Warsaw","updated":null,"numbers":null},{"id":9,"firsName":"Name 9","city":"Warsaw","updated":null,"numbers":null}]

}

func TestProvider_UnMarshall(t *testing.T) {
	testCases := []struct {
		description    string
		json           []byte
		fields         map[string]interface{}
		expectedValues map[string]interface{}
	}{
		{
			description: "unmarshalling",
			fields: map[string]interface{}{
				"Id":    0,
				"Name":  "",
				"Price": 0.5,
			},
			json: []byte(`{
				"Id": 1,
				"Name": "Foo"
			}`),
			expectedValues: map[string]interface{}{
				"Id":   1,
				"Name": "Foo",
			},
		},
	}

	for i, testCase := range testCases {
		provider := gtly.NewProvider(fmt.Sprintf("testcase#%v", i))
		addProviderFields(testCase, provider)
		anObject := provider.UnMarshall(testCase.json)
		for _, field := range anObject.Proto().Fields() {
			value, ok := testCase.expectedValues[field.Name]
			if ok {
				assert.True(t, anObject.SetAt(field.Index), testCase.description)
				assert.Equal(t, value, anObject.Value(field.Name), testCase.description)
			} else {
				assert.False(t, anObject.SetAt(field.Index), testCase.description)
			}
		}
	}
}

func addProviderFields(testCase struct {
	description    string
	json           []byte
	fields         map[string]interface{}
	expectedValues map[string]interface{}
}, provider *gtly.Provider) {
	for k, v := range testCase.fields {
		provider.AddField(&gtly.Field{
			Name: k,
			Type: reflect.TypeOf(v),
		})
	}
}
