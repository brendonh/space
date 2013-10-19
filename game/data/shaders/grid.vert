attribute vec3 aVertexPosition;

uniform mat4 uModelView;
uniform mat4 uPerspective;

uniform bool uActive;
uniform vec2 uActiveCoords;

varying float active;

void main(void) {
  gl_Position = uPerspective * uModelView * vec4(aVertexPosition, 1.0);

  if (uActive && length(aVertexPosition.xy - uActiveCoords) < 1.5) {
    active = 1.0;
  } else {
    active = 0.0;
  }

}
