import cytoscape from 'cytoscape';
import cytoscapeNodeHtmlLabel from 'cytoscape-node-html-label';
import cytoscapeDagre from 'cytoscape-dagre';
import cytoscapeCise from 'cytoscape-cise';

cytoscapeNodeHtmlLabel(cytoscape);
cytoscape.use(cytoscapeDagre);
cytoscape.use(cytoscapeCise);

export default cytoscape;
