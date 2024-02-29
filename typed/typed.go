// A wrapper to make map[string]any a little more type-safe
package typed

import (
	"errors"
	"io/ioutil"
	"strconv"
	"time"

	"src.goblgobl.com/utils/json"
	"src.goblgobl.com/utils/optional"
)

var (
	// Used by ToBytes to indicate that the key was not
	// present in the type
	KeyNotFound = errors.New("Key not found")
	Empty       = Typed(nil)
)

// A Typed type helper for accessing a map
type Typed map[string]any

// Wrap the map into a Typed
func New(m map[string]any) Typed {
	return Typed(m)
}

// Create a Typed helper from the given JSON bytes
func Json(data []byte) (Typed, error) {
	if data == nil {
		return Typed{}, nil
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return Typed(m), nil
}

// Create an array of Typed helpers
// Used for when the root is an array which contains objects
func JsonArray(data []byte) ([]Typed, error) {
	var m []interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	l := len(m)
	if l == 0 {
		return nil, nil
	}

	typed := make([]Typed, l)
	for i := 0; i < l; i++ {
		value := m[i]
		if t, ok := value.(map[string]interface{}); ok {
			typed[i] = t
		} else {
			typed[i] = map[string]interface{}{"0": value}
		}
	}
	return typed, nil
}

// Create a Typed helper from the given JSON bytes, panics on error
func Must(data []byte) Typed {
	m, err := Json(data)
	if err != nil {
		panic(err)
	}
	return m
}

// Create a Typed helper from the given JSON string
func JsonString(data string) (Typed, error) {
	return Json([]byte(data))
}

// Create a Typed helper from the JSON within a file
func JsonFile(path string) (Typed, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Json(data)
}

func (t Typed) Keys() []string {
	keys := make([]string, len(t))
	i := 0
	for k := range t {
		keys[i] = k
		i++
	}
	return keys
}

// Returns a boolean at the key, or false if it
// doesn't exist, or if it isn't a bool
func (t Typed) Bool(key string) bool {
	return t.BoolOr(key, false)
}

func (t Typed) OptionalBool(key string) optional.Bool {
	if value, exists := t.BoolIf(key); exists {
		return optional.NewBool(value)
	}
	return optional.NullBool
}

// Returns a boolean at the key, or the specified
// value if it doesn't exist or isn't a bool
func (t Typed) BoolOr(key string, d bool) bool {
	if value, exists := t.BoolIf(key); exists {
		return value
	}
	return d
}

// Returns a bool or panics
func (t Typed) BoolMust(key string) bool {
	b, exists := t.BoolIf(key)
	if exists == false {
		panic("expected boolean value for " + key)
	}
	return b
}

// Returns a boolean at the key and whether
// or not the key existed and the value was a bolean
func (t Typed) BoolIf(key string) (bool, bool) {
	value, exists := t[key]
	if exists == false {
		return false, false
	}
	switch t := value.(type) {
	case bool:
		return t, true
	case int:
		switch t {
		case 1:
			return true, true
		case 0:
			return false, true
		}
	case string:
		switch t {
		case "true", "TRUE", "True":
			return true, true
		case "false", "FALSE", "False":
			return false, true
		}
	}
	return false, false
}

func (t Typed) Int(key string) int {
	return t.IntOr(key, 0)
}

func (t Typed) OptionalInt(key string) optional.Int {
	if value, exists := t.IntIf(key); exists {
		return optional.NewInt(value)
	}
	return optional.NullInt
}

// Returns a int at the key, or the specified
// value if it doesn't exist or isn't a int
func (t Typed) IntOr(key string, d int) int {
	if value, exists := t.IntIf(key); exists {
		return value
	}
	return d
}

// Returns an int or panics
func (t Typed) IntMust(key string) int {
	i, exists := t.IntIf(key)
	if exists == false {
		panic("expected int value for " + key)
	}
	return i
}

// Returns an int at the key and whether
// or not the key existed and the value was an int
func (t Typed) IntIf(key string) (int, bool) {
	value, exists := t[key]
	if exists == false {
		return 0, false
	}
	return NumericToInt(value)
}

func (t Typed) Float(key string) float64 {
	return t.FloatOr(key, 0)
}

func (t Typed) OptionalFloat(key string) optional.Float {
	if value, exists := t.FloatIf(key); exists {
		return optional.NewFloat(value)
	}
	return optional.NullFloat
}

// Returns a float at the key, or the specified
// value if it doesn't exist or isn't a float
func (t Typed) FloatOr(key string, d float64) float64 {
	if value, exists := t.FloatIf(key); exists {
		return value
	}
	return d
}

// Returns an float or panics
func (t Typed) FloatMust(key string) float64 {
	f, exists := t.FloatIf(key)
	if exists == false {
		panic("expected float value for " + key)
	}
	return f
}

// Returns an float at the key and whether
// or not the key existed and the value was an float
func (t Typed) FloatIf(key string) (float64, bool) {
	value, exists := t[key]
	if exists == false {
		return 0, false
	}
	switch t := value.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	case string:
		f, err := strconv.ParseFloat(t, 10)
		return f, err == nil
	}
	return 0, false
}

