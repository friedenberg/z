package lib

////TODO-P2 swithch to p rintable description
//type CleanActionCheck func(z *Zettel) bool
//type CleanActionPerform func(z *Zettel) (bool, error)

//type CleanAction struct {
//	Check   CleanActionCheck
//	Perform CleanActionPerform
//}

//func GetCleanActions() map[string]CleanAction {
//	return map[string]CleanAction{
//		"delete_if_missing_file": CleanAction{shouldDeleteIfMissingFile, deleteIfMissingFile},
//		"normalize_file": CleanAction{
//			func(z *Zettel) bool {
//				fd, ok := z.LocalFile()

//				if !ok {
//					return false
//				}

//				normalizedFile := path.Base(fd.FileName())

//				return normalizedFile != z.Metadata.File
//			},
//			func(z *Zettel) (shouldWrite bool, err error) {
//				z.Metadata.File = path.Base(z.Metadata.File)
//				shouldWrite = true
//				return
//			},
//		},
//		"rewrite_metadata": CleanAction{
//			func(z *Zettel) bool {
//				oldYaml := z.Data.MetadataYaml
//				//TODO-P3 handle err
//				z.generateMetadataYaml()

//				return oldYaml != z.Data.MetadataYaml
//			},
//			func(z *Zettel) (shouldWrite bool, err error) {
//				shouldWrite = true
//				return
//			},
//		},
//		//index
//		// "remove_from_index":             CleanAction{shouldRemoveFromIndex, removeFromIndex},
//		// "add_to_index":             CleanAction{shouldRemoveFromIndex, removeFromIndex},
//		// "update_in_index":             CleanAction{shouldUpdateInIndex, updateInIndex},
//		//TODO-P3 file attachment
//		//TODO-P3 change file permissions
//		//TODO-P3 reformat yaml
//	}
//}

//func shouldDeleteIfMissingFile(z *Zettel) bool {
//	if !z.HasFile() {
//		return false
//	}

//	return !util.FileExists(z.FilePath())
//}

//func deleteIfMissingFile(z *Zettel) (shouldWrite bool, err error) {
//	err = os.Remove(z.Path)
//	return
//}
