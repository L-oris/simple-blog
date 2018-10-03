package types

type t struct {
	value string
}

func (t t) String() string { return t.value }

var (
	Templates        = t{"templates"}
	FileServer       = t{"fileserver"}
	BucketRepository = t{"bucketrepository"}
	PostRepository   = t{"postrepository"}
	Router           = t{"router"}
	RootController   = t{"rootcontroller"}
	PostController   = t{"postcontroller"}
)
