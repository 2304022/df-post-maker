package dto

// Visibility enum
type Visibility string

const (
	VisibilityPublic       Visibility = "PUBLIC"
	VisibilityUser         Visibility = "USER"
	VisibilityOrganization Visibility = "ORGANIZATION"
)

// FileLinkDto — прикрепляемый файл
type FileLinkDto struct {
	FileId  int    `json:"fileId" binding:"required"`
	Caption string `json:"caption"`
}

// PostTopicViewDto — view внутри jsonView
type PostTopicViewDto struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// PostTopicJsonViewDto — jsonView топика
type PostTopicJsonViewDto struct {
	TypeTopic string           `json:"type_topic"`
	Name      string           `json:"name"`
	Version   string           `json:"version"`
	View      PostTopicViewDto `json:"view"`
}

// PostTopicDto — топик поста (обязательное поле)
type PostTopicDto struct {
	Name            string               `json:"name"`
	VersionTemplate string               `json:"versionTemplate"`
	JsonView        PostTopicJsonViewDto `json:"jsonView"`
	View            PostTopicViewDto     `json:"view"`
}

// CreatePostRequestDto — тело запроса создания поста
type CreatePostRequestDto struct {
	Text             string        `json:"text" binding:"required"`
	ShortText        string        `json:"shortText"`
	Visibility       Visibility    `json:"visibility"`
	Files            []FileLinkDto `json:"files"`
	Topic            PostTopicDto  `json:"topic" binding:"required"`
	Tags             []string      `json:"tags"`
	PublishInFeed    bool          `json:"publishInFeed"`
	OrganisationPost *bool         `json:"organisationPost,omitempty"`
	Nsfw             *bool         `json:"nsfw,omitempty"`
	DisableComments  bool          `json:"disableComments"`
}

// SetDefaults устанавливает значения по умолчанию
func (r *CreatePostRequestDto) SetDefaults() {
	if r.Visibility == "" {
		r.Visibility = VisibilityUser
	}
	if r.Topic.Name == "" {
		r.Topic.Name = "post"
	}
	if r.Topic.VersionTemplate == "" {
		r.Topic.VersionTemplate = "1"
	}
	if r.Topic.JsonView.Name == "" {
		r.Topic.JsonView.Name = "post"
	}
	if r.Topic.JsonView.TypeTopic == "" {
		r.Topic.JsonView.TypeTopic = "article"
	}
	if r.Topic.JsonView.Version == "" {
		r.Topic.JsonView.Version = "1"
	}
}

// UploadFileInfo — информация об одном загруженном файле
type UploadFileInfo struct {
	Id                int                    `json:"id"`
	Size              int                    `json:"size"`
	Number            int                    `json:"number"`
	ContentType       string                 `json:"contentType"`
	HttpPath          string                 `json:"httpPath"`
	FileName          string                 `json:"fileName"`
	Extension         string                 `json:"extension"`
	OriginalExtension string                 `json:"originalExtension"`
	Info              map[string]interface{} `json:"info"`
	Width             int                    `json:"width"`
	Height            int                    `json:"height"`
	PreviewWidth      int                    `json:"previewWidth"`
	PreviewHeight     int                    `json:"previewHeight"`
	ValidImage        bool                   `json:"validImage"`
	ClientTitle       string                 `json:"clientTitle"`
}

// UploadResponseDto — ответ от direct.farm после загрузки файла
type UploadResponseDto struct {
	Files     []UploadFileInfo `json:"files"`
	Errors    []interface{}    `json:"errors"`
	Exception string           `json:"exception"`
	Error     interface{}      `json:"error"`
}

// CreatePostResponseDto — ответ после создания поста
type CreatePostResponseDto struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Visibility  Visibility `json:"visibility"`
	Tags        []string   `json:"tags"`
	Uuid        string     `json:"uuid"`
	TransTitle  string     `json:"transTitle"`
	NewPost     bool       `json:"newPost"`
	Nsfw        bool       `json:"nsfw"`
}
