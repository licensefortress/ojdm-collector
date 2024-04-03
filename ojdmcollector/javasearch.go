package ojdmcollector

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
    "strings"
)

func getJavaSharedLibFileName() string {
	switch runtime.GOOS {
	case "darwin":
		return "libjvm.dylib"
	case "linux":
		return "libjvm.so"
	case "windows":
		return "jvm.dll"
	default:
		return "libjvm.so"
	}
}

func getJavaBinaryFileName() string {
	switch runtime.GOOS {
	case "windows":
		return "java.exe"
	default:
		return "java"
	}
}

func filterPaths(searchPaths []string, defaultPaths []string) []string {
    var filteredPaths []string
    for _, path := range defaultPaths {
        exclude := false
        for _, searchPath := range searchPaths {
            if strings.HasPrefix(path, searchPath) {
                exclude = true
                break
            }
        }
        if !exclude {
            filteredPaths = append(filteredPaths, path)
        }
    }
    return filteredPaths
}

func getJavaFilePaths(searchPaths []string) []string {
    var javaSearchFilename string
    if binary {
        javaSearchFilename = getJavaBinaryFileName()
		
    } else {
        javaSearchFilename = getJavaSharedLibFileName()
    }
    defaultPaths := getSearchPaths()

    searchPaths = append(searchPaths, filterPaths(searchPaths, defaultPaths)...)

	fmt.Println("Java Search Paths: ", searchPaths)

	javaFilesMap := make(map[string]bool)
	var javaFiles []string
	for _, searchPath := range searchPaths {
		filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if os.IsPermission(err) {
					return filepath.SkipDir
				}
				return err
			}

			if !info.IsDir() {
			    if binary {
    				if info.Name() == javaSearchFilename {
	    				if _, exists := javaFilesMap[path]; !exists {
		    				fmt.Printf("Found %s in path %s\n", info.Name(), path)
			    			javaFilesMap[path] = true
				    		javaFiles = append(javaFiles, path)
					    }
				    }
                } else {                 
                    serverFolder := filepath.Base(filepath.Dir(path))
    				if info.Name() == javaSearchFilename && serverFolder == "server" {
	    				if _, exists := javaFilesMap[path]; !exists {
		    				fmt.Printf("Found %s in path %s\n", info.Name(), path)
			    			javaFilesMap[path] = true
				    		javaFiles = append(javaFiles, path)
					    }
				    }
			    }
            }

			return nil
		})
	}

	fmt.Printf("Finished gathering all java related paths!\n")
	return javaFiles
}
