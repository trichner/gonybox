
include <box.scad>


box();
left(250) interior();

left(100+40) insert();

//fwd(100) speaker_in();

fwd(100) left(100) trellis_support();

translate([-130,-100,25]) yrot(180) sample_top();
