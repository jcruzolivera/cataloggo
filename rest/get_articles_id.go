package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/cataloggo/article"
	"github.com/nmarsollier/cataloggo/rest/engine"
)

//	@Summary		Obtener un articulo
//	@Description	Obtener un articulo
//	@Tags			Catalogo
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"bearer {token}"
//	@Param			articleId		path		string				true	"ID de articlo"
//	@Success		200				{object}	article.ArticleData	"Articulo"
//	@Failure		400				{object}	errs.ValidationErr	"Bad Request"
//	@Failure		401				{object}	engine.ErrorData	"Unauthorized"
//	@Failure		404				{object}	engine.ErrorData	"Not Found"
//	@Failure		500				{object}	engine.ErrorData	"Internal Server Error"
//	@Router			/v1/articles/:articleId [get]
//
// Obtener un articulo
func init() {
	engine.Router().GET(
		"/v1/articles/:articleId",
		getArticle,
	)
}

func getArticle(c *gin.Context) {
	articleId := c.Param("articleId")

	ctx := engine.GinCtx(c)
	result, err := article.FindById(articleId, ctx...)
	if err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.JSON(200, result)
}
