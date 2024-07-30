import * as THREE from 'three';
import Stats from 'three/addons/libs/stats.module.js';
import { LightningStrike } from './jsm/geometries/LightningStrike.js';

let container, stats;
let camera, scene, renderer;
let pageCount;
let particles, octahedron;
let raycaster, mouse;
let lightningStrike, lightningStrikeMesh;

const radius = 3000;
let theta = 0;
let sprites = [];
let isZooming = false;
let lightningVisible = false;
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

  // Add octahedron
  const geometry = new THREE.OctahedronGeometry(90);
  const material = new THREE.MeshPhongMaterial({ 
    color: 0x510896, 
    flatShading: true,
    transparent: true,
    opacity: 0.425
  });
  octahedron = new THREE.Mesh(geometry, material);
  octahedron.position.set(0, 50, 0);
  octahedron.scale.set(0.78, 1.175, 0.825);
  scene.add(octahedron);

  octahedron.userData = { isOctahedron: true };

  // Set up lightning parameters
  scene.userData.lightningColor = 0x55a0ff;
  scene.userData.lightningMaterial = new THREE.MeshBasicMaterial({ 
    color: scene.userData.lightningColor,
    transparent: true,
  });

  scene.userData.rayParams = {
    sourceOffset: new THREE.Vector3(0, radius, 0),
    destOffset: new THREE.Vector3(0, -radius, 0),
    radius0: 4,
    radius1: 4,
    minRadius: 2.5,
    maxIterations: 7,
    isEternal: true,

    timeScale: 0.7,
    propagationTimeFactor: 0.05,
    vanishingTimeFactor: 0.95,
    subrayPeriod: 3.5,
    subrayDutyCycle: 0.6,
    maxSubrayRecursion: 3,
    ramification: 7,
    recursionProbability: 0.6,

    roughness: 0.85,
    straightness: 0.6
  };

  createLightningStrike();
  scene.add(lightningStrikeMesh);
  lightningStrikeMesh.visible = false; // Initially hidden

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

  mouse.x = (event.clientX / window.innerWidth) * 2 - 1;
  mouse.y = -(event.clientY / window.innerHeight) * 2 + 1;

  raycaster.setFromCamera(mouse, camera);

  const intersects = raycaster.intersectObjects([...sprites, octahedron]);

  if (intersects.length > 0) {
    const intersectedObject = intersects[0].object;

    if (intersectedObject.userData.isOctahedron) {
      // Make lightning visible when octahedron is clicked
      createLightBeamEffect();
      zoomIntoObject(intersectedObject);
    } else {
      const num = intersectedObject.userData.objectNum;
      window.location.href = `/model${num}`;
    }
  }
}

function createLightningStrike() {
  lightningStrike = new LightningStrike(scene.userData.rayParams);
  lightningStrikeMesh = new THREE.Mesh(lightningStrike, scene.userData.lightningMaterial);
}

function createLightBeamEffect() {
  return new Promise((resolve) => {
    const duration = 600;
    const startTime = Date.now();

    function animateLightBeam() {
      const elapsedTime = Date.now() - startTime;
      const progress = Math.min(elapsedTime / duration, 1);

      // Update lightning strike
      lightningStrike.update(progress * scene.userData.rayParams.timeScale);
      
      if (!lightningVisible) {
        lightningStrikeMesh.visible = lightningVisible = true;
      }

      if (progress < 1) {
        requestAnimationFrame(animateLightBeam);
      } else {
        resolve();
      }
    }

    animateLightBeam();
  });
}

async function zoomIntoObject(object) {
  isZooming = true;

  // Create light beam effect
  await createLightBeamEffect();

  const zoomDuration = 960;
  const zoomStartTime = Date.now();
  const initialPosition = camera.position.clone();
  const finalPosition = new THREE.Vector3().copy(object.position);
  finalPosition.add(new THREE.Vector3(0, 50, 0));

  function zoom() {
    const elapsedTime = Date.now() - zoomStartTime;
    const t = elapsedTime / zoomDuration;

    if (t < 1) {
      camera.position.lerpVectors(initialPosition, finalPosition, t);
      camera.lookAt(object.position);
      requestAnimationFrame(zoom);
    } else {
      camera.position.copy(finalPosition);
      camera.lookAt(object.position);
      // Optionally, fade out the lightning effect here before changing the page
      fadeLightningOut(() => {
        window.location.href = '/story';
      });
    }
  }

  zoom();
}

function fadeLightningOut(callback) {
  const duration = 500; // 0.5 seconds
  const startTime = Date.now();

  function fadeOut() {
    const elapsedTime = Date.now() - startTime;
    const progress = Math.min(elapsedTime / duration, 1);

    if (progress >= 1) {
      lightningStrikeMesh.visible = lightningVisible = false;
      callback();
    } else {
      requestAnimationFrame(fadeOut);
    }
  }

  fadeOut();
}

function animate() {
  if (!isZooming) {
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

    // Rotate the octahedron
    if (octahedron) {
      octahedron.rotation.x += 0.01;
      octahedron.rotation.y += 0.01;
    }

    // update the lightning strike only when visible
    if (lightningVisible) {
      const time = Date.now() / 1000;
      lightningStrike.update(time * scene.userData.rayParams.timeScale);
    }
  }

  renderer.render(scene, camera);
  stats.update();
}

function isPositionValid(position, newSprite, sprites) {
  for (let i = 0; i < sprites.length; ++i) {
    const sprite = sprites[i];
    const distance = position.distanceTo(sprite.position);
    const minDistance = (newSprite.scale.x / 2) + (sprite.scale.x / 2);
    if (distance < minDistance) {
      return false;
    }
  }
  return true;
}
