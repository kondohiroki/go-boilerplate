package exception

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/bytedance/sonic"
)

func Test_createFixedExceptionErrors(t *testing.T) {
	errItem := &ExceptionError{
		Message:      "unknown error but i love you",
		Type:         ERROR_TYPE_BAD_REQUEST,
		ErrorSubcode: SUBCODE_CANNOT_RUN_BATCH_DAILY,
	}

	t.Run("test single common error method implementation", func(t *testing.T) {
		expected := "unknown error but i love you"
		actual := errItem.Error()
		if expected != actual {
			t.Errorf("err msg not equal, expected: %s but got: %s", expected, actual)
		}
	})

	type args struct {
		httpStatusCode int
		t              errorType
		esc            errorSubcode
		m              string
	}

	tests := []struct {
		name string
		args args
		want *ExceptionErrors
	}{
		{
			name: "validation error (name is non-sense)",
			args: args{
				httpStatusCode: 501,
				t:              ERROR_TYPE_BAD_REQUEST,
				esc:            SUBCODE_CANNOT_RUN_BATCH_DAILY,
				m:              "unknown error but i love you",
			},
			want: &ExceptionErrors{
				HttpStatusCode: 501,
				GlobalMessage:  "unknown error but i love you",
				ErrItems: []*ExceptionError{
					errItem,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createFixedExceptionErrors(tt.args.httpStatusCode, tt.args.t, tt.args.esc, tt.args.m)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createFixedExceptionErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ExceptionErrors(t *testing.T) {
	cErrs := createFixedExceptionErrors(
		http.StatusInternalServerError,
		ERROR_TYPE_UNKNOWN_ERROR,
		SUBCODE_UNKNOWN_ERROR,
		"test error",
	)

	t.Run("test error method implementation", func(t *testing.T) {
		expectedMsg := "test error"
		actualMsg := cErrs.Error()
		if expectedMsg != actualMsg {
			t.Errorf("error message not equal, expected: %s but got: %s", expectedMsg, actualMsg)
		}
	})

	t.Run("test JSON marshaller implementation", func(t *testing.T) {
		expected := []byte(`[{"message": "test error","type": "UnknownError","errorsubcode":10800}]`)
		actual, err := sonic.Marshal(cErrs)
		if err != nil {
			t.Errorf("marshal JSON must not error")
		}
		if reflect.DeepEqual(expected, actual) {
			t.Errorf("unexpected JSON marshal result, expected: %s but got: %s", expected, actual)
		}
	})
}
