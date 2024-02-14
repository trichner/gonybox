include <BOSL2/std.scad>
include <BOSL2/rounding.scad>

$fn=32;
module buttons(topwall){
  grid_copies(15,4)
    linear_extrude(topwall){ 
        polygon(round_corners(square(11, center = true),"circle", r=1.5 ));
    };
  grid_copies(30,2) cylinder(h=topwall-2, r=1);
}

module poti(wall){
  throughwall = 2;
  indent_out = 2.5;
  indent_in = wall - throughwall - indent_out;
  
  linear_extrude(indent_in){
      circle(r=10); fwd(20) left(10) square(20);
  };
  cylinder(h=throughwall + indent_in, r=3.5);
  up(indent_in + throughwall) linear_extrude(indent_out){ circle(r=10); };
}

module speaker(wall){
  diameter = 45;
  linear_extrude(wall) {
    round2d(1){
      difference(){
        circle(d = diameter);
        zrot(90) ycopies(6, 7) rect([diameter,3]);
      }
    }
  }
}

module interior(){
    topwall = 5;
    bottomwall = 5;
    trelligo_depth = 16;
    size = 100;
    wall = 10;
    bodyheight = size - trelligo_depth - topwall - bottomwall;
    translate([wall/2, wall/2,0])
      cube([size-wall,size-wall,bottomwall])
        attach(TOP, CTR+BOTTOM)
          cube([size-wall*2,size-wall*2,bodyheight]) {
            attach(TOP, CTR+BOTTOM) 
              cube([61,61,trelligo_depth])
                attach(TOP) buttons(topwall);
            up(trelligo_depth/2) attach(LEFT) poti(wall);
            up(trelligo_depth/2) attach(RIGHT) speaker(wall);
          }
}


module box(){
  difference() {
    cuboid(100, rounding=5, anchor=BOT+LEFT+FRONT,teardrop=true);
    interior();
  }
}

module speaker_in(){
  cube([55,55, 4], anchor=CTR){
    grid_copies(42, 2) cylinder(h=10, r=0.5);
    cylinder(h=4, r=53/2)
      cylinder(h=25, r1=48/2, r2 = 44/2);
  }
}
module insert(){
  topwall = 5;
  bottomwall = 5;
  trelligo_depth = 13;
  size = 100;
  wall = 10;
  socket=4;
  bodyheight = size - trelligo_depth - topwall - bottomwall;
  back(wall/2) cube([size-wall, size-wall, 5]) {
    attach(TOP) difference(){
        cube([size-2*wall,size-2*wall,socket], anchor=CTR);
        cube([size-2*wall-socket*2,size-2*wall -socket*2,socket], anchor=CTR);
      }
    attach(TOP) 
      right(size/2-wall) 
        diff() {
          cube([8, size-2*wall,bodyheight],anchor=BOTTOM+RIGHT)
            attach(RIGHT) 
              force_tag("remove")
                back(trelligo_depth/2)
                  zflip() speaker_in();
        }
    }
}

module trellis_support(){
  width = 79;
  
  difference(){
    cube([width,width,5]){
      attach(TOP,CTR+BOTTOM) ycopies(30, 2) cylinder(h=8,r1=3,r2=2);
      attach(TOP,CTR+BOTTOM) right(15) cylinder(h=8,r1=3,r2=2);
      attach(TOP,CTR+BOTTOM) grid_copies(54, 2) cube([6,6,8]);
    };
    translate([width/2-7,width/2,0]) cylinder(h=10, r=10);
    translate([width/2,width/2,0])
      grid_copies(30, 2)
        cylinder(h=10, r=1.5);
    translate([width/2,width/2,0]) // corner screw holes
      grid_copies(70, 2)
        cylinder(h=10, r=1.5);
  }
}

module sample_top(){
  width = 80;
  difference(){
    //cube([width,width,21]);
    cuboid([width,width,21], rounding=5, anchor=BOT+LEFT+FRONT,teardrop=true, edges=[TOP, "Z"]);
    translate([(width-61)/2,(width-61)/2,0]) 
      cube([61,61,16])
        attach(TOP) buttons(5);
    translate([width/2,width/2,0]) // corner screw holes
      grid_copies(70, 2)
        cylinder(h=10, r=1);
  }
}

box();
//left(150) interior();

left(100+40) insert();

//fwd(100) speaker_in();

fwd(100) left(100) trellis_support();

translate([-200,-100,21]) zflip() sample_top();

/*
diff() cube(100, center = true) {
    force_tag("remove") down(10) attach(TOP) buttons();
    tag("remove") attach(BOTTOM) down(45) cube([62,62,100], center = true);
    tag("remove") attach(BOTTOM) down(30) cube([80,80,100], center = true);
}
*/