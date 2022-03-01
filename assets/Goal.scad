$fn=16;
s=.99;

difference() {
    intersection() {
        cube([s,s,s],center=true);
        rotate([45,0,45])
        cube([s,s,s],center=true);
        rotate([0,45,45])
        cube([s,s,s],center=true);
    }
    intersection() {
        t=s*3/4;
        rotate([45,0,0])
        cube([t,t,t],center=true);
        rotate([0,45,0])
        cube([t,t,t],center=true);
        rotate([0,0,45])
        cube([t,t,t],center=true);
    }
}
