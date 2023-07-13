<?php
/**
 *
 *
 * Created by PhpStorm
 * User: Alex
 * Date: 2023-07-10 22:22
 */
declare(strict_types=1);

require_once __DIR__ . '/SimpleRoundRobin.php';
require_once __DIR__ . '/WeightedRoundRobin.php';
require_once __DIR__ . '/SmoothWeightedRoundRobin.php';

class RobinTest
{

    public function testSimpleRoundRobin()
    {
        $servers = [
            '192.168.0.1',
            '192.168.0.2',
            '192.168.0.3',
        ];
        $balancer = new SimpleRoundRobin($servers);

        $nodes = [];
        for ($i = 1; $i <= 20; $i++) {
            $p = $balancer->getPeer();
            $nodes[$i] = $p;
            echo "i={$i} host={$p} \r\n";
        }

        // output is:
        // i=1 host=192.168.0.1
        // i=2 host=192.168.0.2
        // i=3 host=192.168.0.3
        // i=4 host=192.168.0.1
        // i=5 host=192.168.0.2
        // i=6 host=192.168.0.3
        // i=7 host=192.168.0.1
        // i=8 host=192.168.0.2
        // i=9 host=192.168.0.3
        // i=10 host=192.168.0.1
        // i=11 host=192.168.0.2
        // i=12 host=192.168.0.3
        // i=13 host=192.168.0.1
        // i=14 host=192.168.0.2
        // i=15 host=192.168.0.3
        // i=16 host=192.168.0.1
        // i=17 host=192.168.0.2
        // i=18 host=192.168.0.3
        // i=19 host=192.168.0.1
        // i=20 host=192.168.0.2

        //   '192.168.0.1' => 7,
        //   '192.168.0.2' => 7,
        //   '192.168.0.3' => 6,
        var_export(array_count_values($nodes));

        // 可见，每一个节点都会被调度，但是没法通过权重进行控制

    }

    public function testWeightedRoundRobin()
    {
        $servers = [
            [
                'host' => '192.168.0.1',
                'weight' => 5,
            ],
            [
                'host' => '192.168.0.2',
                'weight' => 1,
            ],
            [
                'host' => '192.168.0.3',
                'weight' => 1,
            ],
        ];

        $balancer = new WeightedRoundRobin($servers);

        $nodes = [];
        for ($i = 1; $i <= 20; $i++) {
            $p = $balancer->getPeer();
            $nodes[$i] = $p['host'];
            echo "i={$i} host={$p['host']} weight={$p['weight']} \r\n";
        }

        // output is:
        // i=1 host=192.168.0.1 weight=5
        // i=2 host=192.168.0.1 weight=5
        // i=3 host=192.168.0.1 weight=5
        // i=4 host=192.168.0.1 weight=5
        // i=5 host=192.168.0.1 weight=5
        // i=6 host=192.168.0.2 weight=1
        // i=7 host=192.168.0.3 weight=1
        // i=8 host=192.168.0.1 weight=5
        // i=9 host=192.168.0.1 weight=5
        // i=10 host=192.168.0.1 weight=5
        // i=11 host=192.168.0.1 weight=5
        // i=12 host=192.168.0.1 weight=5
        // i=13 host=192.168.0.2 weight=1
        // i=14 host=192.168.0.3 weight=1
        // i=15 host=192.168.0.1 weight=5
        // i=16 host=192.168.0.1 weight=5
        // i=17 host=192.168.0.1 weight=5
        // i=18 host=192.168.0.1 weight=5
        // i=19 host=192.168.0.1 weight=5
        // i=20 host=192.168.0.2 weight=1

        //   '192.168.0.1' => 15,
        //   '192.168.0.2' => 3,
        //   '192.168.0.3' => 2,
        var_export(array_count_values($nodes));

        // 可见，虽然会按照一定的权重比例进行轮询，每一个节点也会被调度，但是节点分布不太均匀，尤其是权重偏差越大，调度越不均匀。
    }

    public function testSmoothWeightedRoundRobin()
    {
        $servers = [
            [
                'host' => '192.168.0.1',
                'weight' => 5,
                'current_weight' => 0,
                'effective_weight' => 5,
            ],
            [
                'host' => '192.168.0.2',
                'weight' => 1,
                'current_weight' => 0,
                'effective_weight' => 1,
            ],
            [
                'host' => '192.168.0.3',
                'weight' => 1,
                'current_weight' => 0,
                'effective_weight' => 1,
            ],
        ];

        $balancer = new SmoothWeightedRoundRobin($servers);

        $nodes = [];
        for ($i = 1; $i <= 20; $i++) {
            $p = $balancer->getPeer();
            $nodes[$i] = $p['host'];
            echo "i={$i} host={$p['host']} weight={$p['weight']} current_weight={$p['current_weight']} effective_weight={$p['effective_weight']} \r\n";
        }

        // output is:
        // i=1 host=192.168.0.1 weight=5 current_weight=-2 effective_weight=5
        // i=2 host=192.168.0.1 weight=5 current_weight=-4 effective_weight=5
        // i=3 host=192.168.0.2 weight=1 current_weight=-4 effective_weight=1
        // i=4 host=192.168.0.1 weight=5 current_weight=-1 effective_weight=5
        // i=5 host=192.168.0.3 weight=1 current_weight=-2 effective_weight=1
        // i=6 host=192.168.0.1 weight=5 current_weight=2 effective_weight=5
        // i=7 host=192.168.0.1 weight=5 current_weight=0 effective_weight=5
        // i=8 host=192.168.0.1 weight=5 current_weight=-2 effective_weight=5
        // i=9 host=192.168.0.1 weight=5 current_weight=-4 effective_weight=5
        // i=10 host=192.168.0.2 weight=1 current_weight=-4 effective_weight=1
        // i=11 host=192.168.0.1 weight=5 current_weight=-1 effective_weight=5
        // i=12 host=192.168.0.3 weight=1 current_weight=-2 effective_weight=1
        // i=13 host=192.168.0.1 weight=5 current_weight=2 effective_weight=5
        // i=14 host=192.168.0.1 weight=5 current_weight=0 effective_weight=5
        // i=15 host=192.168.0.1 weight=5 current_weight=-2 effective_weight=5
        // i=16 host=192.168.0.1 weight=5 current_weight=-4 effective_weight=5
        // i=17 host=192.168.0.2 weight=1 current_weight=-4 effective_weight=1
        // i=18 host=192.168.0.1 weight=5 current_weight=-1 effective_weight=5
        // i=19 host=192.168.0.3 weight=1 current_weight=-2 effective_weight=1
        // i=20 host=192.168.0.1 weight=5 current_weight=2 effective_weight=5

        //   '192.168.0.1' => 14,
        //   '192.168.0.2' => 3,
        //   '192.168.0.3' => 3,
        var_export(array_count_values($nodes));

        // 可见，按照了一定权重比例进行了轮询，且每一个节点的调度分布也比较均匀

    }

}

$t = new RobinTest();
// $t->testSimpleRoundRobin();
// $t->testWeightedRoundRobin();
// $t->testSmoothWeightedRoundRobin();