func (t Typed) String(key string) string {
	return t.StringOr(key, "")
}

func (t Typed) OptionalString(key string) optional.String {
	if value, exists := t.StringIf(key); exists {
		return optional.NewString(value)
	}
	return optional.NullString
}

// Returns a string at the key, or the specified
// value if it doesn't exist or isn't a string
func (t Typed) StringOr(key string, d string) string {
	if value, exists := t.StringIf(key); exists {
		return value
	}
	return d
}

// Returns an string or panics
func (t Typed) StringMust(key string) string {
	s, exists := t.StringIf(key)
	if exists == false {
		panic("expected string value for " + key)
	}
	return s
}

// Returns an string at the key and whether
// or not the key existed and the value was an string
func (t Typed) StringIf(key string) (string, bool) {
	value, exists := t[key]
	if exists == false {
		return "", false
	}
	if n, ok := value.(string); ok {
		return n, true
	}
	return "", false
}

func (t Typed) Bytes(key string) []byte {
	return t.BytesOr(key, nil)
}

// Returns a []byte at the key, or the specified
// value if it doesn't exist or isn't a []byte
func (t Typed) BytesOr(key string, d []byte) []byte {
	if value, exists := t.BytesIf(key); exists {
		return value
	}
	return d
}

// Returns an []byte or panics
func (t Typed) BytesMust(key string) []byte {
	s, exists := t.BytesIf(key)
	if exists == false {
		panic("expected []byte value for " + key)
	}
	return s
}

// Returns an []byte at the key and whether
// or not the key existed and the value was an string
func (t Typed) BytesIf(key string) ([]byte, bool) {
	value, exists := t[key]
	if exists == false {
		return nil, false
	}
	switch t := value.(type) {
	case []byte:
		return t, true
	case string:
		return []byte(t), true
	}
	return nil, false
}

func (t Typed) Time(key string) time.Time {
	return t.TimeOr(key, time.Time{})
}

// Returns a time at the key, or the specified
// value if it doesn't exist or isn't a time
func (t Typed) TimeOr(key string, d time.Time) time.Time {
	if value, exists := t.TimeIf(key); exists {
		return value
	}
	return d
}

// Returns a time.Time or panics
func (t Typed) TimeMust(key string) time.Time {
	tt, exists := t.TimeIf(key)
	if exists == false {
		panic("expected time.Time value for " + key)
	}
	return tt
}

// Returns an time.time at the key and whether
// or not the key existed and the value was a time.Time
func (t Typed) TimeIf(key string) (time.Time, bool) {
	value, exists := t[key]
	if exists == false {
		return time.Time{}, false
	}
	switch n := value.(type) {
	case time.Time:
		return n, true
	case string:
		if t, err := time.Parse(time.RFC3339, n); err != nil {
			return time.Time{}, false
		} else {
			return t, true
		}
	case int:
		return time.Unix(int64(n), 0).UTC(), true
	case int64:
		return time.Unix(n, 0).UTC(), true
	}
	return time.Time{}, false
}

// Returns a Typed helper at the key
// If the key doesn't exist, a default Typed helper
// is returned (which will return default values for
// any subsequent sub queries)
func (t Typed) Object(key string) Typed {
	o := t.ObjectOr(key, nil)
	if o == nil {
		return Typed(nil)
	}
	return o
}

// Returns a Typed helper at the key or the specified
// default if the key doesn't exist or if the key isn't
// a map[string]any
func (t Typed) ObjectOr(key string, d map[string]any) Typed {
	if value, exists := t.ObjectIf(key); exists {
		return value
	}
	return Typed(d)
}

// Returns an typed object or panics
func (t Typed) ObjectMust(key string) Typed {
	t, exists := t.ObjectIf(key)
	if exists == false {
		panic("expected map for " + key)
	}
	return t
}

