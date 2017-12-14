<?php
class JsonRPC
{
    private $conn;

    public function __construct($host, $port)
    {
        $this->conn = fsockopen($host, $port, $errno, $errstr, 3);
        if (!$this->conn) {
            var_dump($errstr, $errno);die;
        }
    }

    public function Call($method, $params)
    {
        if (!$this->conn) {
            return false;
        }
        $err = fwrite($this->conn, json_encode(array(
            'method'  => $method,
            'params'  => [$params],
            'id'      => 0,
            'jsonrpc' => '2.0',
        )) . "\n");
        if ($err === false) {
            return false;
        }

        stream_set_timeout($this->conn, 0, 3000);
        // 不知道为什么有时第一次读取不到
        //$i = 0;
        while (true) {
            $line = fgets($this->conn);
            if (empty($line)) {
                continue;
            } else {
                break;
            }
        }
        return json_decode($line, true);
    }
}
/*****----------php cli运行 测试脚本-----------*****/
$client = new JsonRPC("127.0.0.1", 1314);
$a = microtime(true);
$m1 = memory_get_usage();
$r = $client->Call("Excel.ReadExcel", ["FileName" => 'test.xlsx']);
echo microtime(true) - $a;
echo PHP_EOL;
echo (memory_get_usage() - $m1) / 1000;
echo PHP_EOL;
var_dump($r['result']);
