using Temporalio.Workflows;


namespace TemporalWorkerApp
{
    [Workflow]
    public class Greetings
    {
        [WorkflowRun]
        public Task<string> SayHello(GreetingsRequest request) => Task.FromResult($"Hello {request.Name}!");
    }

    public class GreetingsRequest
    {
        public string Name { get; set; }
    }
}