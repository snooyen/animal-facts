package admin

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Endpoints struct {
	HandleSMSEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, logger log.Logger) Endpoints {
	endpoints := Endpoints{
		HandleSMSEndpoint: MakeHandleSMSEndpoint(s),
	}

	endpoints.HandleSMSEndpoint = EndpointLoggingMiddleware(log.With(logger, "method", "HandleSMS"))(endpoints.HandleSMSEndpoint)

	return endpoints
}

func MakeHandleSMSEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(handleSMSRequest)
		_, e := s.HandleSMS(ctx, req)
		return handleSMSResponse{Err: e}, nil
	}
}

type handleSMSRequest struct {
	ToCountry     string
	ToState       string
	SmsMessageSid string
	NumMedia      string
	ToCity        string
	FromZip       string
	SmsSid        string
	FromState     string
	SmsStatus     string
	FromCity      string
	Body          string
	FromCountry   string
	To            string
	ToZip         string
	NumSegments   string
	MessageSid    string
	AccountSid    string
	From          string
	ApiVersion    string
}

type handleSMSResponse struct {
	Err error `json:"err,omitempty"`
}

func (r handleSMSResponse) error() error { return r.Err }
