using System;
using System.Threading.Tasks;
using Temporalio.Client.Interceptors;

namespace TemporalWorkerApp.Auth
{
    public class OAuthInterceptor : IClientInterceptor
    {
        private readonly OAuthTokenProvider _tokenProvider;

        public OAuthInterceptor(OAuthTokenProvider tokenProvider)
        {
            _tokenProvider = tokenProvider;
        }

        public ClientOutboundInterceptor InterceptClient(ClientOutboundInterceptor next)
        {
            return new OAuthHeaderInterceptor(next, _tokenProvider);
        }
    }

    public class OAuthHeaderInterceptor : ClientOutboundInterceptor
    {
        private readonly OAuthTokenProvider _tokenProvider;

        public OAuthHeaderInterceptor(ClientOutboundInterceptor next, OAuthTokenProvider tokenProvider) : base(next)
        {
            _tokenProvider = tokenProvider;
        }

        // TODO: Override methods to inject headers.
        // Example:
        // public override async Task<WorkflowExecution> StartWorkflowAsync(StartWorkflowInput input)
        // {
        //     // Add headers to input.RpcOptions
        //     return await base.StartWorkflowAsync(input);
        // }
    }
}
