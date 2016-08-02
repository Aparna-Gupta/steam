package compiler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/h2oai/steamY/lib/fs"
)

const (
	ArtifactWar       = "war"
	ArtifactJar       = "jar"
	ArtifactPythonWar = "python_war"
)

const (
	fileTypeJava        = "pojo"
	fileTypeJavaDep     = "jar"
	fileTypePythonMain  = "python"
	fileTypePythonOther = "pythonextra"
)

type Service struct {
	Address string
}

func NewService(address string) *Service {
	return &Service{
		address,
	}
}

func (s *Service) urlFor(slug string) string {
	return (&url.URL{Scheme: "http", Host: s.Address, Path: slug}).String()
}

func (s *Service) Ping() error {
	if _, err := http.Get(s.urlFor("ping")); err != nil {
		return fmt.Errorf("Failed connecting to scoring service builder: %s", err)
	}
	return nil
}

func attachFile(w *multipart.Writer, filePath, fileType string) error {
	dst, err := w.CreateFormFile(fileType, path.Base(filePath))
	if err != nil {
		return fmt.Errorf("Failed writing to buffer: %v", err)
	}
	src, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Failed opening file: %v", err)
	}
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("Failed copying file: %v", err)
	}

	return nil
}

func compile(url, javaFilePath, javaDepPath, pythonMainFilePath string, pythonOtherFilePaths []string) (*http.Response, error) {
	javaFilePath, err := fs.ResolvePath(javaFilePath)
	if err != nil {
		return nil, err
	}
	javaDepPath, err = fs.ResolvePath(javaDepPath)
	if err != nil {
		return nil, err
	}

	b := &bytes.Buffer{}
	writer := multipart.NewWriter(b)

	if err := attachFile(writer, javaFilePath, fileTypeJava); err != nil {
		return nil, fmt.Errorf("Failed attaching Java file to compilation request: %s", err)
	}

	if err := attachFile(writer, javaDepPath, fileTypeJavaDep); err != nil {
		return nil, fmt.Errorf("Failed attaching Java dependency to compilation request: %s", err)
	}

	if len(pythonMainFilePath) > 0 {
		if err := attachFile(writer, pythonMainFilePath, fileTypePythonMain); err != nil {
			return nil, fmt.Errorf("Failed attaching Python main file to compilation request: %s", err)
		}
		if len(pythonOtherFilePaths) > 0 {
			for _, p := range pythonOtherFilePaths {
				if err := attachFile(writer, p, fileTypePythonOther); err != nil {
					return nil, fmt.Errorf("Failed attaching Python file to compilation request: %s", err)
				}
			}
		}
	}

	ct := writer.FormDataContentType()
	writer.Close()

	res, err := http.Post(url, ct, b)
	if err != nil {
		return nil, fmt.Errorf("Failed making compilation request: %v", err)
	}

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed reading compilation response: %v", err)
		}
		return nil, fmt.Errorf("Failed compiling scoring service: %s / %s", res.Status, string(body))
	}

	return res, nil
}

func (s *Service) GetPythonFilePaths(wd string, projectId int64, packageName string) (string, []string, error) {
	var pythonMainFilePath string
	var pythonOtherFilePaths []string

	packagePath := fs.GetPackagePath(wd, projectId, packageName)

	if !fs.DirExists(packagePath) {
		return "", nil, fmt.Errorf("Package %s does not exist")
	}

	packageAttrsBytes, err := fs.GetPackageAttributes(wd, projectId, packageName)
	if err != nil {
		return "", nil, fmt.Errorf("Failed reading package attributes: %s", err)
	}

	packageAttrs, err := fs.JsonToMap(packageAttrsBytes)
	if err != nil {
		return "", nil, fmt.Errorf("Failed parsing package attributes: %s", err)
	}

	pythonMain, ok := packageAttrs["main"]
	if !ok {
		return "", nil, fmt.Errorf("Failed determining Python main file from package attributes")
	}

	packageFileList, err := fs.ListFiles(packagePath)
	if err != nil {
		return "", nil, fmt.Errorf("Failed reading package file list: %s", err)
	}

	// Filter .py files; separate ancillary files from the main one.
	pythonOtherFilePaths = make([]string, 0)
	for _, f := range packageFileList {
		if strings.ToLower(path.Ext(f)) == ".py" {
			p := path.Join(packagePath, f)
			if f == pythonMain {
				pythonMainFilePath = p
			} else {
				pythonOtherFilePaths = append(pythonOtherFilePaths, p)
			}
		}
	}

	if len(pythonMainFilePath) == 0 {
		return "", nil, fmt.Errorf("Failed locating Python main file in package file listing")
	}

	return pythonMainFilePath, pythonOtherFilePaths, nil
}

func (s *Service) CompileModel(wd string, modelId int64, modelLogicalName, artifact string) (string, error) {

	genModelPath := fs.GetGenModelPath(wd, modelId)
	javaModelPath := fs.GetJavaModelPath(wd, modelId, modelLogicalName)

	var targetFile, slug string

	switch artifact {
	case ArtifactWar:
		targetFile = fs.GetWarFilePath(wd, modelId, modelLogicalName)
		slug = "makewar"
	case ArtifactPythonWar:
		targetFile = fs.GetPythonWarFilePath(wd, modelId, modelLogicalName)
		slug = "makepythonwar"
	case ArtifactJar:
		targetFile = fs.GetModelJarFilePath(wd, modelId, modelLogicalName)
		slug = "compile"
	}

	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
	} else {
		return targetFile, nil
	}

	if err := s.Ping(); err != nil {
		return "", err
	}

	// XXX
	res, err := compile(s.urlFor(slug), javaModelPath, genModelPath, "", nil)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	dst, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", fmt.Errorf("Download file createion failed: %s: %v", targetFile, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, res.Body); err != nil {
		return "", fmt.Errorf("Download file copy failed: Service to %s: %v", targetFile, err)
	}
	return targetFile, nil
}
