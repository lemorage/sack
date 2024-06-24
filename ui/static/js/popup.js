function togglePopup(popupId) {
  const popup = document.getElementById(popupId);
  if (popup.style.display === 'block') {
      popup.style.display = 'none';
  } else {
      popup.style.display = 'block';
  }
}

function toggleMessage(selector) {
  const popup = document.querySelector(selector);
  if (popup.style.display === 'none' || popup.style.display === '') {
      popup.style.display = 'block';
  } else {
      popup.style.display = 'none';
  }
}

document.addEventListener("DOMContentLoaded", function() {
  const messageBubble = document.querySelector('.message-bubble');

  // Load saved message from localStorage
  const savedMessage = localStorage.getItem('messageBubbleContent');
  if (savedMessage) {
    messageBubble.innerText = savedMessage;
  }

  // Save message to localStorage on input
  messageBubble.addEventListener('input', function() {
    localStorage.setItem('messageBubbleContent', messageBubble.innerText);
  });
});
