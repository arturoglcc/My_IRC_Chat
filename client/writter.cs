using System;
using System.Net.Sockets;
using System.Text;

class Writer
{
    private TcpClient _client;
    private NetworkStream _stream;

    public Writer(int port)
    {
        _client = new TcpClient("127.0.0.1", port);
        _stream = _client.GetStream();
    }

    public void SendMessage(string message)
    {
        byte[] data = Encoding.ASCII.GetBytes(message);
        _stream.Write(data, 0, data.Length);
    }
}