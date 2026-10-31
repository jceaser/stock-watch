/*******************************************************************************
Project: Several
Package: common Logic
File: common.go code sharable to multiple projects
License: Private, use by explicit agreement only
Copyright 2026 - all rights reserved by thomas.cherry@gmail.com
*******************************************************************************/

//2026-10-30

package common

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
)

/* ************************************************************************** */
// MARK: - Data Types

// loggers
type LogType struct {
	Report *log.Logger
	Error  *log.Logger
	Warn   *log.Logger
	Info   *log.Logger
	Debug  *log.Logger
	Stats  *log.Logger
}

var Log LogType

func init() {
	file := os.Stderr
	settings := log.Ldate | log.Ltime | log.Lshortfile
	Log = LogType{
		Report: log.New(file, "REPORT: ", settings),
		Error:  log.New(file, "ERROR: ", settings),
		Warn:   log.New(file, "WARNING: ", settings),
		Info:   log.New(io.Discard, "INFO: ", settings),
		Debug:  log.New(io.Discard, "DEBUG: ", settings),
		Stats:  log.New(io.Discard, "STAT: ", settings)}

	//inner app communication channels
	Talks = Talk{
		Outs: map[string]*os.File{},
		Ins:  map[string]*os.File{},
	}
}

func (self *LogType) SetLevel(value, flag int, logger *log.Logger) {
	if value&flag == flag {
		logger.SetOutput(os.Stderr)
	} else {
		logger.SetOutput(io.Discard)
	}
}

func (self *LogType) Level(value int) {
	//self.SetLevel(value, 1, Log.Report) // should always be on
	self.SetLevel(value, 2, Log.Error)
	self.SetLevel(value, 4, Log.Warn)
	self.SetLevel(value, 8, Log.Info)
	self.SetLevel(value, 16, Log.Debug)
	//self.SetLevel(value, 32, Log.Stats) // should always be controlled by file function
}

func (self *LogType) EnableInfo() {
	self.Info.SetOutput(os.Stderr)
}

func (self *LogType) EnableDebug() {
	self.Debug.SetOutput(os.Stderr)
}

func (self *LogType) EnableStats(file *os.File) {
	self.Stats.SetOutput(file)
}

// TODO: is this still used
type MachineDbReader[K any] interface {
	GetValue() K
	SetValue(K)
	GetItem(name string, dest any, def any) bool
}

type Talk struct {
	Outs map[string]*os.File
	Ins  map[string]*os.File
}

var (
	fifoFile = "/tmp/machine-"
	Talks    Talk
)

func (self *Talk) AddOut(name string) {
	path := fifoFile + name
	os.Remove(path)
	err := syscall.Mkfifo(path, 0640)
	if err != nil {
		fmt.Println("Failed to create pipe", err)
		return
	}

	pipe, err := os.OpenFile(path, os.O_RDWR, os.ModeNamedPipe)
	//O_WRONLY O_RDONLY
	if err != nil {
		fmt.Println("Couldn't open pipe with error: ", err)
		return
	}
	self.Outs[name] = pipe
}

func (self *Talk) AddIn(name string) {
	filePath := fifoFile + name
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Printf("WRITER << opened: %+v|%+v\n", f, err)
		return
	}
	self.Ins[name] = f
}

func (self *Talk) Send(name, msg string) {
	f := self.Outs[name]
	_, err := f.WriteString(fmt.Sprint(msg, "\n"))
	if err != nil {
		fmt.Printf("error sending: %v %v\n", f, err)
	}
}

func (self *Talk) Receive(name string) string {
	pipe := self.Ins[name]
	reader := bufio.NewReader(pipe)

	for {
		line, err := reader.ReadBytes('\n')
		// Close the pipe once EOF is reached
		if err != nil {
			fmt.Println("FINISHED!")
			os.Exit(1)
			//return ""
		}

		// Remove new line char
		nline := string(line)
		//nline = strings.TrimRight(nline, "\r\n")
		return nline
	}
}

func (self *Talk) Done() {
	for key, _ := range self.Outs {
		self.Outs[key].Close()
	}
}

/*
type Ifoo Foo[int]
type Ffoo Foo[float64]
*/

/* ********************************** */
// MARK: -

func Version() float32 {
	return 1.0
}

