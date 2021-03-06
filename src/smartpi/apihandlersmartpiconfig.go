/*
    Copyright (C) Jens Ramhorst
  	This file is part of SmartPi.
    SmartPi is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.
    SmartPi is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.
    You should have received a copy of the GNU General Public License
    along with SmartPi.  If not, see <http://www.gnu.org/licenses/>.
    Diese Datei ist Teil von SmartPi.
    SmartPi ist Freie Software: Sie können es unter den Bedingungen
    der GNU General Public License, wie von der Free Software Foundation,
    Version 3 der Lizenz oder (nach Ihrer Wahl) jeder späteren
    veröffentlichten Version, weiterverbreiten und/oder modifizieren.
    SmartPi wird in der Hoffnung, dass es nützlich sein wird, aber
    OHNE JEDE GEWÄHRLEISTUNG, bereitgestellt; sogar ohne die implizite
    Gewährleistung der MARKTFÄHIGKEIT oder EIGNUNG FÜR EINEN BESTIMMTEN ZWECK.
    Siehe die GNU General Public License für weitere Details.
    Sie sollten eine Kopie der GNU General Public License zusammen mit diesem
    Programm erhalten haben. Wenn nicht, siehe <http://www.gnu.org/licenses/>.
*/
/*
File: apihandlersmomentary.go
Description: Handels API requests
*/

package smartpi

import (
	"fmt"
	// "github.com/gorilla/mux"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/fatih/structs"
	"github.com/gorilla/context"
	"github.com/oleiade/reflections"
)

type writeconfiguration struct {
	Type string
	Msg  interface{}
}

func ReadConfig(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// name := vars["name"]

	// user := context.Get(r,"Username")
	configuration := context.Get(r, "Config")
	if err := json.NewEncoder(w).Encode(configuration.(*Config)); err != nil {
		panic(err)
	}

	// if configuration := r.Context().Value("Config"); configuration != nil {
	// 	if err := json.NewEncoder(w).Encode(configuration.(*Config)); err != nil {
	// 		panic(err)
	// 	}
	// }
}

