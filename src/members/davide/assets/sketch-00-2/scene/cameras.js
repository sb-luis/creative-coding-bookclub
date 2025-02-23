import * as THREE from "three";
import settings from "./settings";

// import gui from "./gui";

export function createPerspectiveCamera() {
  const aspect =
    settings.parentElement.clientWidth / settings.parentElement.clientHeight;
  const camera = new THREE.PerspectiveCamera(75, aspect, 0.1, 100);
  camera.position.z = 5;

  // Handle parent element resize
  window.addEventListener("resize", () => {
    camera.aspect =
      settings.parentElement.clientWidth / settings.parentElement.clientHeight;
    camera.updateProjectionMatrix();
  });

  return camera;
}

export function createOrthographicCamera() {
  const aspectRatio =
    settings.parentElement.clientWidth / settings.parentElement.clientHeight;
  const frustumSize = 3;

  const camera = new THREE.OrthographicCamera(
    -frustumSize * aspectRatio,
    frustumSize * aspectRatio,
    frustumSize,
    -frustumSize,
    0.1,
    100
  );

  camera.position.z = 5;

  // Handle parent element resize
  window.addEventListener("resize", () => {
    const aspectRatio =
      settings.parentElement.clientWidth / settings.parentElement.clientHeight;
    camera.left = -frustumSize * aspectRatio;
    camera.right = frustumSize * aspectRatio;
    camera.top = frustumSize;
    camera.bottom = -frustumSize;
    camera.updateProjectionMatrix();
  });

  // GUI

  // const obj = {
  //   resetCamera: function () {
  //     camera.rotation.set(0, 0, 0);
  //     camera.zoom = 1;
  //     camera.position.set(0, 0, 5);
  //     camera.updateProjectionMatrix();
  //   },
  //   fromBottom: function () {
  //     camera.rotation.set(Math.PI / 2, 0, 0);
  //     camera.zoom = 1;
  //     camera.position.set(0, -5, 0);
  //     camera.updateProjectionMatrix();
  //   },
  // };
  // const folder = gui.addFolder("Camera").close();
  // // add a button to reset the camera position
  // folder.add(obj, "resetCamera").name("Front (Reset)");
  // folder.add(obj, "fromBottom").name("Bottom");

  return camera;
}
