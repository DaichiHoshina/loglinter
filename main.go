package main

import (
    "go/ast"
    "go/parser"
    "go/token"
    "log"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    // カレントディレクトリを取得
    currentDir, err := os.Getwd()
    if err != nil {
        log.Fatalf("Error getting current directory: %v\n", err)
    }

    // カレントディレクトリ以下のファイルを再帰的に検索
    err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Goファイルのみを対象にする
        if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
            checkFileForErrors(path)
        }
        return nil
    })

    if err != nil {
        log.Fatalf("Error walking the path: %v\n", err)
    }
}

// ファイル内に"ERROR"または"error"が含まれるかチェックする関数
func checkFileForErrors(filePath string) {
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
    if err != nil {
        log.Printf("Error parsing file %s: %v\n", filePath, err)
        return
    }

    // ASTをウォークしてログメッセージをチェック
    ast.Inspect(node, func(n ast.Node) bool {
        if callExpr, ok := n.(*ast.CallExpr); ok {
            for _, arg := range callExpr.Args {
                if basicLit, ok := arg.(*ast.BasicLit); ok && basicLit.Kind == token.STRING {
                    if strings.Contains(basicLit.Value, "ERROR") || strings.Contains(basicLit.Value, "error") {
                        log.Printf("Error found in file %s at position %s: %s\n",
                            filePath, fset.Position(n.Pos()), basicLit.Value)
                    }
                }
            }
        }
        return true
    })
}
