attribute vec3 aVertexPosition;
attribute vec3 aVertexNormal;
attribute vec4 aVertexColor;

uniform mat4 uModelView;
uniform mat4 uPerspective;
uniform mat3 uNormalMatrix;
uniform vec3 uLightDirection;

varying vec4 vColor;

void main(void) {
  gl_Position = uPerspective * uModelView * vec4(aVertexPosition, 1.0);

  vec3 eyeNormal = uNormalMatrix * aVertexNormal;
  float diffuse = max(0.1, dot(eyeNormal, uLightDirection));
  vColor = aVertexColor * diffuse;
  vColor.a = 1.0;
}
