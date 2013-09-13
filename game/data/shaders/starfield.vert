attribute vec4 aStarPosition;

uniform mat4 uPMatrix;
uniform mat4 uVMatrix;

varying vec4 vColor;

void main(void) {

  vec4 basePosition = uPMatrix * uVMatrix * vec4(aStarPosition.xyz, 1.0);
  
  float w = basePosition.w;

  vec4 ndcPos = basePosition / w;

  float x = mod(ndcPos.x + 1.0, 2.0) - 1.0;
  float y = mod(ndcPos.y + 1.0, 2.0) - 1.0;

  gl_Position = vec4(x *w, y * w, basePosition.z, basePosition.w);
  gl_PointSize = aStarPosition.w;

  vColor = vec4(1.0, 1.0, 0.99, 1.0);
}
