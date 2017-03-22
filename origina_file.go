package configprovider

//import (
//	"fmt"
//	"io"
//	"os"
//	"path"
//
//	base_config "github.com/james-nesbitt/coach-base/operation/configuration"
//)
//
//type FileResolver interface {
//	ConfigFileName(key, scope string) (filePath string, err error)
//	ListScopes() []string
//	ListKeys(scope string) []string
//}
//
//// FileScopePathsConnector is a Connector that looks for a file in
//// various scope labeled paths, where each scope corresponds to a particular
//// path.
//type FileScopePathsConnector struct {
//	resolver FileResolver
//}
//
//// Config converts this YmlConfig to a Config interface (for clarity and validation)
//func (fspc *FileScopePathsConnector) Connector() base_config.Connector {
//	return base_config.Connector(fspc)
//}
//
//// List lists all config Keys found
//func (fspc *FileScopePathsConnector) Keys() []string {
//	keyMap := map[string]bool{}
//
//	for _, scope := range fspc.resolver.ListScopes() {
//		for _, key := range fspc.resolver.ListKeys(scope) {
//			keyMap[key] = true
//		}
//	}
//
//	keys := []string{}
//	for key, _ := range keyMap {
//		keys = append(keys, key)
//	}
//	return keys
//}
//
//// ListScopes list all scopes for a config
//func (fspc *FileScopePathsConnector) Scopes() []string {
//	return fspc.resolver.ListScopes()
//}
//
//// Get gets a configuration and apply it to a target struct
//func (fspc *FileScopePathsConnector) Get(key string, scope string) (io.ReadCloser, error) {
//	if filePath, err := fspc.resolver.configFileName(key, scope); err != nil {
//		return nil, error(ConfigNotFoundError{Key: key})
//	} else {
//		f, err := os.Open(filePath)
//		return io.ReadCloser(f), err
//	}
//}
//
//// Set sets a Config value by converting a passed struct into a configuration
//func (fspc *FileScopePathsConnector) Set(key string, scope string, source io.ReadCloser) error {
//	if filePath, err := fspc.resolver.configFileName(key, scope); err != nil {
//		return error(ConfigNotFoundError{Key: key, Scope: scope})
//	} else if file, err := os.Create(filePath); err != nil {
//		return err
//	} else {
//		defer file.Close()
//		defer source.Close()
//		_, err := io.Copy(file, source)
//		return err
//	}
//}
//
///**
// * Paths ordered set
// */
//
//// FilePathScopes An ordered list of scope keys, with corresponding file paths
//type FilePathScopes struct {
//	filename string
//	pMap     map[string]string
//	pOrder   []string
//}
//
//// Construct for FileScopePathsConnector with a filename filename that
//// starts with empty paths list
//func NewFilePathScopes(filename string) *FilePathScopes {
//	return &FilePathScopes{filename: filename}
//}
//
//// lazy initializer
//func (fps *FilePathScopes) safe() {
//	if fps.pMap == nil {
//		fps.pMap = map[string]string{}
//		fps.pOrder = []string{}
//	}
//}
//
//// List all scope keys
//func (fps *FilePathScopes) ListScopes() []string {
//	fps.safe()
//	return fps.pOrder
//}
//
//// Get a path for a scope key
//func (fps *FilePathScopes) GetFilePath(scope string) (string, error) {
//	fps.safe()
//	if p, found := fps.pMap[scope]; found {
//		return path.Join(p, fps.filename), nil
//	} else {
//		return "", error(ScopeNotFoundError{Scope: scope})
//	}
//}
//
//// Set a path for a scope key
//func (fps *FilePathScopes) SetPathScope(scope string, path string) {
//	if _, found := fps.pMap[scope]; !found {
//		fps.pOrder = append(fps.pOrder, scope)
//	}
//	fps.pMap[scope] = path
//}
//
///**
// * Errors used
// */
//
//// NoFileError indicated that no file was loaded
//type NoFileError struct {
//	Path string
//}
//
//// Error returns an error string (interface: error)
//func (nfe NoFileError) Error() string {
//	return fmt.Sprintf("No File was found at path : %s", nfe.Path)
//}
//
//// ConfigNotFoundError Config was not found Error
//type ConfigNotFoundError struct {
//	Key   string
//	Scope string
//}
//
//// Error returns an error string (interface: error)
//func (cnfe ConfigNotFoundError) Error() string {
//	return fmt.Sprintf("Config [ %s ] was not found in scope[ %s ]", cnfe.Key, cnfe.Scope)
//}
