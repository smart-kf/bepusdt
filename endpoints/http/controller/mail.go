package controller

import (
	"github.com/gin-gonic/gin"

	"usdtpay/infr/mail"
)

type SendMailRequest struct {
	From       string `json:"from" binding:"required"`
	To         string `json:"to" binding:"required"`
	HtmlString string `json:"html_string" binding:"required"`
}

func SendMail(ctx *gin.Context) {
	var req SendMailRequest
	if err := ctx.ShouldBind(&req); err != nil {
		sendError(ctx, err)
		return
	}

	err := mail.SendMail(
		ctx.Request.Context(), mail.SendMailObject{
			To:         req.To,
			From:       req.From,
			HtmlString: req.HtmlString,
		},
	)

	if err != nil {
		sendError(ctx, err)
		return
	}

	sendSuccess(ctx, nil)
}
