using System;
using System.Data.Common;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Threading;
using Newtonsoft.Json;
using Npgsql;
using StackExchange.Redis;

namespace Worker {
    public class Program {

        public static int Main(string[] args) {
            try {
                var pgsql = OpenDbConnection("Server=db;Username=postgres;");
                var redis = OpenRedisConnection("redis").GetDatabase();

                var definition = new { quote = "", quote_id = "" };
                while (true) {
                    string json = redis.ListLeftPopAsync("quotes").Result;
                    if (json != null) {
                        var quote = JsonConvert.DeserializeAnonymousType(json, definition);
                        Console.WriteLine($"Processing quote for '{quote.quote}' by '{quote.quote_id}'");
                        Updatequote(pgsql, quote.quote_id, quote.quote);
                    }
                }
            } catch (Exception ex) {
                Console.Error.WriteLine(ex.ToString());
                return 1;
            }
        }

        // Open a connection to the database
        private static NpgsqlConnection OpenDbConnection(string connectionString) {
            NpgsqlConnection connection;
            while (true) {
                try {
                    connection = new NpgsqlConnection(connectionString);
                    connection.Open();
                    break;
                } catch (SocketException) {
                    Console.Error.WriteLine("Waiting for db");
                    Thread.Sleep(1000);
                } catch (DbException) {
                    Console.Error.WriteLine("Waiting for db");
                    Thread.Sleep(1000);
                }
            }

            Console.Error.WriteLine("Connected to db");

            var command = connection.CreateCommand();
            command.CommandText = @"CREATE TABLE IF NOT EXISTS quotes (
                                        id VARCHAR(255) NOT NULL UNIQUE,
                                        quote VARCHAR(255) NOT NULL
                                    )";
            command.ExecuteNonQuery();
            return connection;
        }


        // Open a connection to the redis server
        private static ConnectionMultiplexer OpenRedisConnection(string hostname)
        {
            // Use IP address to workaround:
            // https://github.com/StackExchange/StackExchange.Redis/issues/410
            var ipAddress = GetIp(hostname);
            Console.WriteLine($"Found redis at {ipAddress}");

            // This will **eventually** return a redis connection
            while (true) {
                try {
                    Console.Error.WriteLine("Connected to redis");
                    return ConnectionMultiplexer.Connect(ipAddress);
                } catch (RedisConnectionException) {
                    Console.Error.WriteLine("Waiting for redis");
                    Thread.Sleep(1000);
                }
            }
        }


        // Resolve the IP address of a given hostname
        private static string GetIp(string hostname)
            => Dns.GetHostEntryAsync(hostname)
                .Result
                .AddressList
                .First(a => a.AddressFamily == AddressFamily.InterNetwork)
                .ToString();


        //
        private static void Updatequote(NpgsqlConnection connection, string quoteId, string quote) {
            var command = connection.CreateCommand();
            try {
                command.CommandText = "INSERT INTO quotes (id, quote) VALUES (@id, @quote)";
                command.Parameters.AddWithValue("@id", quoterId);
                command.Parameters.AddWithValue("@quote", quote);
                command.ExecuteNonQuery();
            } catch (DbException) {
                command.CommandText = "UPDATE quotes SET quote = @quote WHERE id = @id";
                command.ExecuteNonQuery();
            } finally {
                command.Dispose();
            }
        }
    }
}
