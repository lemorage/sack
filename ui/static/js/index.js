import * as THREE from 'three';
import Stats from 'three/addons/libs/stats.module.js';

let container, stats;
let camera, scene, renderer;
let pageCount;
let particles;
let raycaster, mouse;

const radius = 3000;
let theta = 0;
let sprites = [];
let highlightedObject = null;

init();

async function getYamlPagesCount() {
  try {
    const response = await fetch('./config.yaml');
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }

    const yamlText = await response.text();
    const config = jsyaml.load(yamlText);

    return Object.keys(config.Pages).length;
  } catch (error) {
    console.error('Error fetching or parsing YAML file:', error);
    return 0;
  }
}

async function init() {
  container = document.createElement('div');
  document.body.appendChild(container);

  pageCount = await getYamlPagesCount();

  camera = new THREE.PerspectiveCamera(50, window.innerWidth / window.innerHeight, 1, 10000);
  camera.position.y = 300;

  scene = new THREE.Scene();
  scene.background = new THREE.Color(0x000000);

  const light1 = new THREE.DirectionalLight(0xefefff, 5);
  light1.position.set(1, 1, 1).normalize();
  scene.add(light1);

  const light2 = new THREE.DirectionalLight(0xffefef, 5);
  light2.position.set(-1, -1, -1).normalize();
  scene.add(light2);

  for (let i = 1; i <= pageCount; ++i) {
    let map = new THREE.TextureLoader().load(`/static/obj${i}/object${i}.webp`);
    let material = new THREE.SpriteMaterial({ map: map, color: 0xdfcdcd });
    let sprite = new THREE.Sprite(material);

    // Randomly scale the sprite
    let scale = Math.random() * 200 + 323;
    sprite.scale.set(scale, scale, 1);
    sprite.userData = { objectNum: i };

    // Random position with overlap check
    let position;
    do {
      position = new THREE.Vector3(
        (Math.random() - 0.5) * radius,
        (Math.random() - 0.5) * radius,
        (Math.random() - 0.5) * radius
      );
    } while (!isPositionValid(position, sprite, sprites));

    sprite.position.copy(position);
    sprites.push(sprite);
    scene.add(sprite);
  }

  // Create particles
  const particlesGeometry = new THREE.BufferGeometry();
  const particlesCount = 10000;
  const posArray = new Float32Array(particlesCount * 3);

  for (let i = 0; i < particlesCount * 3; i++) {
    posArray[i] = (Math.random() - 0.5) * radius * 2;
  }

  particlesGeometry.setAttribute('position', new THREE.BufferAttribute(posArray, 3));

  const particlesMaterial = new THREE.PointsMaterial({
    size: 1,
    sizeAttenuation: true,
    color: 0xffffff,
    transparent: true,
    opacity: 0.8
  });

  particles = new THREE.Points(particlesGeometry, particlesMaterial);
  scene.add(particles);

  renderer = new THREE.WebGLRenderer();
  renderer.setPixelRatio(window.devicePixelRatio);
  renderer.setSize(window.innerWidth, window.innerHeight);
  renderer.setAnimationLoop(animate);

  container.appendChild(renderer.domElement);

  stats = new Stats();
  container.appendChild(stats.dom);

  // Initialize raycaster and mouse vector
  raycaster = new THREE.Raycaster();
  mouse = new THREE.Vector2();

  window.addEventListener('resize', onWindowResize);
  window.addEventListener('mousemove', onMouseMove);
  window.addEventListener('click', onObjectClick);
}

function onWindowResize() {
  camera.aspect = window.innerWidth / window.innerHeight;
  camera.updateProjectionMatrix();

  renderer.setSize(window.innerWidth, window.innerHeight);
}

function onMouseMove(event) {
  event.preventDefault();

  mouse.x = (event.clientX / window.innerWidth) * 2 - 1;
  mouse.y = -(event.clientY / window.innerHeight) * 2 + 1;

  raycaster.setFromCamera(mouse, camera);
  const intersects = raycaster.intersectObjects(sprites);

  if (intersects.length > 0) {
    const intersectedObject = intersects[0].object;

    if (highlightedObject && highlightedObject !== intersectedObject) {
      highlightedObject.material.color.set(0xdfcdcd);
    }

    // Highlight new object
    intersectedObject.material.color.set(0xfff8e7);
    highlightedObject = intersectedObject;
  } else if (highlightedObject) {
    // Reset highlighted object if no intersections
    highlightedObject.material.color.set(0xdfcdcd);
    highlightedObject = null;
  }
}

function onObjectClick(event) {
  event.preventDefault();

  const mouse = new THREE.Vector2();
  mouse.x = (event.clientX / window.innerWidth) * 2 - 1;
  mouse.y = -(event.clientY / window.innerHeight) * 2 + 1;

  const raycaster = new THREE.Raycaster();
  raycaster.setFromCamera(mouse, camera);

  const intersects = raycaster.intersectObjects(sprites);

  if (intersects.length > 0) {
    const num = intersects[0].object.userData.objectNum;
    window.location.href = `/model${num}`;
  }
}

function animate() {
  theta += 0.1;

  camera.position.x = radius * Math.sin(THREE.MathUtils.degToRad(theta));
  camera.position.z = radius * Math.cos(THREE.MathUtils.degToRad(theta));

  camera.lookAt(0, 150, 0);

  // Move sprites
  sprites.forEach(sprite => {
    sprite.position.x += (Math.random() - 0.5) * 2;
    sprite.position.y += (Math.random() - 0.5) * 2;
    sprite.position.z += (Math.random() - 0.5) * 2;
  });

  // Animate particles
  particles.rotation.y += 0.002;

  // Update the raycaster for hover effects
  raycaster.setFromCamera(mouse, camera);
  const intersects = raycaster.intersectObjects(sprites);

  if (intersects.length > 0) {
    const intersectedObject = intersects[0].object;

    if (highlightedObject && highlightedObject !== intersectedObject) {
      highlightedObject.material.color.set(0xdfcdcd);
    }

    intersectedObject.material.color.set(0xfff8e7);
    highlightedObject = intersectedObject;
  } else if (highlightedObject) {
    highlightedObject.material.color.set(0xdfcdcd);
    highlightedObject = null;
  }

  renderer.render(scene, camera);
  stats.update();
}

function isPositionValid(position, newSprite, sprites) {
  for (let i = 0; i < sprites.length; i++) {
    const sprite = sprites[i];
    const distance = position.distanceTo(sprite.position);
    const minDistance = (newSprite.scale.x / 2) + (sprite.scale.x / 2);
    if (distance < minDistance) {
      return false;
    }
  }
  return true;
}
