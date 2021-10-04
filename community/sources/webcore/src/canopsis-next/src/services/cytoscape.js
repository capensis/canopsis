import cytoscape from 'cytoscape';
import cytoscapeNodeHtmlLabel from 'cytoscape-node-html-label';
import cytoscapeDagre from 'cytoscape-dagre';

cytoscapeNodeHtmlLabel(cytoscape);
cytoscape.use(cytoscapeDagre);

export default cytoscape;
