package api

var flip_fixture = `
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
  },
  {
    "activated": true,
    "name": "string",
    "description": "",
    "content": "mystring"
  },
  {
    "activated": true,
    "name": "int",
    "description": "",
    "content": 42
  },
  {
    "activated": true,
    "name": "stringslice",
    "description": "",
    "content": ["q", "w", "e", "r", "t", "y"],
    "type": 0
  },
  {
    "activated": true,
    "name": "intslice",
    "description": "",
    "content": [0, 1, 2, 3, 4, 5],
    "type": 1
  }
]
`

type Fixture struct {
	Success bool
	Flip string
	Target string
}

var flipNames = []Fixture{
	{Flip: "boolon", Success: true},
	{Flip: "booloff", Success: false},
	{Flip: "string", Target: "mystring", Success: true},
	{Flip: "string", Target: "plop", Success: false},
	{Flip: "int", Target: "42", Success: true},
	{Flip: "int", Target: "44", Success: false},
	{Flip: "stringslice", Target: "q", Success: true},
	{Flip: "stringslice", Target: "a", Success: false},
	{Flip: "intslice", Target: "0", Success: true},
	{Flip: "intslice", Target: "-1", Success: false},
	// TODO(Integrate this error case)
	//{Flip: "dontexist", Success: false},
}
