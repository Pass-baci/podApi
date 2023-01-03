package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Pass-baci/common"
	"github.com/Pass-baci/pod/proto/pod"
	"github.com/Pass-baci/podApi/proto/podApi"
	"strconv"
)

type PodApi struct {
	PodService pod.PodService
}

func (p *PodApi) FindPodById(ctx context.Context, req *podApi.Request, rsp *podApi.Response) error {
	common.Info("接收到 FindPodById 的请求")

	var (
		podIdPair *podApi.Pair
		ok        bool
		err       error
	)
	if podIdPair, ok = req.Get["pod_id"]; !ok {
		rsp.StatusCode = 400
		return errors.New("参数异常")
	}

	var podId int64
	if podId, err = strconv.ParseInt(podIdPair.Values[0], 10, 64); err != nil {
		rsp.StatusCode = 400
		return errors.New("参数异常")
	}

	var podInfo *pod.PodInfo
	if podInfo, err = p.PodService.FindPodByID(ctx, &pod.PodId{Id: podId}); err != nil {
		rsp.StatusCode = 500
		common.Errorf("FindPodById 查询Pod失败 ID: %d, err: %s \n", podId, err.Error())
		return errors.New("查询Pod失败")
	}

	rsp.StatusCode = 200
	b, _ := json.Marshal(podInfo)
	rsp.Body = string(b)

	return nil
}

func (p *PodApi) AddPod(ctx context.Context, req *podApi.Request, rsp *podApi.Response) error {
	common.Info("接收到 AddPod 的请求")

	rsp.StatusCode = 200
	rsp.Body = "AddPod 的请求处理成功"

	return nil
}

func (p *PodApi) DeletePodById(ctx context.Context, req *podApi.Request, rsp *podApi.Response) error {
	common.Info("接收到 DeletePodById 的请求")

	rsp.StatusCode = 200
	rsp.Body = "DeletePodById 的请求处理成功"

	return nil
}

func (p *PodApi) UpdatePod(ctx context.Context, req *podApi.Request, rsp *podApi.Response) error {
	common.Info("接收到 UpdatePod 的请求")

	rsp.StatusCode = 200
	rsp.Body = "UpdatePod 的请求处理成功"

	return nil
}

//默认接口
func (p *PodApi) Call(ctx context.Context, req *podApi.Request, rsp *podApi.Response) error {
	common.Info("接收到 Call 的请求")

	rsp.StatusCode = 200
	rsp.Body = "Call 的请求处理成功"

	return nil
}
