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

vim.api.nvim_create_user_command("RefreshDB", function()
	local output = {}

	vim.fn.jobstart("make testDb", {
		on_stdout = function(_, data)
			for _, line in ipairs(data) do
				if line ~= "" then
					table.insert(output, line)
				end
			end
		end,
		on_stderr = function(_, data)
			for _, line in ipairs(data) do
				if line ~= "" then
					table.insert(output, line)
				end
			end
		end,
		on_exit = function()
			vim.notify("DB Recreated:\n\t" .. table.concat(output, "\n\t"), "info")
			vim.cmd("DBCompletionClearCache")
		end,
	})
end, {})
