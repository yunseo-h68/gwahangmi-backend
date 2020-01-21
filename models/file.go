package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ImageFile 은 이미지 파일에 대한 정보를 담고 있는 구조체입니다.
type ImageFile struct {
	ID       primitive.ObjectID `bson:"_id"`
	Filename string             `bson:"filename"`
	MetaData FileMeta           `bson:"metadata"`
}

// PostContent 는 글 본문에 대한 정보를 담고 있는 구조체입니다.
type PostContent struct {
	ID       primitive.ObjectID `bson:"_id"`
	Filename string             `bson:"filename"`
	MetaData FileMeta           `bson:"metadata"`
}

// FileMeta 는 Upload할 파일의 메타정보를 담는 구조체입니다.
type FileMeta struct {
	Inode int
	UID   string `bson:"uid" json:"uid"`
	Ext   string `bson:"ext" json:"ext"`
}
