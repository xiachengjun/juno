//
// @Date: 2020-07-08 11:07:37
// @LastEditors: MEX7
// @LastEditTime: 2020-07-09 10:00:30
// @FilePath: /juno/api/apiv1/proxy/proxy.go
//

package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/douyu/juno/internal/pkg/invoker"
	"github.com/douyu/juno/internal/pkg/packages/contrib/output"
	"github.com/douyu/juno/internal/pkg/service/proxy"
	"github.com/douyu/juno/pkg/constx"
	"github.com/douyu/juno/pkg/model/view"
	"github.com/douyu/juno/pkg/pb"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ProxyPost(c echo.Context) (err error) {
	if c.Request().URL.String() == "/api/v1/resource/node/heartbeat" {
		return NodeHeartBeat(c)
	}
	reqModel := view.ReqHTTPProxy{}
	err = c.Bind(&reqModel)
	if err != nil {
		return output.JSON(c, output.MsgErr, err.Error())
	}

	switch reqModel.Type {
	case "POST":
		xlog.Info("post info", xlog.Any("path", c.Request().URL.String()), xlog.Any("req", reqModel))
		resp, err := invoker.Resty.R().SetBody(reqModel.Params).Post(fmt.Sprintf("http://%s%s", reqModel.Address, reqModel.URL))
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}
		return c.HTMLBlob(http.StatusOK, resp.Body())
	case "GET":
		xlog.Info("get info", xlog.Any("path", c.Request().URL.String()), xlog.Any("req", reqModel))
		resp, err := invoker.Resty.R().SetQueryParams(reqModel.Params).Get(fmt.Sprintf("http://%s%s", reqModel.Address, reqModel.URL))
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}
		return c.HTMLBlob(http.StatusOK, resp.Body())
	}
	return c.String(http.StatusOK, "unsupport type")
}

// NodeHeartBeat 只有可用区的创建功能
func NodeHeartBeat(c echo.Context) error {
	var (
		err error
	)
	reqModel := view.ReqNodeHeartBeat{}
	err = c.Bind(&reqModel)
	if err != nil {
		return output.JSON(c, output.MsgErr, err.Error())
	}

	info, err := json.Marshal(reqModel)
	if err != nil {
		return output.JSON(c, output.MsgErr, err.Error())
	}

	if proxy.StreamStore.IsStreamExist() {
		err = proxy.StreamStore.GetStream().Send(&pb.NotifyResp{
			MsgId: constx.MsgNodeHeartBeatResp,
			Msg:   info,
		})
		if err != nil {
			return output.JSON(c, output.MsgErr, err.Error())
		}
	} else {
		return output.JSON(c, output.MsgErr, "stream is not exist")
	}
	return output.JSON(c, output.MsgOk, "success")
}

//
//// pprofInfo
//func PprofInfo(c echo.Context) error {
//	var err error
//	reqModel := view.ReqHTTPProxy{}
//	err = c.Bind(&reqModel)
//	if err != nil {
//		return output.JSON(c, output.MsgErr, err.Error())
//	}
//	fmt.Println("reqModel", reqModel)
//
//	ip, err := checkPara(reqModel.Params, "ip")
//	if err != nil {
//		return output.JSON(c, output.MsgErr, err.Error())
//	}
//	port, err := checkPara(reqModel.Params, "port")
//	if err != nil {
//		return output.JSON(c, output.MsgErr, err.Error())
//	}
//	pprofType, err := checkPara(reqModel.Params, "fileType")
//	if err != nil {
//		return output.JSON(c, output.MsgErr, err.Error())
//	}
//
//	resp, err := getPprof(ip, port, pprofType)
//
//	if err != nil {
//		return output.JSON(c, output.MsgErr, err.Error())
//	}
//	return output.JSON(c, output.MsgOk, "get pprof success", resp)
//	//return c.HTMLBlob(http.StatusOK, resp)
//}
////
////func checkPara(para map[string]interface{}, tar string) (tarStr string, err error) {
////	tmp, ok := para[tar]
////	if !ok {
////		return tarStr, fmt.Errorf("必须传" + tar)
////	}
////	tarStr, ok = tmp.(string)
////	if !ok {
////		return tarStr, fmt.Errorf("%s必须为string类型", tar)
////	}
////	if tarStr == "" {
////		return tarStr, fmt.Errorf("%s不能为空", tar)
////	}
////	return tarStr, nil
////}
//
//func getPprof(ip, port, pprofType string) (resp []byte, err error) {
//	client := resty.New().SetDebug(false).SetTimeout(60*time.Second).SetHeader("X-JUNO-GOVERN", "juno")
//	url := fmt.Sprintf("http://%s:%s/debug/pprof", ip, port)
//	// 检测接口是否ok
//	if _, err = checkPprof(client, url); err != nil {
//		return
//	}
//	// 耗时比较久
//	if pprofType == "profile" {
//		pprofType = pprofType + "?seconds=30"
//	}
//	url = url + "/" + pprofType
//	// 获取数据
//	if resp, err = checkPprof(client, url); err != nil {
//		return
//	}
//	return
//}

//func checkPprof(client *resty.Client, url string) ([]byte, error) {
//	response, err := client.R().Get(url)
//	if err != nil {
//		return []byte{}, err
//	}
//	return response.Body(), nil
//}
