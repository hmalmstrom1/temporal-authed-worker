using System;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using IdentityModel.Client;

namespace TemporalWorkerApp.Auth
{
    public class OAuthTokenProvider
    {
        private readonly string _tokenUrl;
        private readonly string _clientId;
        private readonly string _clientSecret;
        private readonly HttpClient _httpClient;
        private TokenResponse _currentToken;
        private DateTime _tokenExpiration;

        public OAuthTokenProvider(string tokenUrl, string clientId, string clientSecret)
        {
            _tokenUrl = tokenUrl;
            _clientId = clientId;
            _clientSecret = clientSecret;
            _httpClient = new HttpClient();
        }

        public async Task<string> GetTokenAsync(CancellationToken cancellationToken = default)
        {
            if (_currentToken != null && DateTime.UtcNow < _tokenExpiration.AddMinutes(-1))
            {
                return _currentToken.AccessToken;
            }

            // Fetch new token
            var tokenResponse = await _httpClient.RequestClientCredentialsTokenAsync(new ClientCredentialsTokenRequest
            {
                Address = _tokenUrl,
                ClientId = _clientId,
                ClientSecret = _clientSecret
            }, cancellationToken);

            if (tokenResponse.IsError)
            {
                throw new Exception($"Failed to retrieve token: {tokenResponse.Error}");
            }

            _currentToken = tokenResponse;
            _tokenExpiration = DateTime.UtcNow.AddSeconds(tokenResponse.ExpiresIn);

            return _currentToken.AccessToken;
        }
    }
}
