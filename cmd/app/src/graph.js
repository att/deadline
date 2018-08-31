var cy = window.cy = cytoscape({
  container: document.getElementById('cy'),

  boxSelectionEnabled: false,
  autounselectify: true,

  layout: {
    name: 'dagre'
  },

  style: cytoscape.stylesheet() 
      .selector('node')
        .css({
          'content': 'data(id)',
          'text-opacity': 0.5,
          'text-valign': 'center',
          'text-halign': 'right',
          'background-color': '#468499'
        })
      .selector('edge')
        .css({
          'curve-style': 'bezier',
          'width': 4,
          'target-arrow-shape': 'triangle',
          'line-color': '#9dbaea',
          'target-arrow-color': '#9dbaea'
        })
      .selector('.bad')
        .css({
        'background-color': '#f6546a',
        'line-color': '#f6546a',
        'target-arrow-color': '#f6546a',
        'transition-property': 'background-color, line-color, target-arrow-color',
        'transition-duration': '0.5s'
        })
      .selector('.good')
	.css({
	'background-color': '#006700',
	'line-color': '#006700',
	'target-arrow-color': '#006700',
	'transition-property': 'background-color, line-color, target-arrow-color',
	'transition-duration': '0.5s'
	})
	.selector('.wait')
        .css({
        'background-color': '#ffe062',
        'line-color': '#ffe062',
        'target-arrow-color': '#ffe062',
        'transition-property': 'background-color, line-color, target-arrow-color',
        'transition-duration': '0.5s'
        }),
      
      
  elements: {
    nodes: [

      { data: { id: 'schedule1' } },
      { data: { id: 'first event' } },
      { data: { id: 'second event' } },
      { data: { id: 'third event' } },
      { data: { id: 'fourth event' } },
      { data: { id: 'action' } },
      { data: { id: 'error' } },
      { data: { id: 'end' } },
    ],
    edges: [
      { data: { id: 's1', source: 'schedule1', target: 'first event' } },
      { data: { id: 's2', source: 'schedule1', target: 'second event' } },
      { data: { id: 's3', source: 'schedule1', target: 'third event' } },
      { data: { id: 's4', source: 'schedule1', target: 'fourth event' } },
      { data: { id: '1a', source: 'first event' , target: 'action' } },
      { data: { id: '2a', source: 'second event' , target: 'action' } },
      { data: { id: 'ae', source: 'action' , target: 'end' } },
      { data: { id: '3e', source: 'third event' , target: 'error' } },
    ]
  },

});


setTimeout(function() {
  cy.$id('first event').addClass('good')
}

, 3*1000);
setTimeout(function() {
  cy.$id('second event').addClass('good');
}
, 6 *1000);
setTimeout(function() {
cy.$id('action').addClass('good');
}
, 9 *1000);


setTimeout(function() {
cy.$id('third event').addClass('bad');
}
, 12*1000);

setTimeout(function() {
  cy.$id('error').addClass('bad');
  }
  , 15*1000);

setTimeout(function() {
cy.$id('fourth event').addClass('wait'); }
, 18*1000);


