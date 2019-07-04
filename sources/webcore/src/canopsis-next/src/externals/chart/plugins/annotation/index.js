import Chart from 'chart.js';
import elementAnnotationGenerator from 'chartjs-plugin-annotation/src/element';
import lineAnnotationGenerator from 'chartjs-plugin-annotation/src/types/line';
import boxAnnotationGenerator from 'chartjs-plugin-annotation/src/types/box';

import annotationGenerator from './annotation';

export const AnnotationChart = { ...Chart };

AnnotationChart.Annotation = {};

AnnotationChart.Annotation.drawTimeOptions = {
  afterDraw: 'afterDraw',
  afterDatasetsDraw: 'afterDatasetsDraw',
  beforeDatasetsDraw: 'beforeDatasetsDraw',
};

AnnotationChart.Annotation.defaults = {
  drawTime: 'afterDatasetsDraw',
  dblClickSpeed: 350, // ms
  events: [],
  annotations: [],
};

AnnotationChart.Annotation.labelDefaults = {
  backgroundColor: 'rgba(0,0,0,0.8)',
  fontFamily: Chart.defaults.global.defaultFontFamily,
  fontSize: Chart.defaults.global.defaultFontSize,
  fontStyle: 'bold',
  fontColor: '#fff',
  xPadding: 6,
  yPadding: 6,
  cornerRadius: 6,
  position: 'center',
  xAdjust: 0,
  yAdjust: 0,
  enabled: false,
  content: null,
};

AnnotationChart.Annotation.Element = elementAnnotationGenerator(AnnotationChart);

AnnotationChart.Annotation.types = {
  line: lineAnnotationGenerator(AnnotationChart),
  box: boxAnnotationGenerator(AnnotationChart),
};

export default annotationGenerator(AnnotationChart);
