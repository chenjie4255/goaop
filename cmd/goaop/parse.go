package main

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var fileContent []byte

func init() {
}

func parseFile(path string) (renderData, error) {
	data := renderData{}
	fc, err := os.Open(path)
	if err != nil {
		return data, err
	}
	fileContent, _ = ioutil.ReadAll(fc)
	fc.Close()

	fset := token.NewFileSet()

	fmt.Println("parser.ParseFile", path)
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return data, err
	}

	packageName, err := parsePackageName(f)
	if err != nil {
		return data, err
	}
	data.Package = packageName

	parseInterfaces(f, fset, &data)

	return data, nil
}

func parsePackageName(f *ast.File) (string, error) {
	if f.Name == nil {
		return "", errors.New("no package name found")
	}
	return f.Name.Name, nil
}

func parseInterfaces(f *ast.File, fset *token.FileSet, data *renderData) error {
	importMap := map[string]string{}
	fmt.Println(f.Scope.String())
	for _, imp := range f.Imports {
		if imp.Name != nil {
			importMap[imp.Name.Name] = imp.Path.Value
		} else {
			dirset := token.FileSet{}
			goPath := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))
			goRoot := build.Default.GOROOT
			if goRoot != "" {
				goPath = append(goPath, goRoot)
			}

			if len(goPath) == 0 {
				goPath = append(goPath, build.Default.GOPATH)
			} else if goPath[0] == "" {
				goPath[0] = build.Default.GOPATH
			}

			build.Default.SrcDirs()

			modFilePath := ""
			if os.Getenv("GO111MODULE") == "on" {
				fmt.Println("GO111MODULE on")
				// find go.mod file
				//info, err := os.Stat(inputFile)
				//if err != nil {
				//	panic(err)
				//}
				wd, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				fPath := filepath.Join(wd, inputFile)
				goModDir := filepath.Dir(fPath)
				for {
					info, err := os.Stat(filepath.Join(goModDir, "go.mod"))
					if err == nil && !info.IsDir() {
						break
					}
					d := filepath.Dir(goModDir)
					if len(d) >= len(goModDir) {
						goModDir = ""
						break
					}
					goModDir = d
				}

				if goModDir != "" {
					fmt.Println("find mod.file file's dir", goModDir)
					// read go.mod file
					modF, err := os.Open(filepath.Join(goModDir, "go.mod"))
					if err != nil {
						panic(err)
					}
					firstLine, err := bufio.NewReader(modF).ReadString('\n')
					if err != nil {
						panic(err)
					}
					modF.Close()
					results := strings.Split(firstLine, " ")
					moduleRoot := strings.TrimSpace(results[1])
					importPath := strings.Trim(imp.Path.Value, `"`)
					if strings.HasPrefix(importPath, moduleRoot) {
						modFilePath = filepath.Join(goModDir, strings.TrimPrefix(importPath, moduleRoot))
					}
					fmt.Println("go module path", modFilePath)
				}
			}

			for _, gp := range goPath {
				dirPath := gp+"/src/"+strings.Trim(imp.Path.Value, `"`)
				pkgs, err := parser.ParseDir(&dirset, dirPath, nil, parser.PackageClauseOnly)
				if err != nil {
					fmt.Println(dirPath, err)
					continue
				}
				for k := range pkgs {
					importMap[k] = imp.Path.Value
				}
				fmt.Println("pkgs", pkgs)
			}
			if modFilePath != ""{
				pkgs, err := parser.ParseDir(&dirset, modFilePath, nil, parser.PackageClauseOnly)
				if err != nil {
					fmt.Println(imp.Path.Value, err)
					continue
				}
				for k := range pkgs {
					importMap[k] = imp.Path.Value
				}
				fmt.Println("pkgs", pkgs)
			}
		}
	}

	ctx := parseContext{}
	ctx.importMap = importMap
	ctx.fset = fset

	for _, decl := range f.Decls {
		//
		switch t := decl.(type) {
		case *ast.FuncDecl:
			fmt.Println("FuncDecl", t.Name)
		case *ast.GenDecl:
			fmt.Println("Gen Decl", t.Tok)
		}

		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		// need to be 'type xxx interface'
		if genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			ifType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			if genDecl.Doc == nil {
				continue
			}

			if hasSignComment := checkComment(genDecl.Doc.Text()); !hasSignComment {
				fmt.Println("no comment", genDecl.Doc.Text())
				continue
			}

			if data.InterfaceName == "" {
				fmt.Printf("interface: %s\n", typeSpec.Name.Name)
				data.InterfaceName = typeSpec.Name.Name
			} else {
				break
			}

			for _, method := range ifType.Methods.List {
				if len(method.Names) == 0 {
					continue
				}

				funcType, ok := method.Type.(*ast.FuncType)
				if !ok {
					continue
				}

				ps := fset.Position(method.Pos())
				pe := fset.Position(method.End())

				renderM := renderMethod{}
				renderM.Name = method.Names[0].Name
				renderM.Raw = string(fileContent[ps.Offset:pe.Offset])

				for _, param := range funcType.Params.List {
					// check if xxx.bbbb, we need to get the selector xxx
					impObj := getImportObj(&ctx, param.Pos(), param.End())
					if impObj != nil {
						data.addImportPath(impObj.Name, impObj.Path)
					}

					hasVarargs := containsVarargs(&ctx, param.Pos(), param.End())

					for _, paramName := range param.Names { //processing case with func(v1, v2 int)
						renderP := renderParam{}
						renderP.Name = paramName.Name
						renderP.Varargs = hasVarargs

						renderM.Params = append(renderM.Params, renderP)

						// paramStr := paramName.Name

						// switch t := param.Type.(type) {
						// case *ast.Ident:
						// 	fmt.Println("ident", t.Name)
						// 	paramStr += " " + t.Name
						// case *ast.SelectorExpr:
						// 	fmt.Println("SelectorExpr", t.Sel.Name, t.Sel.IsExported())
						// 	if t.Sel.IsExported() {
						// 		fmt.Println("X", t.X)
						// 		paramStr += " " + t.X.(*ast.Ident).String() + "." + t.Sel.Name
						// 	} else {
						// 		paramStr += " " + t.Sel.Name
						// 	}
						// case *ast.StarExpr:
						// 	fmt.Println("StarExpr.x", t.X)
						// 	switch txx := t.X.(type) {
						// 	case *ast.SelectorExpr:
						// 		if txx.Sel.IsExported() {
						// 			paramStr += " *" + txx.X.(*ast.Ident).String() + "." + txx.Sel.Name
						// 		} else {
						// 			paramStr += " *" + txx.Sel.Name
						// 		}
						// 	case *ast.Ident:
						// 		paramStr += " *" + txx.Name
						// 	}
						// case *ast.UnaryExpr:
						// 	fmt.Println("UnaryExpr", t.Op.String())
						// case *ast.SliceExpr:
						// 	fmt.Println("SliceExpr", t.X)
						// case *ast.IndexExpr:
						// 	fmt.Println("IndexExpr", t.X)
						// case *ast.ArrayType:
						// 	fmt.Println("ArrayType", t)
						// default:
						// 	fmt.Println(t.(*ast.IndexExpr).Pos())
						// 	fmt.Println("?????", t)
						// }

						// fmt.Println("param:", paramStr, "\n")
					}
				}

				if funcType.Results != nil {
					renderM.ResultCount = len(funcType.Results.List)
					renderM.ResultErrorIndex = -1
					for i, result := range funcType.Results.List {
						impObj := getImportObj(&ctx, result.Pos(), result.End())
						if impObj != nil {
							data.addImportPath(impObj.Name, impObj.Path)
						}
						switch t := result.Type.(type) {
						case *ast.Ident:
							if t.String() == "error" {
								renderM.ResultErrorIndex = i
							}
						}
					}
				}

				data.Methods = append(data.Methods, renderM)
			}
			break
		}
	}

	return nil
}

