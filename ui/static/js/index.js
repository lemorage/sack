const imageContainer = document.getElementById('image-container');
const imagePath = '/static/'; // Path to your .webp images
const numberOfImages = 1;
const imageSize = 200; // Size of images (assuming square)
const [gridCols, gridRows] = getGrid(numberOfImages);

// Ensure that the container is large enough to fit the grid
imageContainer.style.position = 'relative';
imageContainer.style.width = `${gridCols * imageSize}px`;
imageContainer.style.height = `${gridRows * imageSize}px`;

for (let i = 1; i <= numberOfImages; i++) {
  const img = document.createElement('img');
  const folderName = `obj${i}`;

  img.src = `${imagePath}${folderName}/object${i}.webp`;
  img.alt = `Image ${i}`;
  img.style.position = 'absolute';
  img.style.width = `${imageSize}px`;
  img.style.height = `${imageSize}px`;

  // Calculate the grid position
  const col = (i - 1) % gridCols;
  const row = Math.floor((i - 1) / gridCols);

  img.style.left = `${col * imageSize}px`;
  img.style.top = `${row * imageSize}px`;

  img.addEventListener('click', () => {
      window.location.href = `/model${i}`;
  });

  imageContainer.appendChild(img);
}

function getDivisors(n) {
    const res = [];
    
    for (let i = 2; i < Math.floor(n / i); ++i) {
        if (n % i === 0) {
            res.push(i);
            if (i !== Math.floor(n / i)) {
                res.push(Math.floor(n / i));
            }
        }
    }
    
    res.sort((a, b) => a - b);
    return res;
}

function getGrid(num) {
    if (num <= 5) return [num, 1];

    const divisors = getDivisors(num);

    if (divisors.length == 0) {
        return getGrid(num - 1);
    }

    const col = divisors[divisors.length / 2];
    const row = divisors[divisors.length / 2 - 1];

    return [col, row];
}
