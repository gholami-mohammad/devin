package models

type Repository struct {
	ID                     uint64
	OwnerID                uint64
	Owner                  *User `doc:"مالک این مخزن میتواند یک فرد یا یک سازمان باشد."`
	ProjectID              uint64
	Project                *Project `doc:"هر مخزن به یک پروژه متصل است و جزء جدانشدنی آن است"`
	Name                   string   `doc:"Unique name, Unique rule is : a-z, A-Z, 0-9, dash(-), underscore(_) "`
	Title                  string   `doc:"Like name but without any unique constraint"`
	Description            string   `doc:"A breif description about the project repositpry. It is diffrent from README.md file"`
	Website                string   `doc:"Offical website of the project"`
	DefaultBranch          string
	ByteSize               float32
	WatchesCount           uint
	StarsCount             uint
	IssuesCount            uint
	ForksCount             uint
	ClosedIssuesNumber     uint
	IsBare                 bool
	IsMirror               bool
	IsPrivate              bool `doc:"درصورت خصوصی بودن تنها در دسترس افراد پروژه خواهد بود"`
	EnablePullRequest      bool
	IsForked               bool
	ForkedFromRepositoryID uint64
	ForkedFromRepository   *Repository
	Contributers           []*RepositoryContributer
}
