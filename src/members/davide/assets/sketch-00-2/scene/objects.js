import * as THREE from "three";

import settings from "./settings";
// import gui from "./gui";
import { getRandomFloat, getRandomElement } from "./utils/random";

export function createCube() {
  const geometry = new THREE.BoxGeometry(1, 1, 1);
  const material = new THREE.MeshStandardMaterial({
    color: settings.color1,
  });
  const cube = new THREE.Mesh(geometry, material);

  // const folder = gui.addFolder("Cube").close();
  // folder.addColor(material, "color").name("Color");

  return cube;
}

export function createPlane() {
  const scaleFactor = 5;
  let aspect =
    settings.parentElement.clientWidth / settings.parentElement.clientHeight;
  let scaleX = scaleFactor * aspect;
  let scaleY = scaleFactor * 1;

  const geometry = new THREE.PlaneGeometry(1, 1);
  const material = new THREE.MeshStandardMaterial({
    color: settings.color2,
    side: THREE.DoubleSide,
  });
  const plane = new THREE.Mesh(geometry, material);
  plane.scale.set(scaleX, scaleY, 1);
  plane.position.z = -2;

  // Handle window resize
  window.addEventListener("resize", () => {
    aspect =
      settings.parentElement.clientWidth / settings.parentElement.clientHeight;
    scaleX = scaleFactor * aspect;
    scaleY = scaleFactor;
    plane.scale.set(scaleX, scaleY, 1);
  });

  // const folder = gui.addFolder("Plane").close();
  // folder.addColor(material, "color").name("Color");

  return plane;

  // Info
  // https://discourse.threejs.org/t/which-method-is-right-to-update-geometry-of-mesh/48676
}

export function createCurveTube() {
  const points = [];
  for (let x = -5; x < 5; x = x + 1) {
    let y = getRandomFloat(-1, 1);
    points.push(new THREE.Vector3(x, y, 0));
  }

  const curve = new THREE.CatmullRomCurve3(points);
  const geometry = new THREE.TubeGeometry(curve, 100, 0.02, 10, false);
  const material = new THREE.MeshStandardMaterial({
    color: settings.color1,
  });
  const tube = new THREE.Mesh(geometry, material);

  // const folder = gui.addFolder("Tube").close();
  // folder.addColor(material, "color").name("Color");

  return tube;
}

export function createWalker() {
  const geometry = new THREE.SphereGeometry(0.25, 16, 16);
  const material = new THREE.MeshStandardMaterial({
    color: settings.color1,
  });
  const walker = new THREE.Mesh(geometry, material);
  walker.position.set(0, 0, 0);

  // const folder = gui.addFolder("Walker").close();
  // folder.addColor(material, "color").name("Color");

  // Movement

  const speed = 2;
  const boundary = 2.5;
  let angle = getRandomFloat(0, 2 * Math.PI);
  let lastTime = 0;

  walker.update = function (time) {
    // Circular movement
    // walker.position.x = Math.sin(time) * 1.5;
    // walker.position.y = Math.cos(time) * 1.5;

    const deltaTime = Math.min(time - lastTime, 0.1); // Cap this to avoid large deltas when pausing animation loop
    lastTime = time;

    // Add a small variation to the angle to create a curved trajectory
    angle += getRandomFloat(-0.1, 0.1);

    // Update position based on angle
    walker.position.x += Math.cos(angle) * deltaTime * speed;
    walker.position.y += Math.sin(angle) * deltaTime * speed;
    // Add z movement
    walker.position.z = Math.sin(time * 2) * 3;
    // Constrain within boundaries
    walker.position.x = Math.max(
      -boundary,
      Math.min(boundary, walker.position.x)
    );
    walker.position.y = Math.max(
      -boundary,
      Math.min(boundary, walker.position.y)
    );

    // Constrain within boundaries and pick a new direction if needed
    if (
      walker.position.x >= boundary ||
      walker.position.x <= -boundary ||
      walker.position.y >= boundary ||
      walker.position.y <= -boundary
    ) {
      angle = getRandomFloat(0, 2 * Math.PI);
    }

    lastTime = time;
  };

  return walker;
}
