//
// Unless otherwise noted, this project is licensed under the Creative
// Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License. Please
// see the README file.
//
// Copyright (c) 2012 The ggit Authors
//
package objects

// ================================================================= //
// TAGS
// ================================================================= //

type Tag struct {
	hdr     *ObjectHeader // the size of the tag
	oid     *ObjectId     // the oid of the tag itself
	name    string        // the tag name
	object  *ObjectId     // the object this tag is pointing at
	otype   ObjectType    // the object type
	tagger  *WhoWhen      // the tagger
	message string        // the tag message
}

func NewTag(tag, target *ObjectId, targetType ObjectType, hdr *ObjectHeader, name, msg string, tagger *WhoWhen) *Tag {
	return &Tag{
		hdr,
		tag,
		name,
		target,
		targetType,
		tagger,
		msg,
	}
}

func (t *Tag) Header() *ObjectHeader {
	return t.hdr
}

func (t *Tag) ObjectId() *ObjectId {
	return t.oid
}

func (t *Tag) Name() string {
	return t.name
}

func (t *Tag) Object() *ObjectId {
	return t.object
}

func (t *Tag) ObjectType() ObjectType {
	return t.otype
}

func (t *Tag) Tagger() *WhoWhen {
	return t.tagger
}

func (t *Tag) Message() string {
	return t.message
}
