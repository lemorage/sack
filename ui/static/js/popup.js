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
  const icon = document.getElementById('toolbox-icon');
  const popup = document.getElementById('toolbox-popup');

    // Hide popup when clicking outside
  document.addEventListener('click', (event) => {
    if (!popup.contains(event.target) && event.target !== icon) {
      popup.style.display = 'none';
    }
  });

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