// KeysFrom extracts all keys from a given map into a slice.
//
// It is a generic function that works with any comparable key type and any value type.
// The function creates a new slice with a pre-allocated capacity equal to the map's length,
// which provides efficient memory allocation.
//
// Parameters:
//   - m: The input map from which keys will be extracted
//
// Returns:
//   - A slice containing all keys from the input map
//
// Time Complexity: O(n), where n is the number of keys in the map
// Space Complexity: O(n) to store the keys slice
//
// Example:
//
//	m := map[string]int{"a": 1, "b": 2}
//	keys := KeysFrom(m)  // keys will be ["a", "b"]
func KeysFrom[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

/*
Returns a sorted list of keys from a map.
Generic type is the map type:

	names := common.SortStringKeys[Command](Cmd.List)
*/
func SortStringKeys[V any](m map[string]V) []string {
	keys := KeysFrom[string, V](m)
	sort.Strings(keys)
	return keys
}

func IndexOf[V comparable](source []V, key V) int {
	for i, item := range source {
		if item == key {
			return i
		}
	}
	return -1
}

/**************************************/
// MARK: - Marshal functions

func StructFromJson[K any](someJson []byte) K {
	var value K
	err := json.Unmarshal(someJson, &value)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return value
}

func JsonFromStructPretty[K any](data K, pretty bool) []byte {
	result, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return result
}

func JsonFromStruct[K any](data K) []byte {
	result, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return result
}

/**************************************/
// MARK: - File Functions

func PathExpand(filePath string) string {
	if strings.HasPrefix(filePath, "~") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(homeDir, filePath[1:])
		}
	}
	return filePath
}

/*
*
Util Function Write a file
Not tested!
*/
func FileWrite(filePath string, content string) {
	d1 := []byte(content)
	filePath = PathExpand(filePath)
	err := os.WriteFile(filePath, d1, 0644)
	if err != nil {
		Log.Warn.Printf("%v", err.Error())
	}
}

/*
*
Util Function Read a file
Not tested!
@param full path to read
@return empty string on error, file contents otherwise
*/
func FileRead(filePath string) string {
	filePath = PathExpand(filePath)
	content, err := os.ReadFile(filePath)
	if err != nil {
		Log.Warn.Printf("%v", err.Error())
		return ""
	}
	return string(content)
}

/*
*
Test if a file exists
@param full path to file to be tested
@return true if no error
*/
func FileExists(filePath string) bool {
	filePath = PathExpand(filePath)
	_, err := os.Stat(filePath)
	return err == nil
}

/** convert a reader to a string */
func ReaderToString(reader io.Reader) string {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, reader)
	if err != nil {
		Log.Warn.Printf("%v", err.Error())
	}
	output := buf.String()
	return output
}

/**************************************/
// MARK: - Timer/Stat Functions

func TimeStart() int64 {
	return time.Now().UnixMicro()
}

func TimeLog(who string, start int64) {
	stop := TimeStart()
	Log.Stats.Printf("%s took %d μs", who, stop-start)
}

func GetCurrentDate(format string) string {
	if format == "" {
		format = "2006-01-02" // default format
	}
	return time.Now().Format(format)
}

/**************************************/
// MARK: - math functions

type Numbers interface {
	int | int64 | float64
}

func Min[Type Numbers](items ...Type) Type {
	_, smallest := MinMax(items...)
	return smallest
}

func Max[Type Numbers](items ...Type) Type {
	largest, _ := MinMax(items...)
	return largest
}

func MinMax[Type Numbers](items ...Type) (Type, Type) {
	var smallest, largest Type
	for i, item := range items {
		if i == 0 {
			smallest = item
			largest = item
			continue
		}
		if item < smallest {
			smallest = item
		}
		if item > largest {
			largest = item
		}
	}
	return largest, smallest
}

/**************************************/
// MARK: - other Functions

/** Find a string in a slice */
func StringInSlice(list []string, question string) bool {
	for _, item := range list {
		if strings.EqualFold(item, question) {
			return true
		}
	}
	return false
}

func PrettyMap[Type any](data map[string]Type) string {
	o := new(strings.Builder)
	keys := SortStringKeys[Type](data)
	for _, key := range keys {
		o.WriteString(fmt.Sprintf("%s: %v, ", key, data[key]))
	}
	return o.String()
}

func GenerateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

/*
func init() {
	type Test struct {
		Name string
		Value string
		Note string
	}

	data :=  map[string]float64{"energy": 12.34, "fuel": 23.3}
	fmt.Println(PrettyMap[float64](data))
	fmt.Println(PrettyMap[int](map[string]int{"energy": 12, "fuel": 23}))
	fmt.Println(PrettyMap[bool](map[string]bool{"lock": true, "run": false}))
	fmt.Println(PrettyMap[string](map[string]string{
		"loc": "home", "test": "no"}))

	fmt.Printf("unmarshal: %v\n", StructFromJson[Test]([]byte(
		"{\"Name\": \"name\", \"Value\": \"some value\"}")))
}
*/
