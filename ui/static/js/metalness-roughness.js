const modelViewerParameters = document.querySelector("model-viewer#transformer");

modelViewerParameters.addEventListener("load", (ev) => {
  let material = modelViewerParameters.model.materials[0];

  let metalnessDisplay = document.querySelector("#metalness-value");
  let roughnessDisplay = document.querySelector("#roughness-value");

  // Set initial values to match the model's state
  document.querySelector('#metalness').value = material.pbrMetallicRoughness.metallicFactor;
  document.querySelector('#roughness').value = material.pbrMetallicRoughness.roughnessFactor;

  // Update display to match initial values
  metalnessDisplay.textContent = material.pbrMetallicRoughness.metallicFactor;
  roughnessDisplay.textContent = material.pbrMetallicRoughness.roughnessFactor;

  // Event listeners for input changes
  document.querySelector('#metalness').addEventListener('input', (event) => {
    material.pbrMetallicRoughness.setMetallicFactor(event.target.value);
    metalnessDisplay.textContent = event.target.value;
  });

  document.querySelector('#roughness').addEventListener('input', (event) => {
    material.pbrMetallicRoughness.setRoughnessFactor(event.target.value);
    roughnessDisplay.textContent = event.target.value;
  });
});
