package go_code_gen

import code_gen_utils "github.com/buildkite/pipeline-sdk/pkg/codegen/utils"

func newFile() *code_gen_utils.CodeGenFile {
	file := code_gen_utils.File.New(code_gen_utils.Languages.Go())
	file.AppendToHeader("package buildkite")
	return file
}