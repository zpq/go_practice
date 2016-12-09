<?php
// 0, 1, 1, 2, 3, 5, 8, 13, 21...

// function fb($max) {
//     $a = 0;
//     $b = 1;
//     while ($a < $max) {
//         echo $a, ' ';
//         $tmp = $a;
//         $a = $b;
//         $b += $tmp;
//     }
// }

// fb(1000000);

// echo "\n", memory_get_usage(true);



function gen($n) {
    $i = 0;
    while($i < $n) {
        yield $i.":)";
        $i++;
    }
}

// $g = gen(3);
// $g->current();
// echo $g->send(0)."\n";
// echo $g->send(1)."\n";
// echo $g->send(2)."\n";

// foreach($g as $v) {
//     echo $v;
// }

function gen2() {
    $ret = (yield 'yield1');
    var_dump($ret);
    $ret = (yield 'yield2');
    var_dump($ret);
    yield "yield3";
}

$g2 = gen2();
// var_dump($g2->current());
var_dump($g2->send("hello"));
// var_dump($g2->current());
var_dump($g2->send("ret2"));





