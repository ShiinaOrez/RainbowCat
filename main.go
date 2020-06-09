package main

import (
    "context"
    "fmt"
    "github.com/ShiinaOrez/RainbowCat/crawlers/crad"
    "github.com/ShiinaOrez/kylin"
    _const "github.com/ShiinaOrez/kylin/const"
    "github.com/ShiinaOrez/kylin/crawler"
    "github.com/ShiinaOrez/kylin/interceptor"
    "github.com/ShiinaOrez/kylin/logger"
    "github.com/ShiinaOrez/kylin/param"
    "github.com/ShiinaOrez/kylin/render"
    "os"
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

func main() {
    var (
        k            kylin.Kylin             = kylin.NewKylin()
        i            interceptor.Interceptor = &ParamInterceptor{}
        cradCrawler  crawler.Crawler         = &crawler.BaseCrawler{ID: "CRAD-Crawler"}
    )

    err := k.AddInputInterceptor(&i, "tail")
    if err != nil {
        logger.GetLogger(nil).Fatal(err.Error())
        return
    }
    cradCrawler.SetProc(crad.Proc)

    err = k.RegisterCrawlerWithRender(&cradCrawler, render.FileRender{})
    if err != nil {
        logger.GetLogger(nil).Fatal(err.Error())
        return
    }
    kw := "计算机"
    path, _ := os.Getwd()
    p := param.NewJSONParam(fmt.Sprintf(`{"content": {"keyword": "%s", "page": 1, "path": "%s"}}`, kw, path))

    ch := k.StartOn(p)
    defer k.Stop()

    select {
    case result := <-ch:
        switch result {
        case _const.Success:
            logger.GetLogger(nil).Info("Success")
        case _const.Failed:
            logger.GetLogger(nil).Info("Failed")
        }
    }
}