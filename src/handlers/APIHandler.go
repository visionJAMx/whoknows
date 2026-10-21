package handlers

import (

	 "github.com/gin-gonic/gin"
	 "net/http"

)


func Register(ctx *gin.Context){

	ctx.String(http.StatusAccepted, ("h"))


}



/*
/api/search


/api/weather



/api/login


/api/logout



/api/register
*/
