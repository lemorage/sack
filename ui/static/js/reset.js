function reset() {
  const modelViewer = document.querySelector("model-viewer#transformer");
  const originalColor = [1, 1, 1, 1]; // Assuming the original color is white
  // const originalColor = modelViewer.model.materials[0].pbrMetallicRoughness.baseColorFactor.slice(); // Store original color

  document.querySelector('#color-controls').addEventListener('click', (event) => {
    const colorString = event.target.dataset.color;
    if (colorString) {
      const [material] = modelViewer.model.materials;
      material.pbrMetallicRoughness.setBaseColorFactor(colorString);
    }
  });

  const neutralCheckbox = document.querySelector('#neutral');
  const outlineEffect = modelViewer.querySelector("outline-effect");
  const outlineCheckbox = document.querySelector('#outline');

  document.querySelector('#reset-button').addEventListener('click', () => {
    const [material] = modelViewer.model.materials;
    material.pbrMetallicRoughness.setBaseColorFactor(originalColor);

    neutralCheckbox.checked = !!checkbox.attributes.getNamedItem("checked");
    modelViewer.environmentImage = '';

    outlineCheckbox.checked = !!outlineEffect.attributes.getNamedItem("checked");
    outlineEffect.blendMode = 'skip';

    if (0.5 != material.pbrMetallicRoughness.metallicFactor ||
      0.5 != material.pbrMetallicRoughness.roughnessFactor) {
        window.location.reload();
    }
  });
}

document.addEventListener('DOMContentLoaded', () => {
  reset();
});
