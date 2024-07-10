import * as THREE from 'three';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import Stats from 'three/addons/libs/stats.module.js';

let container, stats;
let camera, scene, renderer;
let mesh, mixer;
let pageCount;

const radius = 600;
let theta = 0;
let prevTime = Date.now();

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

  const loader = new GLTFLoader();

  for (let i = 1; i <= pageCount; ++i) {
    loader.load( `/static/obj${i}/object${i}.glb`, function (gltf) {
      dim = Math.floor(Math.random() * 10 + 5);

      mesh = gltf.scene.children[0];
      mesh.scale.set(dim, dim, dim);
      scene.add(mesh);

      mixer = new THREE.AnimationMixer(mesh);
    });
  }

  renderer = new THREE.WebGLRenderer();
  renderer.setPixelRatio(window.devicePixelRatio);
  renderer.setSize(window.innerWidth, window.innerHeight);
  renderer.setAnimationLoop(animate);

  container.appendChild(renderer.domElement);

  stats = new Stats();
  container.appendChild(stats.dom);

  window.addEventListener('resize', onWindowResize);
}

function onWindowResize() {
  camera.aspect = window.innerWidth / window.innerHeight;
  camera.updateProjectionMatrix();

  renderer.setSize(window.innerWidth, window.innerHeight);
}

function animate() {
  render();
  stats.update();
}

function render() {
  theta += 0.1;

  camera.position.x = radius * Math.sin(THREE.MathUtils.degToRad(theta));
  camera.position.z = radius * Math.cos(THREE.MathUtils.degToRad(theta));

  camera.lookAt(0, 150, 0);

  if (mixer) {
    const time = Date.now();
    mixer.update((time - prevTime) * 0.001);
    prevTime = time;
  }

  renderer.render(scene, camera);
}
