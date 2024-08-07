<h2>This is HTML</h2>
<?php
namespace All\Tokens;

use Other\Tokens;

require_once "settings.php"
require "settings.php"
include_once "settings.php"
include "settings.php"

class Token extends Tokenb implements TokenR {
    public function __construct(private $name)
    {
    }
}

$a = 2;
$b = $a + 1.55;

$g =new Token("MyName");
$f = clone $g;
eval("echo 'Hellow';");

print "that's it";

?>
<h2>aefaef</h2>
<?php

