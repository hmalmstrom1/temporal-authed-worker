using System;
using System.IO;
using System.Security.Cryptography.X509Certificates;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Temporalio.Client;
using Temporalio.Worker;
using Temporalio.Activities;
using TemporalWorkerApp.Auth;

namespace TemporalWorkerApp
{
    class Program
    {
        static async Task Main(string[] args)
        {
            // Setup logging
            using var loggerFactory = LoggerFactory.Create(builder =>
            {
                builder.AddConsole();
                builder.SetMinimumLevel(LogLevel.Information);
            });
            var logger = loggerFactory.CreateLogger<Program>();

            // Load config from env vars
            var temporalAddress = Environment.GetEnvironmentVariable("TEMPORAL_ADDRESS") ?? "localhost:7233";
            var temporalNamespace = Environment.GetEnvironmentVariable("TEMPORAL_NAMESPACE") ?? "default";
            var oauthClientId = Environment.GetEnvironmentVariable("OAUTH_CLIENT_ID") ?? "client";
            var oauthClientSecret = Environment.GetEnvironmentVariable("OAUTH_CLIENT_SECRET") ?? "secret";
            var oauthTokenUrl = Environment.GetEnvironmentVariable("OAUTH_TOKEN_URL") ?? "http://localhost:8080/token";
            
            var tlsCertPath = Environment.GetEnvironmentVariable("TLS_CERT_PATH");
            var tlsKeyPath = Environment.GetEnvironmentVariable("TLS_KEY_PATH");
            var tlsRootCaPath = Environment.GetEnvironmentVariable("TLS_SERVER_ROOT_CA");

            logger.LogInformation($"Connecting to Temporal at {temporalAddress} Namespace: {temporalNamespace}");

            // Setup OAuth
            var tokenProvider = new OAuthTokenProvider(oauthTokenUrl, oauthClientId, oauthClientSecret);
            var token = await tokenProvider.GetTokenAsync();

	    logger.LogInformation($"Received a token of  {token}");
            // Setup TLS
            TlsOptions tlsOptions = null;
            if (!string.IsNullOrEmpty(tlsRootCaPath))
            {
                var rootCaBytes = await File.ReadAllBytesAsync(tlsRootCaPath);
                tlsOptions = new TlsOptions
                {
                    ServerRootCACert = rootCaBytes
                };

                if (!string.IsNullOrEmpty(tlsCertPath) && !string.IsNullOrEmpty(tlsKeyPath))
                {
                    var certPem = await File.ReadAllBytesAsync(tlsCertPath);
                    var keyPem = await File.ReadAllBytesAsync(tlsKeyPath);
                    
                    tlsOptions.ClientCert = certPem;
                    tlsOptions.ClientPrivateKey = keyPem;
                }
            }

            // Create Client Options with RpcMetadata
            var clientOptions = new TemporalClientConnectOptions(temporalAddress)
            {
                Namespace = temporalNamespace,
                Tls = tlsOptions,
                RpcMetadata = new Dictionary<string, string>
                {
                    { "Authorization", $"Bearer {token}" }
                }
            };

            try
            {
                var client = await TemporalClient.ConnectAsync(clientOptions);
                logger.LogInformation("Successfully connected to Temporal Service!");

                // Create Worker
                using var worker = new TemporalWorker(
                    client,
                    new TemporalWorkerOptions(taskQueue: "my-task-queue")
                    .AddActivity(Hello));

                logger.LogInformation("Starting worker...");
                await worker.ExecuteAsync(CancellationToken.None);
            }
            catch (Exception ex)
            {
                logger.LogError(ex, "Failed to start worker");
            }
        }

        [Activity]
        public static string Hello(string name) => $"Hello {name}!";
    }
}
