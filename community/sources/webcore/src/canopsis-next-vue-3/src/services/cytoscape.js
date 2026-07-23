import cytoscape from 'cytoscape';
import cytoscapeCise from 'cytoscape-cise';
import cytoscapeDagre from 'cytoscape-dagre';
import cytoscapeNodeHtmlLabel from 'cytoscape-node-html-label';

cytoscapeNodeHtmlLabel(cytoscape);
cytoscape.use(cytoscapeDagre);
cytoscape.use(cytoscapeCise);

export default cytoscape;
