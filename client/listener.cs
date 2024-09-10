using System;
using System.Net.Sockets;
using System.Text;

class Listener
{
    private TcpClient _client;
    private NetworkStream _stream;

    public Listener(int port)
    {
        _client = new TcpClient("127.0.0.1", port);
        _stream = _client.GetStream();
    }

    public void Listen()
    {
        byte[] buffer = new byte[1024];
        int bytesRead;

        while ((bytesRead = _stream.Read(buffer, 0, buffer.Length)) != 0)
        {
            string message = Encoding.ASCII.GetString(buffer, 0, bytesRead);
            ServerMessages.Display(message);
        }
    }
}