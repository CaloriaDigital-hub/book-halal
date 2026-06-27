package reader

import (
	readerCommands "book_halal/internal/application/reader/commands"
	readerQueries "book_halal/internal/application/reader/queries"
)

type ReaderController struct {
	saveProgress    readerCommands.SaveProgressHandler
	getProgress     readerQueries.GetProgressHandler
	addBookmark     readerCommands.AddBookmarkHandler
	getBookmarks    readerQueries.GetBookmarksHandler
	removeBookmark  readerCommands.RemoveBookmarkHandler
}

func NewReaderController(
	saveProgress readerCommands.SaveProgressHandler,
	getProgress readerQueries.GetProgressHandler,
	addBookmark readerCommands.AddBookmarkHandler,
	getBookmarks readerQueries.GetBookmarksHandler,
	removeBookmark readerCommands.RemoveBookmarkHandler,
) *ReaderController {
	return &ReaderController{
		saveProgress:   saveProgress,
		getProgress:    getProgress,
		addBookmark:    addBookmark,
		getBookmarks:   getBookmarks,
		removeBookmark: removeBookmark,
	}
}
