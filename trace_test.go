package google_cloud_trace

import (
	"net/http"
	"testing"
)

func TestDeconstructXCloudTraceContext(t *testing.T) {
	tcs := []struct {
		v       string
		traceID string
		spanID  string
		sampled bool
	}{
		{
			v:       "105445aa7843bc8bf206b120001000/000000000000004a;o=1",
			traceID: "projects//traces/105445aa7843bc8bf206b120001000",
			spanID:  "000000000000004a",
			sampled: true,
		},
		{
			v:       "105445aa7843bc8bf206b120001000/000000000000004a;o=0",
			traceID: "projects//traces/105445aa7843bc8bf206b120001000",
			spanID:  "000000000000004a",
			sampled: false,
		},
		{
			v:       "/0;o=1",
			traceID: "",
			spanID:  "",
			sampled: true,
		},
		{
			v:       "/;o=1",
			traceID: "",
			spanID:  "",
			sampled: true,
		},
		{
			v:       "105445aa7843bc8bf206b120001000/;o=1",
			traceID: "projects//traces/105445aa7843bc8bf206b120001000",
			spanID:  "",
			sampled: true,
		},
		{
			v:       "105445aa7843bc8bf206b120001000/0",
			traceID: "projects//traces/105445aa7843bc8bf206b120001000",
			spanID:  "",
			sampled: false,
		},
		{
			v:       "105445aa7843bc8bf206b120001000",
			traceID: "projects//traces/105445aa7843bc8bf206b120001000",
			spanID:  "",
			sampled: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.v, func(t *testing.T) {
			r, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			r.Header.Set("X-Cloud-Trace-Context", tc.v)

			fields := Fields(r)

			var tr, sp string
			var sm bool
			for _, f := range fields {
				switch f.Key {
				case traceKey:
					tr = f.String
				case spanKey:
					sp = f.String
				case traceSampledKey:
					sm = f.Integer == 1
				}
			}

			//tr, sp, sm := deconstructXCloudTraceContext(ctx)
			//tr, sp, sm := deconstructXCloudTraceContextValue(tc.v)
			if tr != tc.traceID {
				t.Error("invalid trace id", tr, tc.traceID)
			}
			if sp != tc.spanID {
				t.Error("invalid span id", sp, tc.spanID)
			}
			if sm != tc.sampled {
				t.Error("invalid sampled", sm, tc.sampled)
			}
		})
	}
}
