import * as THREE from "three";
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls.js";

import settings from "./settings";
// import gui, { getController } from "./gui";

import { createRenderer } from "./renderer";
import { createAmbientLight, createDirectionalLight } from "./lights";
import {
  createCube,
  createPlane,
  createCurveTube,
  createWalker,
} from "./objects";
import { createPerspectiveCamera, createOrthographicCamera } from "./cameras";

export default function init() {
  // Scene
  const scene = new THREE.Scene();

  // Renderer
  const renderer = createRenderer();
  settings.parentElement.appendChild(renderer.domElement);

  // Camera
  // const camera = createPerspectiveCamera();
  const camera = createOrthographicCamera();
  const controls = new OrbitControls(camera, renderer.domElement);

  // Lights
  const ambientLight = createAmbientLight();
  scene.add(ambientLight);
  const directionalLight = createDirectionalLight();
  scene.add(directionalLight);

  // Objects
  const plane = createPlane();
  scene.add(plane);

  const cube = createCube();
  scene.add(cube);

  const tube = createCurveTube();
  scene.add(tube);

  const walker = createWalker();
  scene.add(walker);

  // Helpers
  // const axesHelper = new THREE.AxesHelper(5);
  // scene.add(axesHelper);

  // Animation
  const clock = new THREE.Clock();
  const animate = () => {
    const elapsedTime = clock.getElapsedTime();

    // Update objects
    // cube.rotation.x = Math.cos(elapsedTime);
    // cube.rotation.y = Math.sin(elapsedTime);

    walker.update(elapsedTime);

    // Render
    renderer.render(scene, camera);
  };

  // Render
  renderer.render(scene, camera);
  // Animation loop
  renderer.setAnimationLoop(animate);
  // Add logic to controller
  // const animationController = getController("animation");
  // animationController.onChange((value) => {
  //   if (value) {
  //     renderer.setAnimationLoop(animate);
  //   } else {
  //     renderer.setAnimationLoop(null);
  //   }
  // });
}
