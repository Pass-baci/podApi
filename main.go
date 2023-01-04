package main

import (
	"fmt"
	"github.com/Pass-baci/common"
	"github.com/Pass-baci/pod/proto/pod"
	handler2 "github.com/Pass-baci/podApi/handler"
	hystrix2 "github.com/Pass-baci/podApi/plugin/hystrix"
	"github.com/Pass-baci/podApi/proto/podApi"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/go-micro/plugins/v3/registry/consul"
	ratelimiter "github.com/go-micro/plugins/v3/wrapper/ratelimiter/uber"
	"github.com/go-micro/plugins/v3/wrapper/select/roundrobin"
	microOpentracing "github.com/go-micro/plugins/v3/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"net"
	"net/http"
)

var (
	// 注册中心配置
	consulHost       = "192.168.230.135"
	consulPort int64 = 8500
	// 熔断器
	hystrixPort = "9093"
	// 监控端口
	prometheusPort = 9193
)

func main() {
	// 配置consul
	consulClient := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			fmt.Sprintf("%s:%d", consulHost, consulPort),
		}
	})

	// 配置中心consul
	conf, err := common.GetConsulConfig(consulHost, consulPort, "/micro/config")
	if err != nil {
		common.Fatalf("连接配置中心consul失败 err: %s \n", err.Error())
	}
	common.Info("添加配置中心consul成功")

	// 使用配置中心获取Jaeger配置
	var jaegerConfig = &common.JaegerConfig{}
	if jaegerConfig, err = common.GetJaegerConfigFromConsul(conf, "jaeger-podApi"); err != nil {
		common.Fatalf("配置中心获取Jaeger配置失败 err: %s \n", err.Error())
	}
	common.Info("配置中心获取Jaeger配置成功")

	// 添加链路追踪
	tracer, closer, err := common.NewTracer(jaegerConfig.ServiceName, jaegerConfig.Address)
	if err != nil {
		common.Fatalf("添加链路追踪失败 err: %s \n", err.Error())
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	common.Info("添加链路追踪成功")

	// 添加熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	common.Info("添加熔断器成功")

	// 启动监听程序
	go func() {
		// http://ip:port/turbine/turbine.stream
		// 看板访问地址：http://192.168.230.135:9002/hystrix
		if err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", hystrixPort), hystrixStreamHandler); err != nil {
			common.Errorf("启动监听程序失败 err: %s \n", err.Error())
		}
	}()

	// 添加日志中心
	// 需要把程序日志打入到日志文件中
	//common.Info("添加日志中心成功")
	// 在程序代码中添加filebeat.yml文件
	// 启动filebeat 启动命令 ./filebeat -e -c filebeat.yml

	common.PrometheusBoot(prometheusPort)

	// Create service
	srv := micro.NewService(
		micro.Server(server.NewServer(func(options *server.Options) {
			options.Advertise = "192.168.230.135:8082"
		})),
		micro.Name("go.micro.api.podApi"),
		micro.Version("latest"),
		micro.Address(":8082"),
		// 添加consul
		micro.Registry(consulClient),
		// 添加链路追踪
		micro.WrapHandler(microOpentracing.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(microOpentracing.NewClientWrapper(opentracing.GlobalTracer())),
		// 只作为客户端的时候起作用
		micro.WrapClient(hystrix2.NewClientHystrixWrapper()),
		// 添加限流
		micro.WrapHandler(ratelimiter.NewHandlerWrapper(1000)),
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	srv.Init()

	podService := pod.NewPodService("go.micro.service.pod", srv.Client())

	if err = podApi.RegisterPodApiHandler(srv.Server(), &handler2.PodApi{podService}); err != nil {
		common.Fatalf("RegisterPodApiHandler失败 err: %s \n", err.Error())
	}

	// Run service
	if err = srv.Run(); err != nil {
		common.Fatalf("Run service失败 err: %s \n", err.Error())
	}
}
