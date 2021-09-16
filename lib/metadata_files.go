package lib

// func (z *Zettel) RemoveFile(fd metadata.File) (err error) {
// 	found := -1
// 	tags := z.Metadata.Tags
// 	tag := fd.Tag()

// 	for i, t := range tags {
// 		if t == tag {
// 			found = i
// 			break
// 		}
// 	}

// 	if found == -1 {
// 		err = xerrors.Errorf("tag not found: %s", tag)
// 		return
// 	}

// 	tags[found] = tags[len(tags)-1]
// 	z.Metadata.Tags = tags[:len(tags)-1]

// 	return
// }
