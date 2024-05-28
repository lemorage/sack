const effectsViewer = document.querySelector("model-viewer#transformer");
const outlineEffect = effectsViewer.querySelector("outline-effect");

document.querySelector("#outline").addEventListener('change', (e) => {
  outlineEffect.blendMode = e.target.checked ? 'default' : 'skip';
});
