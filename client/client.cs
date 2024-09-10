using System;
using System.Threading;

class Client
{
    private int _port;
    private Writer _writer;
    private Listener _listener;

    public Client(int port)
    {
        _port = port;
        _writer = new Writer(_port);
        _listener = new Listener(_port);
    }

    public void Start()
    {
        // Start the listener in a separate thread
        Thread listenerThread = new Thread(_listener.Listen);
        listenerThread.Start();

        // Keep receiving user input and send through the writer
        while (true)
        {
            Console.Write("Enter message: ");
            string input = Console.ReadLine();
            _writer.SendMessage(input);
        }
    }
}