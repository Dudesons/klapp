package flip

var Flip_fixture = `
[
  {
    "activated": true,
    "name": "boolon",
    "description": ""
  },
  {
    "activated": false,
    "name": "booloff",
    "description": ""
  }
]
`
var Flip_string_fixture = `
[
  {
    "activated": true,
    "name": "string",
    "description": "",
    "content": "mystring"
  },
  {
    "activated": true,
    "name": "stringslice",
    "description": "",
    "content": ["q", "w", "e", "r", "t", "y"]
  }
]
`
var Flip_int_fixture = `
[
  {
    "activated": true,
    "name": "int",
    "description": "",
    "content": 42
  },
  {
    "activated": true,
    "name": "intslice",
    "description": "",
    "content": [0, 1, 2, 3, 4, 5]
  }
]
`

type Fixture struct {
	Success bool
	Flip string
	Target interface{}
	FlipType int // 0=bool 1=string 2=stringslice 3=int 4=intslice
}

var FlipNamesForBool = []Fixture{
	{Flip: "boolon", Success: true},
	{Flip: "booloff", Success: false},
	{Flip: "string", Target: "mystring", Success: false},
	{Flip: "string", Target: "plop", Success: false},
	{Flip: "int", Target: "42", Success: false},
	{Flip: "int", Target: "44", Success: false},
	{Flip: "stringslice", Target: "q", Success: false},
	{Flip: "stringslice", Target: "a", Success: false},
	{Flip: "intslice", Target: "0", Success: false},
	{Flip: "intslice", Target: "-1", Success: false},
	// TODO(Integrate this error case)
	//{Flip: "dontexist", Success: false},
}

var FlipsFixture = []Fixture{
	{Flip: "boolon", Success: true, FlipType: 0},
	{Flip: "booloff", Success: false},
	{Flip: "string", Target: "mystring", Success: true, FlipType: 1},
	{Flip: "string", Target: "plop", Success: false},
	{Flip: "int", Target: 42, Success: true, FlipType: 2},
	{Flip: "int", Target: 44, Success: false},
	{Flip: "stringslice", Target: "q", Success: true, FlipType: 3},
	{Flip: "stringslice", Target: "a", Success: false},
	{Flip: "intslice", Target: 0, Success: true, FlipType: 4},
	{Flip: "intslice", Target: -1, Success: false},
	// TODO(Integrate this error case)
	//{Flip: "dontexist", Success: false},
}
