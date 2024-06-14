function colorControl() {
  const modelViewerColor = document.querySelector("model-viewer#transformer");

  document.querySelector('#color-controls').addEventListener('click', (event) => {
    const colorString = event.target.dataset.color;
    const [material] = modelViewerColor.model.materials;
    material.pbrMetallicRoughness.setBaseColorFactor(colorString);
  });
}

// Call the function to ensure the controls are initialized when the script is loaded
document.addEventListener('DOMContentLoaded', () => {
  colorControl();
});
