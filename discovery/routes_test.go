package discovery

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func Test_checkPath(t *testing.T) {
	srv := genServerContent(func() string {
		return `{}`
	})
	defer srv.Close()

	srvErr := genServerStatus(func() int {
		return http.StatusNotFound
	})
	defer srvErr.Close()

	type fields struct {
		client   *http.Client
		location string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "location not valid - no connection",
			fields: fields{client: &http.Client{Timeout: time.Millisecond}, location: "example.com:8080"},
			want:   "",
		},
		{
			name:   "location not valid - url not found",
			fields: fields{client: srvErr.Client(), location: url2host(srvErr.URL)},
			want:   "",
		},
		{
			name:   "location is valid",
			fields: fields{client: srv.Client(), location: url2host(srv.URL)},
			want:   fmt.Sprintf("http://%s/swagger.json", url2host(srv.URL)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkPath(tt.fields.client, tt.fields.location)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
