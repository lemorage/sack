document.addEventListener('keydown', function(event) {
    if (event.key === 'ArrowLeft') {
        // Left arrow key pressed
        const prevLink = document.getElementById('prev-model');
        if (prevLink) {
            window.location.href = prevLink.href;
        }
    } else if (event.key === 'ArrowRight') {
        // Right arrow key pressed
        const nextLink = document.getElementById('next-model');
        if (nextLink) {
            window.location.href = nextLink.href;
        }
    }
});
