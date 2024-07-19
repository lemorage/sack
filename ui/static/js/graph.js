const width = document.getElementById('graph-container').offsetWidth;
const height = document.getElementById('graph-container').offsetHeight;

function createGraph(data) {
  const svg = d3.select('#graph-container').append('svg')
    .attr('width', width)
    .attr('height', height);

  const simulation = d3.forceSimulation(data.nodes)
    .force("link", d3.forceLink(data.links).id(d => d.id).distance(100))
    .force('charge', d3.forceManyBody().strength(-300))
    .force('center', d3.forceCenter(width / 2, height / 2));

  const link = svg.append('g')
    .selectAll('line')
    .data(data.links)
    .enter().append('line')
    .attr('stroke-width', 2)
    .attr('stroke', '#999');

  // Create a 'g' element for each node to group the circle and text elements
  const node = svg.append('g')
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

  simulation.on('tick', () => {
    link
        .attr('x1', d => d.source.x)
        .attr('y1', d => d.source.y)
        .attr('x2', d => d.target.x)
        .attr('y2', d => d.target.y);

    node
        .attr('transform', d => `translate(${d.x},${d.y})`);
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

document.addEventListener("DOMContentLoaded", function() {
  // Fetch the graph data from the JSON file
  fetch('./graph.json')
      .then(response => response.json())
      .then(data => {
          createGraph(data);
      })
      .catch(error => console.error('Error loading the JSON file:', error));
});
