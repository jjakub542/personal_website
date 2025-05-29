import * as THREE from 'three';
import { PointerLockControls } from 'three/addons/controls/PointerLockControls.js';

const scene = new THREE.Scene();
const camera = new THREE.PerspectiveCamera(75, window.innerWidth/window.innerHeight, 0.1, 1000);
const renderer = new THREE.WebGLRenderer();

document.body.appendChild(renderer.domElement);
renderer.setSize(window.innerWidth, window.innerHeight);

// PointerLockControls
const controls = new PointerLockControls(camera, document.body);
scene.add(controls.getObject());

// Lock pointer on click
document.addEventListener('click', () => {
  controls.lock();
});

// Track key presses
const keys = {};
document.addEventListener('keydown', e => keys[e.code] = true);
document.addEventListener('keyup', e => keys[e.code] = false);


const velocity = new THREE.Vector3();
const direction = new THREE.Vector3();
const speed = 5.0;
const clock = new THREE.Clock();

camera.position.z = 10;
camera.position.y = 5;

const axesHelper = new THREE.AxesHelper(50); // size of the axes
axesHelper.setColors(new THREE.Color(0xff0000), new THREE.Color(0x00ff00), new THREE.Color(0x0000ff));
scene.add(axesHelper);

const geometry = new THREE.PlaneGeometry(10, 10);
const material = new THREE.MeshBasicMaterial({
  color: 0x808080,
  side: THREE.DoubleSide,
  wireframe: false
});
const plane = new THREE.Mesh(geometry, material);
plane.rotation.x = -Math.PI / 2; // or Math.PI / 2 for reverse facing
scene.add(plane);

// Animation loop
function animate() {
  requestAnimationFrame(animate);

  const delta = clock.getDelta();
  direction.set(0, 0, 0);
  if (keys['KeyW']) controls.moveForward(speed * delta);
  if (keys['KeyS']) controls.moveForward(-speed * delta);
  if (keys['KeyA']) direction.x -= 1;
  if (keys['KeyD']) direction.x += 1;
  if (keys['Space']) camera.position.y += speed * delta;
  if (keys['ShiftLeft']) camera.position.y -= speed * delta;

  direction.normalize(); // to ensure consistent speed
  velocity.copy(direction).multiplyScalar(speed * delta);

  controls.moveRight(velocity.x);
  controls.moveForward(velocity.z);

  renderer.render(scene, camera);
  renderer.render(scene, camera);
}
animate();