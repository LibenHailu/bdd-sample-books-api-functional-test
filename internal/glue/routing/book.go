package routing

import (
	"net/http"

	"github.com/LibenHailu/sample-books/internal/http/rest"
	"github.com/LibenHailu/sample-books/platform/gin"
	ginRouter "github.com/LibenHailu/sample-books/platform/gin"
)

func BookRouting(handler rest.BookHandler) []ginRouter.Router {
	return []gin.Router{
		{
			Method:  http.MethodPost,
			Path:    "books",
			Handler: handler.InsertBook,
		},

		{
			Method:  http.MethodGet,
			Path:    "books",
			Handler: handler.GetBooks,
		},

		{
			Method:  http.MethodGet,
			Path:    "books/:id",
			Handler: handler.GetBook,
		},

		{
			Method:  http.MethodDelete,
			Path:    "books/:id",
			Handler: handler.DeleteBook,
		},
		{
			Method:  http.MethodPatch,
			Path:    "books/:id",
			Handler: handler.UpdateBook,
		},
	}
}
