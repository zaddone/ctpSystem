package route
import(
	"github.com/zaddone/ctpSystem/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init(){
	Router := gin.Default()
	Router.Static("/"+config.Conf.Static,"./"+config.Conf.Static)
	Router.LoadHTMLGlob(config.Conf.Templates+"/*")
	Router.GET("/",func(c *gin.Context){
		c.HTML(http.StatusOK,"index.tmpl",nil)
	})
	Router.GET("/trun",func(c *gin.Context){
		words := c.DefaultQuery("word","")
		c.JSON(http.StatusOK,gin.H{"msg":"Success","word":words})
	})
	Router.Run(config.Conf.Port)

}
