#!/bin/bash
base_dir=$(dirname "$0")
base_dir_full=$(realpath $base_dir)
echo $base_dir_full
project_root=$(realpath $base_dir_full/..)
application_exe=${project_root}/tmp/main
rm -rf /tmp/ramdisk/test_project_*
test_project_dir="/tmp/ramdisk/test_project_"$(date +%Y%m%d_%H%M)
rm -rf $test_project_dir
mkdir -p $test_project_dir
cat << EOF > ${test_project_dir}/.nvimrc.lua
vim.filetype.add({
  extension = {
    tsf = "timesheet",
  },
})

vim.lsp.config("timesheet_lsp", {
  cmd = function(dispatchers, config)
    local root = config.root_dir or vim.fn.getcwd()
    return vim.lsp.rpc.start({
      "$application_exe",
      "-c", root .. "/config.toml",
      "--project-root", root,
    }, dispatchers, {
      cwd = root,
    })
  end,
  filetypes = { "timesheet" },
  root_markers = { ".git", "config.toml" },
})

vim.lsp.enable("timesheet_lsp")
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
