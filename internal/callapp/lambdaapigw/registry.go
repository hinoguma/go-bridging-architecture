package lambdaapigw

type HandlerRegistry interface {
	DepositHandler() LambdaAPIGWHandler
	WithdrawHandler() LambdaAPIGWHandler
}

type MiddlewareRegistry interface {
	AuthenticateBeforeMiddleware() BeforeMiddleware
}
