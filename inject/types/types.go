package types

type t struct {
	value string
}

func (t t) String() string { return t.value }

var (
	BucketRepository = t{"bucketrepository"}
	FileServer       = t{"fileserver"}
	DB               = t{"db"}
	PostRepository   = t{"postrepository"}
	RootController   = t{"rootcontroller"}
	PostController   = t{"postcontroller"}
	PostService      = t{"postservice"}
	Router           = t{"router"}
	Template         = t{"template"}
)
