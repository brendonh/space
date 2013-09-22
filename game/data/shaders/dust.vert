//attribute vec3 aPosition;

uniform mat4 uPerspective;
uniform mat4 uView;

uniform vec3 uBasePosition;
uniform vec3 uCenterPosition;

varying vec4 vColor;
  
float rand(vec2 co){
    return fract(sin(dot(co.xy ,vec2(12.9898,78.233))) * 43758.5453);
}

float DUST_BOX = 10.0;

void main(void) {

  float boxMod = mod(((3 * uBasePosition.x) +
                      (3 * uBasePosition.y) +
                      (uBasePosition.z)) / 100.0,
                     1.0);
    
  float rand0 = float(gl_InstanceID) / 100.0;

  float rand1 = rand(vec2(rand0, boxMod));
  float rand2 = rand(vec2(rand0, rand1));
  float rand3 = rand(vec2(rand0, rand2));

  vec4 pos = vec4( uBasePosition.x + rand1 * DUST_BOX,
                   uBasePosition.y + rand2 * DUST_BOX,
                   uBasePosition.z + rand3 * DUST_BOX,
                   1.0 );

  gl_Position = uPerspective * uView * pos;
  gl_PointSize = rand(vec2(rand3, 0.0)) * 3.0;

  float alpha = max(0.0, length(uCenterPosition.xy - pos.xy) / DUST_BOX);
  //vColor = vec4(1.0, 1.0, 1.0, 1 - alpha);
  vColor = vec4(1.0, 1.0, 1.0, 1.0 - (alpha * 0.0001));
}