func WriteConfig(w http.ResponseWriter, r *http.Request) {
	var wc writeconfiguration

	b, _ := ioutil.ReadAll(r.Body)

	if err := json.Unmarshal(b, &wc); err != nil {
		log.Fatal(err)
	}

	configuration := context.Get(r, "Config")
	// if configuration := r.Context().Value("Config"); configuration != nil {

	keys := make([]string, 0, len(wc.Msg.(map[string]interface{})))
	for k := range wc.Msg.(map[string]interface{}) {
		keys = append(keys, k)
	}
	fmt.Printf("%+v\n", keys)

	confignames := structs.Names(configuration.(*Config))

	for i := range confignames {
		for j := range keys {
			if keys[j] == confignames[i] {

				// fmt.Println("Treffer: Key: " + keys[j] + " Configname: " + confignames[i])
				// fmt.Println(reflect.TypeOf(wc.Msg.(map[string]interface{})[keys[j]]))
				// fmt.Println(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]))

				var err error
				var fieldtype string
				fieldtype, err = reflections.GetFieldType(configuration.(*Config), confignames[i])
				// fmt.Println("Fieldtype: " + fieldtype)

				switch fieldtype {
				case "int":
					switch wc.Msg.(map[string]interface{})[keys[j]].(type) {
					case float64:
						err = reflections.SetField(configuration.(*Config), confignames[i], int(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Float()))
					case string:
						intval, _ := strconv.Atoi(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).String())
						err = reflections.SetField(configuration.(*Config), confignames[i], intval)
					case int:
						err = reflections.SetField(configuration.(*Config), confignames[i], int(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Int()))
					case bool:
						err = reflections.SetField(configuration.(*Config), confignames[i], b2i(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Bool()))
					}
				case "float64":
					switch wc.Msg.(map[string]interface{})[keys[j]].(type) {
					case float64:
						err = reflections.SetField(configuration.(*Config), confignames[i], float64(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Float()))
					case string:
						floatval, _ := strconv.ParseFloat(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).String(), 64)
						err = reflections.SetField(configuration.(*Config), confignames[i], floatval)
					case int:
						err = reflections.SetField(configuration.(*Config), confignames[i], float64(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Int()))
					case bool:
						err = reflections.SetField(configuration.(*Config), confignames[i], float64(b2i(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Bool())))
					}
				case "string":
					switch wc.Msg.(map[string]interface{})[keys[j]].(type) {
					case float64:
						err = reflections.SetField(configuration.(*Config), confignames[i], strconv.FormatFloat(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Float(), 'f', -1, 64))
					case string:
						err = reflections.SetField(configuration.(*Config), confignames[i], reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).String())
					case int:
						err = reflections.SetField(configuration.(*Config), confignames[i], strconv.FormatInt(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Int(), 16))
					case bool:
						err = reflections.SetField(configuration.(*Config), confignames[i], strconv.FormatBool(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Bool()))
					}
				case "bool":
					switch wc.Msg.(map[string]interface{})[keys[j]].(type) {
					case float64:
						err = reflections.SetField(configuration.(*Config), confignames[i], !(int(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Float()) == 0))
					case string:
						boolval, _ := strconv.ParseBool(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).String())
						err = reflections.SetField(configuration.(*Config), confignames[i], boolval)
					case int:
						err = reflections.SetField(configuration.(*Config), confignames[i], !(int(reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Int()) == 0))
					case bool:
						err = reflections.SetField(configuration.(*Config), confignames[i], reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Bool())
					}
				case "map[string]int":
					values := reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Interface().(map[string]interface{})
					aData := make(map[string]int)
					for k, v := range values {
						// fmt.Println(reflect.TypeOf(v))
						switch v.(type) {
						case float64:
							aData[k] = int(v.(float64))
						case string:
							intval, _ := strconv.Atoi(v.(string))
							aData[k] = intval
						case int:
							aData[k] = int(v.(int))
						case bool:
							aData[k] = b2i(v.(bool))
						}
						// fmt.Printf("key[%s] value[%s]\n", k, v)
					}
					err = reflections.SetField(configuration.(*Config), confignames[i], aData)
				case "map[string]float":
					values := reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Interface().(map[string]interface{})
					aData := make(map[string]float64)
					for k, v := range values {
						// fmt.Println(reflect.TypeOf(v))
						switch v.(type) {
						case float64:
							aData[k] = float64(v.(float64))
						case string:
							floatval, _ := strconv.ParseFloat(v.(string), 64)
							aData[k] = floatval
						case int:
							aData[k] = float64(v.(int))
						case bool:
							aData[k] = float64(b2i(v.(bool)))
						}
						// fmt.Printf("key[%s] value[%s]\n", k, v)
					}
					err = reflections.SetField(configuration.(*Config), confignames[i], aData)
				case "map[string]string":
					values := reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Interface().(map[string]interface{})
					aData := make(map[string]string)
					for k, v := range values {
						// fmt.Println(reflect.TypeOf(v))
						switch v.(type) {
						case float64:
							aData[k] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
						case string:
							aData[k] = v.(string)
						case int:
							aData[k] = strconv.FormatInt(v.(int64), 16)
						case bool:
							aData[k] = strconv.FormatBool(v.(bool))
						}
						// fmt.Printf("key[%s] value[%s]\n", k, v)
					}
					err = reflections.SetField(configuration.(*Config), confignames[i], aData)
				case "map[string]bool":
					values := reflect.ValueOf(wc.Msg.(map[string]interface{})[keys[j]]).Interface().(map[string]interface{})
					aData := make(map[string]bool)
					for k, v := range values {
						// fmt.Println(reflect.TypeOf(v))
						switch v.(type) {
						case float64:
							aData[k] = !(int(v.(float64)) == 0)
						case string:
							boolval, _ := strconv.ParseBool(v.(string))
							aData[k] = boolval
						case int:
							aData[k] = !(int(v.(int)) == 0)
						case bool:
							aData[k] = v.(bool)
						}
						// fmt.Printf("key[%s] value[%s]\n", k, v)
					}
					err = reflections.SetField(configuration.(*Config), confignames[i], aData)
				}
				if err != nil {
					log.Fatal(err)
				}

			}
		}
	}
	fmt.Printf("%+v\n", configuration.(*Config))
	configuration.(*Config).SaveParameterToFile()
	// }
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func GetStringValueByFieldName(n interface{}, field_name string) (string, bool) {
	s := reflect.ValueOf(n)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	if s.Kind() != reflect.Struct {
		return "", false
	}
	f := s.FieldByName(field_name)
	if !f.IsValid() {
		return "", false
	}
	switch f.Kind() {
	case reflect.String:
		return f.Interface().(string), true
	case reflect.Int:
		return strconv.FormatInt(f.Int(), 10), true
	// add cases for more kinds as needed.
	default:
		return "", false
		// or use fmt.Sprint(f.Interface())
	}
}
