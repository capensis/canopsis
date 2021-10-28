import { generateChart } from '../helpers/generate-chart';

export const Bar = generateChart('bar-chart', 'bar');
export const Doughnut = generateChart('doughnut-chart', 'doughnut');
export const Line = generateChart('line-chart', 'line');
export const Pie = generateChart('pie-chart', 'pie');
export const PolarArea = generateChart('polar-chart', 'polarArea');
export const Radar = generateChart('radar-chart', 'radar');
export const Bubble = generateChart('bubble-chart', 'bubble');
export const Scatter = generateChart('scatter-chart', 'scatter');

export default {
  Bar,
  Doughnut,
  Line,
  Pie,
  PolarArea,
  Radar,
  Bubble,
  Scatter,
};
