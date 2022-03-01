#ifdef GL_ES
  #define MEDIUMP mediump
  precision MEDIUMP float;
#else
  #define MEDIUMP
#endif

uniform vec3 u_LightPos;
uniform vec4 u_Color;
varying vec3 v_Position;
varying vec3 v_Normal;

void main() {
    vec3 diff = u_LightPos - v_Position;
    vec3 lightVector = normalize(diff);
    float diffuse = (dot(v_Normal, lightVector) + 1.0) / 2.0;
    gl_FragColor.rgb = u_Color.rgb * diffuse;
    gl_FragColor.a = u_Color.a;
}