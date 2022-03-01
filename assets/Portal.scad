$fn=6;
s=.99;
h=s/2;
o=s/4+.1;

intersection() {
    difference() {
        cube([s,s,s],center=true);
        translate([0,0,-o])
        cylinder(d1=s,d2=0,h=h,center=true);
        translate([0,0,o])
        cylinder(d2=s,d1=0,h=h,center=true);
        translate([-o,0,0])
        rotate([90,0,90])
        cylinder(d1=s,d2=0,h=h,center=true);
        translate([o,0,0])
        rotate([-90,0,-90])
        cylinder(d2=s,d1=0,h=h,center=true);
        translate([0,-o,0])
        rotate([0,90,90])
        cylinder(d1=s,d2=0,h=h,center=true);
        translate([0,o,0])
        rotate([0,-90,-90])
        cylinder(d2=s,d1=0,h=h,center=true);
    }
    intersection() {
        t=s*4/3;
        rotate([0,0,45])
        cube([t,t,t],center=true);
        rotate([0,45,0])
        cube([t,t,t],center=true);
        rotate([45,0,0])
        cube([t,t,t],center=true);
    }
}
