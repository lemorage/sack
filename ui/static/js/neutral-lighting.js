const modelViewer = document.querySelector('#transformer');
const checkbox = document.querySelector('#neutral');

checkbox.addEventListener('change', () => {
  modelViewer.environmentImage = checkbox.checked ? '' : 'legacy';
});
