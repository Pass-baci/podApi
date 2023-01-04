package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pass-baci/common"
	"github.com/Pass-baci/pod/proto/pod"
	"github.com/Pass-baci/podApi/plugin/from"
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
	fmt.Printf("%+v", *req)

	addPodInfo := &pod.PodInfo{}

	podPortPair, ok := req.Get["pod_port"]
	if ok {
		for _, v := range podPortPair.Values {
			i, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				common.Errorf("AddPod 创建Pod失败, err: %s \n", err.Error())
				return errors.New("创建Pod失败")
			}
			port := &pod.PodPort{
				ContainerPort: int32(i),
				Protocol:      "TCP",
			}
			addPodInfo.PodPort = append(addPodInfo.PodPort, port)
		}
	}

	from.FromToPodStruct(req.Get, addPodInfo)
	fmt.Printf("%+v", *addPodInfo)
	response, err := p.PodService.AddPod(ctx, addPodInfo)
	if err != nil {
		common.Errorf("AddPod 创建Pod失败, err: %s \n", err.Error())
		return errors.New("创建Pod失败")
	}

	rsp.StatusCode = 200
	b, _ := json.Marshal(response)
	rsp.Body = string(b)

	return nil
}

func (p *PodApi) DeletePodById(ctx context.Context, req *podApi.Request, rsp *podApi.Response) error {
	common.Info("接收到 DeletePodById 的请求")

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

	var podInfo *pod.Response
	if podInfo, err = p.PodService.DeletePod(ctx, &pod.PodId{Id: podId}); err != nil {
		rsp.StatusCode = 500
		common.Errorf("DeletePodById 删除Pod失败 ID: %d, err: %s \n", podId, err.Error())
		return errors.New("删除Pod失败")
	}

	rsp.StatusCode = 200
	b, _ := json.Marshal(podInfo)
	rsp.Body = string(b)

	return nil
}

func (p *PodApi) UpdatePod(ctx context.Context, req *podApi.Request, rsp *podApi.Response) error {
	common.Info("接收到 UpdatePod 的请求")

	addPodInfo := &pod.PodInfo{}

	podPortPair, ok := req.Get["pod_port"]
	if ok {
		for _, v := range podPortPair.Values {
			i, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				common.Errorf("UpdatePod 更新Pod失败, err: %s \n", err.Error())
				return errors.New("更新Pod失败")
			}
			port := &pod.PodPort{
				ContainerPort: int32(i),
				Protocol:      "TCP",
			}
			addPodInfo.PodPort = append(addPodInfo.PodPort, port)
		}
	}

	from.FromToPodStruct(req.Get, addPodInfo)
	response, err := p.PodService.UpdatePod(ctx, addPodInfo)
	if err != nil {
		common.Errorf("UpdatePod 更新Pod失败, err: %s \n", err.Error())
		return errors.New("更新Pod失败")
	}

	rsp.StatusCode = 200
	b, _ := json.Marshal(response)
	rsp.Body = string(b)

	return nil

	return nil
}

//默认接口
func (p *PodApi) Call(ctx context.Context, req *podApi.Request, rsp *podApi.Response) error {
	common.Info("接收到 Call 的请求")

	allPod, err := p.PodService.FindAllPod(ctx, &pod.FindAll{})
	if err != nil {
		common.Errorf("Call 查询All Pod失败, err: %s \n", err.Error())
		return errors.New("查询All Pod失败")
	}

	rsp.StatusCode = 200
	b, _ := json.Marshal(allPod)
	rsp.Body = string(b)

	return nil
}
