package main

import (
    "context"
    "fmt"
    "github.com/ShiinaOrez/RainbowCat/crawlers/crad"
    "github.com/ShiinaOrez/kylin"
    _const "github.com/ShiinaOrez/kylin/const"
    "github.com/ShiinaOrez/kylin/crawler"
    "github.com/ShiinaOrez/kylin/interceptor"
    "github.com/ShiinaOrez/kylin/param"
)

type ParamInterceptor struct {}

func(i *ParamInterceptor) Run(ctx context.Context) context.Context {
    if kw := ctx.Value("keyword").(string); kw == "" {
        ctx = context.WithValue(ctx, "break", i.GetID())
    }
    return ctx
}

func (i *ParamInterceptor) GetID() string {
    return "keyword-interceptor"
}

type CradCrawler struct {
    crawler.BaseCrawler
}

func (ic CradCrawler) GetID() string {
    return "CRAD-Crawler"
}

func main() {
    var (
        k            kylin.Kylin             = kylin.NewKylin()
        i            interceptor.Interceptor = &ParamInterceptor{}
        cradCrawler  crawler.Crawler         = &CradCrawler{}
    )

    err := k.AddInputInterceptor(&i, "tail")
    if err != nil {
        k.GetLogger().Fatal(err.Error())
        return
    }
    cradCrawler.SetProc(crad.Proc)

    err = k.RegisterCrawler(&cradCrawler)
    if err != nil {
        k.GetLogger().Fatal(err.Error())
        return
    }
    kw := "计算机"
    p := param.NewJSONParam(fmt.Sprintf(`{"content": {"keyword": "%s", "page": 2}}`, kw))

    ch := k.StartOn(p)
    defer k.Stop()

    select {
    case result := <-ch:
        switch result {
        case _const.Success:
            k.GetLogger().Info("Success")
        case _const.Failed:
            k.GetLogger().Info("Failed")
        }
    }
}