type importObj struct {
	Name string // may be alise
	Path string
}

type parseContext struct {
	fset      *token.FileSet
	importMap map[string]string
}

func getImportObj(ctx *parseContext, pos token.Pos, end token.Pos) *importObj {
	startOffset := ctx.fset.Position(pos)
	endOffset := ctx.fset.Position(end)

	content := string(fileContent[startOffset.Offset:endOffset.Offset])

	reg := regexp.MustCompile(selectorReg)
	result := reg.FindStringSubmatch(content)
	if len(result) != 2 {
		return nil
	}

	importName := result[1]
	val, found := ctx.importMap[importName]
	if !found {
		panic("invalid selector: " + importName)
	}

	return &importObj{importName, val}
}

func containsVarargs(ctx *parseContext, pos token.Pos, end token.Pos) bool {
	startOffset := ctx.fset.Position(pos)
	endOffset := ctx.fset.Position(end)

	content := string(fileContent[startOffset.Offset:endOffset.Offset])

	reg := regexp.MustCompile(varargs)
	result := reg.FindStringSubmatch(content)
	return len(result) == 1
}

const (
	parseReg    = `@ifmeasure`
	selectorReg = `(\w+)\.`
	varargs     = `\.\.\.`
)

func checkComment(text string) bool {
	reg := regexp.MustCompile(parseReg)
	result := reg.FindStringSubmatch(text)
	return len(result) == 1
}
