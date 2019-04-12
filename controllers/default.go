package controllers

import "github.com/astaxie/beego"

/**
@author: agamgn
@date:	2019-04-12
 */

type HomeContraller struct {
	beego.Controller
	
}

func (h *HomeContraller) Get()  {
	h.TplName="index.html"
}