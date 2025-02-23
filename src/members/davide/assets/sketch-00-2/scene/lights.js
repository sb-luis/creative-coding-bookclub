import * as THREE from "three";

export function createAmbientLight() {
  return new THREE.AmbientLight(0xffffff, 1);
}

export function createDirectionalLight() {
  const light = new THREE.DirectionalLight(0xffffff, 1);
  light.position.set(0, 0, 10);

  return light;
}
