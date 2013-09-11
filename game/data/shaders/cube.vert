attribute vec3 aVertexPosition;
attribute vec3 aNormal;
attribute vec4 aVertexColor;

uniform mat4 uMVMatrix;
uniform mat4 uPMatrix;
uniform mat3 uNormalMatrix;
uniform vec3 uLightDirection;

varying vec4 vColor;

void main(void) {
  gl_Position = uPMatrix * uMVMatrix * vec4(aVertexPosition, 1.0);

  vec3 eyeNormal = uNormalMatrix * aNormal;
  float diffuse = max(0.1, dot(eyeNormal, uLightDirection));
  vColor = aVertexColor * diffuse;
  vColor.a = 1.0;
}
