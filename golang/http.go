package main

import (
	"bytes"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/deepch/vdk/format/ts"

	"github.com/gin-gonic/gin"
)

func serveHTTP() {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

	gin.SetMode(gin.DebugMode)
	router.LoadHTMLGlob("web/templates/*")
	router.GET("/", func(c *gin.Context) {
		fi, all := Config.list()
		sort.Strings(all)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"port":     Config.Server.HTTPPort,
			"suuid":    fi,
			"suuidMap": all,
			"version":  time.Now().String(),
		})
	})
	router.GET("/player/:suuid", func(c *gin.Context) {
		_, all := Config.list()
		sort.Strings(all)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"port":     Config.Server.HTTPPort,
			"suuid":    c.Param("suuid"),
			"suuidMap": all,
			"version":  time.Now().String(),
		})
	})

	router.GET("/play/hls/:suuid/index.m3u8", PlayHLS)
	router.GET("/play/hls/:suuid/segment/:seq/file.ts", PlayHLSTS)
	//添加调试接口列出所有 suuid
	router.GET("/debug/streams", func(c *gin.Context) {
        _, all := Config.list()
        c.JSON(200, gin.H{
            "streams": all,
        })
    })
	router.StaticFS("/static", http.Dir("web/static"))
	err := router.Run(Config.Server.HTTPPort)
	if err != nil {
		log.Fatalln(err)
	}
}

func PlayHLS(c *gin.Context) {
	suuid := c.Param("suuid")
	log.Println("收到播放请求 suuid =", suuid)
	if !Config.ext(suuid) {
	    log.Println("流不存在：", suuid)
    	c.Status(http.StatusNotFound)
    	return
    }
	Config.RunIFNotRun(suuid)
	for i := 0; i < 40; i++ {
		index, seq, err := Config.StreamHLSm3u8(suuid)
		if err != nil {
			log.Println(err)
			return
		}
		if seq >= 6 {
			_, err := c.Writer.Write([]byte(index))
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
		log.Println("Play list not ready wait or try update page")
		time.Sleep(1 * time.Second)
	}
}

//PlayHLSTS send client ts segment
func PlayHLSTS(c *gin.Context) {
	suuid := c.Param("suuid")
	if !Config.ext(suuid) {
		return
	}
	codecs := Config.coGe(c.Param("suuid"))
	if codecs == nil {
		return
	}
	outfile := bytes.NewBuffer([]byte{})
	Muxer := ts.NewMuxer(outfile)
	err := Muxer.WriteHeader(codecs)
	if err != nil {
		log.Println(err)
		return
	}
	Muxer.PaddingToMakeCounterCont = true
	seqData, err := Config.StreamHLSTS(c.Param("suuid"), stringToInt(c.Param("seq")))
	if err != nil {
		log.Println(err)
		return
	}
	if len(seqData) == 0 {
		log.Println(err)
		return
	}
	for _, v := range seqData {
		v.CompositionTime = 1
		err = Muxer.WritePacket(*v)
		if err != nil {
			log.Println(err)
			return
		}
	}
	err = Muxer.WriteTrailer()
	if err != nil {
		log.Println(err)
		return
	}
	_, err = c.Writer.Write(outfile.Bytes())
	if err != nil {
		log.Println(err)
		return
	}
}
