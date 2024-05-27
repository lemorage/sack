function togglePopup(popupId) {
  const popup = document.getElementById(popupId);
  if (popup.style.display === 'block') {
      popup.style.display = 'none';
  } else {
      popup.style.display = 'block';
  }
}
