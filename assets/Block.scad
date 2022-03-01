s=.99;

intersection() {
    cube([s,s,s],center=true);
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
