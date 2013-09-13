attribute vec4 aStarPosition;

uniform mat4 uPMatrix;
uniform mat4 uVMatrix;

varying vec4 vColor;

void main(void) {

  gl_Position = uPMatrix * uVMatrix * vec4(aStarPosition.xyz, 1.0);
  gl_PointSize = aStarPosition.w;

  vColor = vec4(1.0, 1.0, 1.0, 1.0);
}
