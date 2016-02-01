package main

type User struct {
	Id						int    			`json:"id"`
	FirstName				string    		`json:"first_name"`
	LastName				string    		`json:"last_name"`
	UserName				string   		`json:"username"`
}

type Chat struct {
	Id						int    			`json:"id"`
	ChatType				string 			`json:"type"`
	Title					string    		`json:"title"`
	UserName				string    		`json:"username"`
	FirstName				string    		`json:"first_name"`
	LastName				string			`json:"last_name"`
}

type PhotoSize struct {
	FileId					string          `json:"file_id"`
	Width					int             `json:"width"`
	Height					int 			`json:"height"`
	FileSize				int        		`json:"file_size"`
}

type Audio struct {
	FileId					string    		`json:"file_id"`
	Duration				int    			`json:"duration"`
	Performer				string    		`json:"performer"`
	Title					string    		`json:"title"`
	MimeType				string    		`json:"mime_type"`
	FileSize				int    			`json:"file_size"`
}

type Document struct {
	File_id					string          `json:"file_id"`
	Thumb					PhotoSize       `json:"thumb"`
	FileName				string          `json:"file_name"`
	MimeType				string          `json:"mime_type"`
	FileSize				int             `json:"file_size"`
}

type Sticker struct {
	FileId					string          `json:"file_id"`
	Width					int            	`json:"width"`
	Height					int            	`json:"height"`
	Thumb					PhotoSize       `json:"thumb"`
	FileSize				int            	`json:"file_size"`
}

type Video struct {
	FileId					string       	`json:"file_id"`
	Width					int        		`json:"width"`
	Height					int        		`json:"height"`
	Duration				int    			`json:"duration"`
	Thumb					PhotoSize    	`json:"thumb"`
	MimeType				string    		`json:"mime_type"`
	FileSize				int    			`json:"file_size"`
}

type Voice struct {
	FileId					string          `json:"file_id"`
	Duration				int             `json:"duration"`
	MimeType				string          `json:"mime_type"`
	FileSize				int             `json:"file_size"`
}

type Contact struct {
	PhoneNumber				string   		`json:"phone_number"`
	FirstName				string    		`json:"first_name"`
	LastName				string    		`json:"last_name"`
	UserId					int        		`json:"user_id"`
}

type Location struct {
	Longitude 				float64    		`json:"longitude"`
	Latitude    			float64    		`json:"latitude"`
}

type UserProfilePhotos struct {
	TotalCount 				int 			`json:"total_count"`
	Photos    				[][]PhotoSize   `json:"photos"`
}

type MessageWithoutReply struct {
	MessageId 				int 			`json:"message_id"`
	From					User    		`json:"from"`
	Date					int 			`json:"date"`
	Chat					Chat    		`json:"chat"`
	ForwardFrom				User    		`json:"forward_from"`
	ForwardDate				int				`json:"forward_date"`
	Text					string			`json:"text"`
	Audio					Audio    		`json:"audio"`
	Document				Document    	`json:"document"`
	Photo					[]PhotoSize 	`json:"photo"`
	Sticker					Sticker    		`json:"sticker"`
	Video					Video    		`json:"video"`
	Voice					Voice   		`json:"voice"`
	Caption					string    		`json:"caption"`
	Contact					Contact    		`json:"contact"`
	Location				Location    	`json:"location"`
	NewChatParticipant		User    		`json:"new_chat_participant"`
	LeftChatParticipant		User    		`json:"left_chat_participant"`
	NewChatTitle			string    		`json:"new_chat_title"`
	NewChatPhoto			[]PhotoSize 	`json:"new_chat_photo"`
	DeleteChatPhoto			bool    		`json:"delete_chat_photo"`
	GroupChatCreated		bool    		`json:"group_chat_created"`
	SupergroupChatCreated	bool    		`json:"supergroup_chat_created"`
	ChannelChatCreated		bool    		`json:"channel_chat_created"`
	MigrateToChatId			int    			`json:"migrate_to_chat_id"`
	MigrateFromChatId		int    			`json:"migrate_from_chat_id"`
}

type MessageWithReply struct {
	MessageId 				int 			`json:"message_id"`
	From					User    		`json:"from"`
	Date					int 			`json:"date"`
	Chat					Chat    		`json:"chat"`
	ForwardFrom				User    		`json:"forward_from"`
	ForwardDate				int				`json:"forward_date"`
	Text					string			`json:"text"`
	Audio					Audio    		`json:"audio"`
	Document				Document    	`json:"document"`
	Photo					[]PhotoSize 	`json:"photo"`
	Sticker					Sticker    		`json:"sticker"`
	Video					Video    		`json:"video"`
	Voice					Voice   		`json:"voice"`
	Caption					string    		`json:"caption"`
	Contact					Contact    		`json:"contact"`
	Location				Location    	`json:"location"`
	NewChatParticipant		User    		`json:"new_chat_participant"`
	LeftChatParticipant		User    		`json:"left_chat_participant"`
	NewChatTitle			string    		`json:"new_chat_title"`
	NewChatPhoto			[]PhotoSize 	`json:"new_chat_photo"`
	DeleteChatPhoto			bool    		`json:"delete_chat_photo"`
	GroupChatCreated		bool    		`json:"group_chat_created"`
	SupergroupChatCreated	bool    		`json:"supergroup_chat_created"`
	ChannelChatCreated		bool    		`json:"channel_chat_created"`
	MigrateToChatId			int    			`json:"migrate_to_chat_id"`
	MigrateFromChatId		int    			`json:"migrate_from_chat_id"`

	replyToMessage 			MessageWithoutReply    		`json:"reply_to_message"`
}

type InlineQuery struct  {
	Id						string 			`json:"id"`
	From					User            `json:"from"`
	Query					string          `json:"query"`
	Offset					string          `json:"offset"`
}

type ChosenInlineResult struct {
	ResultId				string    		`json:"result_id"`
	From					User    		`json:"from"`
	Query					string    		`json:"query"`
}

type Update struct {
	UpdateId				uint64				`json:"update_id"`
	Message					MessageWithReply    `json:"message"`
	InlineQuery				InlineQuery         `json:"inline_query"`
	ChosenInlineResult		ChosenInlineResult	`json:"chosen_inline_result"`
}