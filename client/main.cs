using System;

class Program
{
    static void Main(string[] args)
    {
        int port = 1234; // Default port
        if (args.Length > 0 && int.TryParse(args[0], out int userPort))
        {
            port = userPort;
        }

        // Instantiate client
        Client client = new Client(port);
        client.Start();
    }
}