// Returns a Typed helper at the key and whether
// or not the key existed and the value was an map[string]any
func (t Typed) ObjectIf(key string) (Typed, bool) {
	value, exists := t[key]
	if exists == false {
		return nil, false
	}
	switch t := value.(type) {
	case map[string]any:
		return Typed(t), true
	case Typed:
		return t, true
	}
	return nil, false
}

func (t Typed) Any(key string) any {
	return t.AnyOr(key, nil)
}

// Returns a string at the key, or the specified
// value if it doesn't exist or isn't a strin
func (t Typed) AnyOr(key string, d any) any {
	if value, exists := t.AnyIf(key); exists {
		return value
	}
	return d
}

// Returns an interface or panics
func (t Typed) AnyMust(key string) any {
	i, exists := t.AnyIf(key)
	if exists == false {
		panic("expected any for " + key)
	}
	return i
}

// Returns an string at the key and whether
// or not the key existed and the value was an string
func (t Typed) AnyIf(key string) (any, bool) {
	value, exists := t[key]
	if exists == false {
		return nil, false
	}
	return value, true
}

func (t Typed) Anys(key string) []any {
	return t.AnysOr(key, nil)
}

// Returns a string at the key, or the specified
// value if it doesn't exist or isn't a strin
func (t Typed) AnysOr(key string, d []any) []any {
	if value, exists := t.AnysIf(key); exists {
		return value
	}
	return d
}

// Returns an interface or panics
func (t Typed) AnysMust(key string) []any {
	i, exists := t.AnysIf(key)
	if exists == false {
		panic("expected array of anys for " + key)
	}
	return i
}

// Returns an string at the key and whether
// or not the key existed and the value was an string
func (t Typed) AnysIf(key string) ([]any, bool) {
	value, exists := t[key].([]any)
	if exists == false {
		return nil, false
	}
	return value, true
}

// Returns a map[string]any at the key
// or a nil map if the key doesn't exist or if the key isn't
// a map[string]interface
func (t Typed) Map(key string) map[string]any {
	return t.MapOr(key, nil)
}

// Returns a map[string]any at the key
// or the specified default if the key doesn't exist
// or if the key isn't a map[string]interface
func (t Typed) MapOr(key string, d map[string]any) map[string]any {
	if value, exists := t.MapIf(key); exists {
		return value
	}
	return d
}

// Returns a map[string]interface at the key and whether
// or not the key existed and the value was an map[string]any
func (t Typed) MapIf(key string) (map[string]any, bool) {
	value, exists := t[key]
	if exists == false {
		return nil, false
	}
	if n, ok := value.(map[string]any); ok {
		return n, true
	}
	return nil, false
}

// Returns an slice of boolean, or an nil slice
func (t Typed) Bools(key string) []bool {
	return t.BoolsOr(key, nil)
}

// Returns an slice of boolean, or the specified slice
func (t Typed) BoolsOr(key string, d []bool) []bool {
	n, ok := t.BoolsIf(key)
	if ok {
		return n
	}
	return d
}

// Returns a boolean slice + true if valid
// Returns nil + false otherwise
// (returns nil+false if one of the values is not a valid boolean)
func (t Typed) BoolsIf(key string) ([]bool, bool) {
	value, exists := t[key]
	if exists == false {
		return nil, false
	}
	if n, ok := value.([]bool); ok {
		return n, true
	}
	if a, ok := value.([]any); ok {
		l := len(a)
		n := make([]bool, l)
		var ok bool
		for i := 0; i < l; i++ {
			if n[i], ok = a[i].(bool); ok == false {
				return nil, false
			}
		}
		return n, true
	}
	return nil, false
}

// Returns an slice of ints, or the specified slice
// Some conversion is done to handle the fact that JSON ints
// are represented as floats.
func (t Typed) Ints(key string) []int {
	return t.IntsOr(key, nil)
}

// Returns an slice of ints, or the specified slice
// if the key doesn't exist or isn't a valid []int.
// Some conversion is done to handle the fact that JSON ints
// are represented as floats.
func (t Typed) IntsOr(key string, d []int) []int {
	n, ok := t.IntsIf(key)
	if ok {
		return n
	}
	return d
}

// Returns a int slice + true if valid
// Returns nil + false otherwise
// (returns nil+false if one of the values is not a valid int)
func (t Typed) IntsIf(key string) ([]int, bool) {
	value, exists := t[key]
	if exists == false {
		return nil, false
	}
	if n, ok := value.([]int); ok {
		return n, true
	}
	if a, ok := value.([]any); ok {
		l := len(a)
		if l == 0 {
			return nil, false
		}

		n := make([]int, l)
		for i := 0; i < l; i++ {
			switch t := a[i].(type) {
			case int:
				n[i] = t
			case float64:
				n[i] = int(t)
			case string:
				_i, err := strconv.Atoi(t)
				if err != nil {
					return nil, false
				}
				n[i] = _i
			default:
				return nil, false
			}
		}
		return n, true
	}
	return nil, false
}

