#!/bin/bash
base_dir=$(dirname "$0")
base_dir_full=$(realpath $base_dir)
echo $base_dir_full
project_root=$(realpath $base_dir_full/..)
application_exe=${project_root}/tmp/main
test_project_dir="/mnt/ramdisk/test_project"
rm -rf $test_project_dir
mkdir -p $test_project_dir
cat << EOF > ${test_project_dir}/.nvimrc.lua
vim.filetype.add({
  extension = {
    tsf = "timesheet"
  }
})
local lspconfig = require("lspconfig")

lspconfig.timesheet_lsp.setup({
  on_attach = function(client, bufnr)
  end,
  root_dir = function(fname)
    return lspconfig.util.root_pattern(".git")(fname) or vim.fn.getcwd()
  end,
})

vim.api.nvim_create_autocmd("FileType", {
  pattern = "timesheet",
  callback = function()
    local root_dir = lspconfig.util.root_pattern(".git")(vim.fn.expand("%:p")) or vim.fn.getcwd()
    vim.lsp.start({
      name = "timesheet_lsp",
      cmd = {"$application_exe", "-c", "$test_project_dir/config.toml", "--project-root", root_dir},
      root_dir = root_dir,
    })
  end,
})
EOF
cat << EOF > $test_project_dir/config.toml
[categories]
regular=["categoryA", "categoryB"]
overtime=["overtimeA"]
[holidays]
repeatable=["11-11","05-01","01-01", "12-25"]
addhoc=["2021-02-01"]
[tasks]
prefix="task-"
onlyNumbers=true
EOF
mkdir -p $test_project_dir/2025/02
touch $test_project_dir/2025/02/01.tsf
touch $test_project_dir/2025/02/02.tsf
pushd $test_project_dir || exit
git init .
nvim . 
popd || exit
