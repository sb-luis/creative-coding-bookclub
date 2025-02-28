import * as THREE from "three";

import settings from "./settings";
// import gui from "./gui";

export function createRenderer() {
  const renderer = new THREE.WebGLRenderer();
  renderer.setSize(
    settings.parentElement.clientWidth,
    settings.parentElement.clientHeight
  );
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
  renderer.setClearColor(new THREE.Color(settings.backgroundColor));
  renderer.antialias = true;

  // Handle window resize
  window.addEventListener("resize", () => {
    renderer.setSize(
      settings.parentElement.clientWidth,
      settings.parentElement.clientHeight
    );
  });

  // GUI
  // const folder = gui.addFolder("Renderer").close();
  // const obj = { backgroundColor: settings.backgroundColor, animation: true };
  // folder
  //   .addColor(obj, "backgroundColor")
  //   .name("Background Color")
  //   .onChange((value) => {
  //     renderer.setClearColor(new THREE.Color(value));
  //   });

  // folder.add(obj, "animation").name("Animation");
  // Logic added in _init to start/stop animation

  return renderer;
}