// Returns an slice of ints64, or the specified slice
// Some conversion is done to handle the fact that JSON ints
// are represented as floats.
func (t Typed) Ints64(key string) []int64 {
	return t.Ints64Or(key, nil)
}

// Returns an slice of ints, or the specified slice
// if the key doesn't exist or isn't a valid []int.
// Some conversion is done to handle the fact that JSON ints
// are represented as floats.
func (t Typed) Ints64Or(key string, d []int64) []int64 {
	n, ok := t.Ints64If(key)
	if ok {
		return n
	}
	return d
}

// Returns a boolean slice + true if valid
// Returns nil + false otherwise
// (returns nil+false if one of the values is not a valid boolean)
func (t Typed) Ints64If(key string) ([]int64, bool) {
	value, exists := t[key]
	if exists == false {
		return nil, false
	}
	if n, ok := value.([]int64); ok {
		return n, true
	}
	if a, ok := value.([]any); ok {
		l := len(a)
		if l == 0 {
			return nil, false
		}

		n := make([]int64, l)
		for i := 0; i < l; i++ {
			switch t := a[i].(type) {
			case int64:
				n[i] = t
			case float64:
				n[i] = int64(t)
			case int:
				n[i] = int64(t)
			case string:
				_i, err := strconv.ParseInt(t, 10, 10)
				if err != nil {
					return nil, false
				}
				n[i] = _i
			default:
				return nil, false
			}
		}
		return n, true
	}
	return nil, false
}

// Returns an slice of floats, or a nil slice
func (t Typed) Floats(key string) []float64 {
	return t.FloatsOr(key, nil)
}

// Returns an slice of floats, or the specified slice
// if the key doesn't exist or isn't a valid []float64
func (t Typed) FloatsOr(key string, d []float64) []float64 {
	n, ok := t.FloatsIf(key)
	if ok {
		return n
	}
	return d
}

// Returns a float slice + true if valid
// Returns nil + false otherwise
// (returns nil+false if one of the values is not a valid float)
func (t Typed) FloatsIf(key string) ([]float64, bool) {
	value, exists := t[key]
	if exists == false {
		return nil, false
	}
	if n, ok := value.([]float64); ok {
		return n, true
	}
	if a, ok := value.([]any); ok {
		l := len(a)
		n := make([]float64, l)
		for i := 0; i < l; i++ {
			switch t := a[i].(type) {
			case float64:
				n[i] = t
			case string:
				f, err := strconv.ParseFloat(t, 10)
				if err != nil {
					return nil, false
				}
				n[i] = f
			default:
				return nil, false
			}
		}
		return n, true
	}
	return nil, false
}

// Returns an slice of strings, or a nil slice
func (t Typed) Strings(key string) []string {
	return t.StringsOr(key, nil)
}

// Returns an slice of strings, or the specified slice
// if the key doesn't exist or isn't a valid []string
func (t Typed) StringsOr(key string, d []string) []string {
	n, ok := t.StringsIf(key)
	if ok {
		return n
	}
	return d
}

// Returns a string slice + true if valid
// Returns nil + false otherwise
// (returns nil+false if one of the values is not a valid string)
func (t Typed) StringsIf(key string) ([]string, bool) {
	value, exists := t[key]
	if exists == false {
		return nil, false
	}
	if n, ok := value.([]string); ok {
		return n, true
	}
	if a, ok := value.([]any); ok {
		l := len(a)
		n := make([]string, l)
		var ok bool
		for i := 0; i < l; i++ {
			if n[i], ok = a[i].(string); ok == false {
				return nil, false
			}
		}
		return n, true
	}
	return nil, false
}

// Returns an slice of Typed helpers, or a nil slice
func (t Typed) Objects(key string) []Typed {
	value, _ := t.ObjectsIf(key)
	return value
}

// Returns a slice of Typed helpers and true if exists, otherwise; nil and false.
func (t Typed) ObjectsIf(key string) ([]Typed, bool) {
	value, exists := t[key]
	if exists == true {
		switch t := value.(type) {
		case []any:
			l := len(t)
			n := make([]Typed, l)
			for i := 0; i < l; i++ {
				switch it := t[i].(type) {
				case map[string]any:
					n[i] = Typed(it)
				case Typed:
					n[i] = it
				}
			}
			return n, true
		case []map[string]any:
			l := len(t)
			n := make([]Typed, l)
			for i := 0; i < l; i++ {
				n[i] = Typed(t[i])
			}
			return n, true
		case []Typed:
			return t, true
		}
	}
	return nil, false
}

