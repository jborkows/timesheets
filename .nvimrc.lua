vim.g.dbs = {
	{
		name = "myTestDB",
		url = "sqlite:temp/timesheets.db",
	},
}
vim.g.db_ui_use_nerd_fonts = 1
vim.g.db_ui_execute_on_save = 0

vim.api.nvim_create_autocmd("FileType", {
	pattern = "sql",
	callback = function()
		vim.bo.omnifunc = "vim_dadbod_completion#omni"
	end,
})
vim.g.db = vim.g.dbs[1]

vim.api.nvim_create_autocmd("FileType", {
	pattern = { "dbui", "sql" },
	callback = function()
		vim.keymap.set("n", "<F8>", "<Plug>(DBUI_ExecuteQuery)", { silent = true, buffer = true })
	end,
})
