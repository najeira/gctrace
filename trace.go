package gctrace

import (
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

const (
	traceKey        = "logging.googleapis.com/trace"
	spanKey         = "logging.googleapis.com/spanId"
	traceSampledKey = "logging.googleapis.com/trace_sampled"
	headerName      = "X-Cloud-Trace-Context"
)

func Fields(r *http.Request) (fields []zap.Field) {
	ctx := r.Context()
	traceID, spanID, traceSampled := deconstructXCloudTraceContext(r)
	if len(traceID) > 0 {
		if projectID := projectIDWithContext(ctx); len(projectID) >= 0 {
			traceIDValue := fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)
			fields = append(fields, zap.String(traceKey, traceIDValue))
		}
	}
	if len(spanID) > 0 {
		fields = append(fields, zap.String(spanKey, spanID))
	}
	if traceSampled {
		fields = append(fields, zap.Bool(traceSampledKey, traceSampled))
	}
	return fields
}

func deconstructXCloudTraceContext(r *http.Request) (traceID, spanID string, traceSampled bool) {
	v := r.Header.Get(headerName)
	return deconstructXCloudTraceContextValue(v)
}

// https://cloud.google.com/appengine/docs/standard/go/writing-application-logs#writing_app_logs
// アプリログのエントリをリクエストログに関連付ける場合は、
// 構造化アプリログのエントリにリクエストのトレース ID を含める必要があります。
// トレース ID は X-Cloud-Trace-Context リクエスト ヘッダーから抽出できます。
// 構造化ログエントリで、ID を logging.googleapis.com/trace という名前のフィールドに書き込みます。
// X-Cloud-Trace-Context ヘッダーの詳細については、リクエストを強制的にトレースするをご覧ください。
//
// https://cloud.google.com/trace/docs/setup#force-trace
// X-Cloud-Trace-Context: TRACE_ID/SPAN_ID;o=TRACE_TRUE
// TRACE_ID は、128 ビットの番号を表す 32 文字の 16 進数値です。
// リクエストを束ねるつもりがないのであれば、リクエスト間で一意の値にする必要があります。これには UUID を使用できます。
// SPAN_ID は、（符号なしの）スパン ID の 10 進表現です。
// これはランダムに生成され、トレースで一意である必要があります。
// 後続のリクエストでは、SPAN_ID を親リクエストのスパン ID に設定します。
// ネストされたトレースの詳細については、TraceSpan（REST、RPC）の説明をご覧ください。
// このリクエストをトレースするには、TRACE_TRUE を 1 に設定する必要があります。
// リクエストをトレースしない場合は 0 を指定します。
func deconstructXCloudTraceContextValue(v string) (traceID, spanID string, traceSampled bool) {
	parts := strings.Split(v, "/")
	if len(parts) <= 0 {
		return
	}
	traceID = parts[0]

	if len(parts) <= 1 {
		return
	}
	parts = strings.Split(parts[1], ";")

	if len(parts) <= 0 {
		return
	}
	spanID = parts[0]
	if spanID == "0" {
		spanID = ""
	}

	if len(parts) <= 1 {
		return
	}
	traceSampled = parts[1] == "o=1"
	return
}
