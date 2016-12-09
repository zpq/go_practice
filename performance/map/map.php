<?php
ini_set("memory_limit", "2048M");
$arr = [];
$l = 100000;

for($i=0;$i<$l;$i++) {
    $arr[] = $i;
}
echo "down\n";

echo "begin\n";
$start = microtime(true);
for($i=0;$i<$l;$i++) {
    if (isset($arr[$i])) {

    }
}
echo microtime(true) - $start;
echo "\nend\n";

$start = microtime(true);
foreach($arr as $v) {
    foreach ($arr as $vv) {
        if ($v == $vv) {
            break;
        }
    }
}
echo microtime(true) - $start;

