const width = document.getElementById('graph-container').offsetWidth;
const height = document.getElementById('graph-container').offsetHeight;

let zoom; // Declare zoom globally to access it in the reset function
let brushing = false; // Track if brushing mode is enabled

function createGraph(data) {
  const svg = d3.select('#graph-container').append('svg')
    .attr('width', width)
    .attr('height', height);

  zoom = d3.zoom()
    .scaleExtent([0.5, 10]) // Set the minimum and maximum zoom scale
    .on('zoom', (event) => {
      container.attr('transform', event.transform);
      updateVisibility(event.transform.k);
    });

  svg.call(zoom);

  const container = svg.append('g'); // Add a 'g' element to apply zoom transform

  const simulation = d3.forceSimulation(data.nodes)
    .force("link", d3.forceLink(data.links).id(d => d.id).distance(100))
    .force('charge', d3.forceManyBody().strength(-300))
    .force('center', d3.forceCenter(width / 2, height / 2));

  const link = container.append('g')
    .selectAll('line')
    .data(data.links)
    .enter().append('line')
    .attr('stroke-width', 2)
    .attr('stroke', '#999');

  const node = container.append('g')
    .selectAll('g')
    .data(data.nodes)
    .enter().append('g')
    .call(drag(simulation));

  const displayStory = function (event, d) {
    node.selectAll('circle').attr('fill', '#69b3a2');
    link.attr('stroke', '#999');

    // Highlight the clicked node
    d3.select(this).select('circle').attr('fill', '#ff5722');

    // Highlight the in-degree edges
    link.attr('stroke', link => link.target.id === d.id ? '#ff5722' : '#999');

    document.getElementById('story-content').innerHTML = d.story;
  };

  function updateVisibility(zoomLevel) {
    const threshold = 2.7;
    console.log("Update Visibility called with zoomLevel:", zoomLevel); 
    node.each(function(d) {
        const nodeElement = d3.select(this);
        const hasMedia = d.images || d.videos || d.audios;

        console.log("Node has media:", hasMedia);
        if (zoomLevel >= threshold && hasMedia) {
            nodeElement.select('circle').style('display', 'none');
            nodeElement.select('text').style('display', 'none');
            nodeElement.select('.media-group').style('display', 'block');
            
            // Add click events to show media in preview window
            nodeElement.selectAll('image').on('click', function() {
                showMediaPreview('image', d3.select(this).attr('xlink:href'));
            });

            nodeElement.selectAll('foreignObject video').on('click', function() {
                showMediaPreview('video', d3.select(this).attr('src'));
            });

            nodeElement.selectAll('.media-group text').on('click', function() {
                showMediaPreview('audio', d.audios[d3.select(this).index()]);
            });

        } else {
            nodeElement.select('circle').style('display', 'block');
            nodeElement.select('text').style('display', 'block');
            nodeElement.select('.media-group').style('display', 'none');
        }
    });
  }

  function showMediaPreview(type, src) {
    // Create a full-screen overlay
    const overlay = d3.select('body').append('div')
        .attr('class', 'media-overlay')
        .style('position', 'fixed')
        .style('top', 0)
        .style('left', 0)
        .style('width', '100%')
        .style('height', '100%')
        .style('background-color', 'rgba(0, 0, 0, 0.8)')
        .style('display', 'flex')
        .style('align-items', 'center')
        .style('justify-content', 'center')
        .style('z-index', 1000);

    let mediaElement;
    if (type === 'image') {
        mediaElement = overlay.append('img')
            .attr('src', src)
            .style('max-width', '90%')
            .style('max-height', '90%');
    } else if (type === 'video') {
        mediaElement = overlay.append('video')
            .attr('src', src)
            .attr('controls', true)
            .style('max-width', '90%')
            .style('max-height', '90%')
            .style('background-color', 'black');
    } else if (type === 'audio') {
        mediaElement = overlay.append('audio')
            .attr('src', src)
            .attr('controls', true);
    }

    // Add a close button
    overlay.append('div')
        .attr('class', 'close-button')
        .text('✕')
        .style('position', 'absolute')
        .style('top', '20px')
        .style('right', '20px')
        .style('font-size', '30px')
        .style('color', 'white')
        .style('cursor', 'pointer')
        .on('click', function() {
            overlay.remove();
        });
  }

  // Append circle elements to each 'g' element
  node.append('circle')
    .attr('r', 27)
    .attr('fill', '#69b3a2');

  // Append text elements to each 'g' element
  node.append('text')
    .attr('x', 0)
    .attr('y', 5) // Adjust this value to center the text vertically
    .attr('text-anchor', 'middle')
    .attr('font-size', '10px')
    .attr('fill', '#fff')
    .text(d => d.keyword);

  node.on('click', displayStory);

  // Append media elements (initially hidden)
  node.each(function(d) {
    const mediaGroup = d3.select(this).append('g')
      .attr('class', 'media-group')
      .style('display', 'none');

    // Add images
    if (d.images) {
      d.images.forEach((img, i) => {
        mediaGroup.append('image')
          .attr('xlink:href', img)
          .attr('x', -20)
          .attr('y', -20 + i * 40)
          .attr('width', 40)
          .attr('height', 40);
      });
    }

    // Add video thumbnails
    if (d.videos) {
      d.videos.forEach((video, i) => {
        const videoGroup = mediaGroup.append('foreignObject')
          .attr('x', -20)
          .attr('y', -20 + (d.images ? d.images.length * 40 : 0) + i * 40)
          .attr('width', 160)
          .attr('height', 90);
   
        videoGroup.append('xhtml:video')
          .attr('src', video)
          .attr('width', '80%')
          .attr('height', '80%')
          .attr('controls', true);
      });
    }

    // Add audio indicators (since audio doesn't have a visual representation)
    if (d.audios) {
      d.audios.forEach((audio, i) => {
        mediaGroup.append('text')
          .attr('x', 0)
          .attr('y', -20 + (d.images ? d.images.length * 40 : 0) + (d.videos ? d.videos.length * 40 : 0) + i * 20)
          .attr('text-anchor', 'middle')
          .attr('fill', 'white')
          .text('🎵');
      });
    }
  });

  node.on('click', displayStory);

  simulation.on('tick', () => {
    link
        .attr('x1', d => d.source.x)
        .attr('y1', d => d.source.y)
        .attr('x2', d => d.target.x)
        .attr('y2', d => d.target.y);

    node
        .attr('transform', d => `translate(${d.x},${d.y})`);
  });

  simulation.on('tick', () => {
    link
        .attr('x1', d => d.source.x)
        .attr('y1', d => d.source.y)
        .attr('x2', d => d.target.x)
        .attr('y2', d => d.target.y);

    node
        .attr('transform', d => `translate(${d.x},${d.y})`);
  });

  function brushed(event) {
    if (event.selection === null) return;
    const [[x0, y0], [x1, y1]] = event.selection;
    node.classed("brushed", d => x0 <= d.x && d.x <= x1 && y0 <= d.y && d.y <= y1);
  }

  function brushended(event) {
    if (event.selection === null) {
      node.classed("brushed", false);
    } else {
      const selectedNodes = node.filter(function() {
        return d3.select(this).classed("brushed");
      });

      if (selectedNodes.size() > 0) {
        const stories = selectedNodes.data().map(d => d.story).join("<br><br>");
        document.getElementById('story-content').innerHTML = stories;
      }
    }
  }

  const brush = d3.brush()
    .extent([[0, 0], [width, height]])
    .on("start brush", brushed)
    .on("end", brushended);

  function enableBrush() {
    svg.on(".zoom", null);

    svg.append("g")
      .attr("class", "brush")
      .call(brush);
  }

  function disableBrush() {
    svg.select(".brush").remove();
    sys.call(zoom);
  }

  document.getElementById('brush-mode').addEventListener('click', () => {
    brushing = !brushing;
    if (brushing) {
      enableBrush();
    } else {
      disableBrush();
    }
  });
}

function drag(simulation) {
  function dragstarted(event, d) {
    if (!event.active) simulation.alphaTarget(0.3).restart();
    d.fx = d.x;
    d.fy = d.y;
  }

  function dragged(event, d) {
    d.fx = event.x;
    d.fy = event.y;
  }

  function dragended(event, d) {
    if (!event.active) simulation.alphaTarget(0);
    d.fx = null;
    d.fy = null;
  }

  return d3.drag()
    .on('start', dragstarted)
    .on('drag', dragged)
    .on('end', dragended);
}

function zoomIn() {
  d3.select('svg').transition().duration(750).call(
    zoom.scaleBy, 1.5
  );
}

function zoomOut() {
  d3.select('svg').transition().duration(750).call(
    zoom.scaleBy, 0.5
  );
}

document.addEventListener("DOMContentLoaded", function() {
  // Fetch the graph data from the JSON file
  fetch('./graph.json')
    .then(response => response.json())
    .then(data => {
      createGraph(data);
    })
    .catch(error => console.error('Error loading the JSON file:', error));

  document.getElementById('zoom-in').addEventListener('click', zoomIn);
  document.getElementById('zoom-out').addEventListener('click', zoomOut);
});
