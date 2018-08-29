//cytoscape.js
var cy = window.cy = cytoscape({
    container: document.getElementById('cy'),
  
    boxSelectionEnabled: false,
    autounselectify: true,
  
    style: cytoscape.stylesheet()
      .selector('node')
        .css({
          'content': 'data(id)'
        })
      .selector('edge')
        .css({
          'curve-style': 'bezier',
          'target-arrow-shape': 'triangle',
          'width': 4,
          'line-color': '#ddd',
          'target-arrow-color': '#ddd'
        })
      .selector('.highlighted')
        .css({
          'background-color': '#61bffc',
          'line-color': '#61bffc',
          'target-arrow-color': '#61bffc',
          'transition-property': 'background-color, line-color, target-arrow-color',
          'transition-duration': '0.5s'
        }),
  
    elements: {
        nodes: [
          { data: { id: 'start' } },
          { data: { id: 'e1' } },
          { data: { id: 'e2' } },
          { data: { id: 'e3' } },
          { data: { id: 'end' } }
        ],
  
        edges: [
          { data: { id: 'start"e1', weight: 1, source: 'start', target: 'e1' } },
          { data: { id: 'starte2', weight: 3, source: 'start', target: 'e2' } },
          { data: { id: 'starte3', weight: 4, source: 'start', target: 'e3' } },
          { data: { id: 'e3end', weight: 5, source: 'e3', target: 'end' } },
          { data: { id: 'e2end', weight: 5, source: 'e1', target: 'end' } },
          { data: { id: 'e1end', weight: 5, source: 'e2', target: 'end' } }
         
        ]
      },
  
    layout: {
      name: 'breadthfirst',
      directed: true,
      roots: '#a',
      padding: 10
    }
  });
  
  var bfs = cy.elements().bfs('#a', function(){}, true);
  
  var i = 0;
  var highlightNextEle = function(){
    if( i < bfs.path.length ){
      bfs.path[i].addClass('highlighted');
  
      i++;
      setTimeout(highlightNextEle, 1000);
    }
  };
  
  // kick off first highlight
  highlightNextEle();