func (t Typed) ObjectsMust(key string) []Typed {
	value, exists := t.ObjectsIf(key)
	if exists == false {
		panic("expected objects value for " + key)
	}

	return value
}

// Returns an slice of map[string]interfaces, or a nil slice
func (t Typed) Maps(key string) []map[string]any {
	value, exists := t[key]
	if exists == true {
		if a, ok := value.([]any); ok {
			l := len(a)
			n := make([]map[string]any, l)
			for i := 0; i < l; i++ {
				n[i] = a[i].(map[string]any)
			}
			return n
		}
	}
	return nil
}

// Returns an map[string]bool
func (t Typed) StringBool(key string) map[string]bool {
	raw, ok := t.getmap(key)
	if ok == false {
		return nil
	}
	m := make(map[string]bool, len(raw))
	for k, value := range raw {
		m[k] = value.(bool)
	}
	return m
}

// Returns an map[string]int
// Some work is done to handle the fact that JSON ints
// are represented as floats.
func (t Typed) StringInt(key string) map[string]int {
	raw, ok := t.getmap(key)
	if ok == false {
		return nil
	}
	m := make(map[string]int, len(raw))
	for k, value := range raw {
		switch t := value.(type) {
		case int:
			m[k] = t
		case float64:
			m[k] = int(t)
		case string:
			i, err := strconv.Atoi(t)
			if err != nil {
				return nil
			}
			m[k] = i
		}
	}
	return m
}

// Returns an map[string]float64
func (t Typed) StringFloat(key string) map[string]float64 {
	raw, ok := t.getmap(key)
	if ok == false {
		return nil
	}
	m := make(map[string]float64, len(raw))
	for k, value := range raw {
		switch t := value.(type) {
		case float64:
			m[k] = t
		case string:
			f, err := strconv.ParseFloat(t, 10)
			if err != nil {
				return nil
			}
			m[k] = f
		default:
			return nil
		}
	}
	return m
}

// Returns an map[string]string
func (t Typed) StringString(key string) map[string]string {
	raw, ok := t.getmap(key)
	if ok == false {
		return nil
	}
	m := make(map[string]string, len(raw))
	for k, value := range raw {
		m[k] = value.(string)
	}
	return m
}

// Returns an map[string]Typed
func (t Typed) StringObject(key string) map[string]Typed {
	raw, ok := t.getmap(key)
	if ok == false {
		return nil
	}
	m := make(map[string]Typed, len(raw))
	for k, value := range raw {
		m[k] = Typed(value.(map[string]any))
	}
	return m
}

func (t Typed) Exists(key string) bool {
	_, exists := t[key]
	return exists
}

func (t Typed) Json(key string) (Typed, error) {
	value, exists := t[key]
	if !exists {
		return nil, nil
	}

	switch v := value.(type) {
	case []byte:
		return Json(v)
	case string:
		return JsonString(v)
	default:
		return nil, errors.New("expected []byte or string")
	}
}

func (t Typed) JsonMust(key string) Typed {
	t, err := t.Json(key)
	if err != nil {
		panic(err)
	}
	return t
}

func (t Typed) ToBytes(key ...string) ([]byte, error) {
	var o any
	if len(key) == 0 {
		o = t
	} else {
		exists := false
		o, exists = t[key[0]]
		if exists == false {
			return nil, KeyNotFound
		}
	}

	if o == nil {
		return nil, nil
	}
	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (t Typed) MustBytes(key string) []byte {
	data, err := t.ToBytes(key)
	if err != nil {
		panic(err)
	}
	return data
}

func (t Typed) getmap(key string) (raw map[string]any, exists bool) {
	value, exists := t[key]
	if exists == false {
		return
	}
	raw, exists = value.(map[string]any)
	return
}

func NumericToInt(value any) (int, bool) {
	switch t := value.(type) {
	case int:
		return t, true
	case int16:
		return int(t), true
	case int32:
		return int(t), true
	case int64:
		return int(t), true
	case float64:
		nt := int(t)
		if t == float64(nt) {
			return nt, true
		}
		return 0, false
	case string:
		i, err := strconv.Atoi(t)
		return i, err == nil
	}
	return 0, false
}

func Or[T any](raw any, dflt T) T {
	if raw == nil {
		return dflt
	}
	v, ok := raw.(T)
	if !ok {
		return dflt
	}
	return v
}
