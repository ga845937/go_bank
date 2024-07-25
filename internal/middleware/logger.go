package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
)

func Logger() gin.HandlerFunc {
	if _, err := os.Stat(os.Getenv("LOG_PATH")); os.IsNotExist(err) {
		os.MkdirAll(os.Getenv("LOG_PATH"), 0755)
	}
	fileName := fmt.Sprintf("go_bank_%s.log", time.Now().Format("2006-01-02"))
	logFileName := path.Join(os.Getenv("LOG_PATH"), fileName)

	writeFile := zapcore.AddSync(&lumberjack.Logger{
		Filename:  logFileName,
		MaxSize:   10,
		MaxAge:    3,
		LocalTime: true,
		Compress:  true,
	})
	writeStdout := zapcore.AddSync(os.Stdout)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.NewMultiWriteSyncer(writeFile, writeStdout),
		zap.InfoLevel,
	)

	logger = zap.New(core)
	defer logger.Sync()

	return handle_logger()
}

func handle_logger() gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		query := handle_query(context.Request.URL.RawQuery)
		body := handle_body(context)
		param := handle_param(context.Params)

		var traceID string
		tx := apm.TransactionFromContext(context.Request.Context())
		if tx != nil {
			traceID = tx.TraceContext().Trace.String()
		}

		logger.Info("request",
			zap.String("program", "golang"),
			zap.String("project", "go_bank"),
			zap.String("env", os.Getenv("GOLANG_ENV")),
			zap.String("service", strings.Split(path, "/")[1]),
			zap.String("traceID", traceID),
			zap.String("clientIP", context.ClientIP()),
			zap.String("requestMethod", context.Request.Method),
			zap.String("requestURI", context.Request.RequestURI),
			zap.Any("responseHeader", context.Request.Header),
			zap.Any("requestQuery", query),
			zap.Any("requestBody", body),
			zap.Any("requestParam", param),
		)

		response := handle_response(context)
		logger.Info("response",
			zap.String("program", "golang"),
			zap.String("project", "go_bank"),
			zap.String("env", os.Getenv("GOLANG_ENV")),
			zap.String("service", strings.Split(path, "/")[1]),
			zap.String("traceID", traceID),
			zap.String("clientIP", context.ClientIP()),
			zap.String("requestMethod", context.Request.Method),
			zap.String("requestURI", context.Request.RequestURI),
			zap.Int("httpStatus", context.Writer.Status()),
			zap.Any("responseHeader", context.Writer.Header()),
			zap.Any("responseData", response),
			zap.Any("error", context.Errors.JSON()),
		)
	}
}

func handle_query(rawQuery string) map[string]string {
	queryMap, err := url.ParseQuery(rawQuery)
	if err != nil {
		log.Println("Error Parse Query")
		return nil
	}

	simplifiedQueryMap := make(map[string]string)
	for key, values := range queryMap {
		if len(values) > 0 {
			simplifiedQueryMap[key] = values[0]
		}
	}
	return simplifiedQueryMap
}

func handle_body(context *gin.Context) interface{} {
	bodyByte, _ := io.ReadAll(context.Request.Body)
	var bodyJson map[string]interface{}
	var body interface{}

	if context.Request.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		formData, err := url.ParseQuery(string(bodyByte))
		if err == nil {
			formJson := make(map[string]interface{})
			for key, values := range formData {
				if len(values) > 0 {
					formJson[key] = values[0]
				} else {
					formJson[key] = nil
				}
			}
			body = formJson
		} else {
			body = string(bodyByte)
		}
	} else {
		if err := json.Unmarshal(bodyByte, &bodyJson); err != nil {
			body = string(bodyByte)
		} else {
			body = bodyJson
		}
	}

	context.Request.Body = io.NopCloser(bytes.NewBuffer(bodyByte))

	return body
}

func handle_param(params gin.Params) map[string]string {
	paramsJson := make(map[string]string)
	for _, param := range params {
		paramsJson[param.Key] = param.Value
	}
	return paramsJson
}

type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func handle_response(context *gin.Context) interface{} {
	responseWriter := &ResponseWriter{
		body:           &bytes.Buffer{},
		ResponseWriter: context.Writer,
	}
	context.Writer = responseWriter

	context.Next()

	data := responseWriter.body.String()
	responseJson := make(map[string]interface{})
	json.Unmarshal([]byte(data), &responseJson)

	return responseJson
}
