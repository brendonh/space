attribute vec3 aVertexPosition;

uniform mat4 uModelView;
uniform mat4 uPerspective;

void main(void) {
  gl_Position = uPerspective * uModelView * vec4(aVertexPosition, 1.0);
}
