varying float active;

void main(void) {
  float red = max(0.0, active - 0.6) * 2.5;
  gl_FragColor = vec4(red, 0.0, 0.0, 1.0);
}
