package checkout

import (
	"context"
	"fmt"

	"github.com/Tiger-Coders/tigerlily-bff/internal/helper"
	"github.com/Tiger-Coders/tigerlily-bff/internal/manager/grpc_client"
	"github.com/Tiger-Coders/tigerlily-bff/internal/pkg/logger"
	"github.com/Tiger-Coders/tigerlily-payment/api/rpc"
)

type CheckoutService struct {
	GRPCClient *grpc_client.GRPCClient
	logs       logger.Logger
}

// Inits a new CheckoutService instance ...
func NewCheckoutService(grpc *grpc_client.GRPCClient) *CheckoutService {
	return &CheckoutService{
		logs:       *logger.NewLogger(),
		GRPCClient: grpc,
	}
}

// Checkout method for CheckoutService. Inits a new GRPC clien and invokes strategy algo based on payment type...
func (srv CheckoutService) Checkout(ctx context.Context, req *rpc.CheckoutReq) (resp *rpc.CheckoutResp, err error) {
	defer srv.GRPCClient.Conn.Close()

	srv.GRPCClient.SetStrategy(grpc_client.NewGRPCCheckoutClient(srv.GRPCClient.Conn))
	res, resErr := srv.GRPCClient.Strategy.Execute(ctx, req.PaymentType, req)
	if resErr != nil {
		srv.logs.ErrorLogger.Printf("[SERVICE] Error getting response from RPC via strategy : %+v\n", resErr)
		return nil, resErr
	}
	fmt.Println("CHECKOUT SERVICE BEFORE TRANSFORMING")
	resp, err = helper.TransformCheckoutresp(res)
	if err != nil {
		srv.logs.ErrorLogger.Printf("[SERVICE] Error transforming response to proper format for checkout service: %+v\n", err)
	}
	resp = &rpc.CheckoutResp{
		Success:   true,
		Message:   "Checkout Success",
		StatusUrl: resp.StatusUrl,
	}
	return
}
