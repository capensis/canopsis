/**
 * @typedef {Object} ChartEmptyPieFontOptions
 * @property {number} size
 * @property {string} family
 * @property {string} color
 */

/**
 * @typedef {Object} ChartEmptyPieOptions
 * @property {number} width
 * @property {string} color
 * @property {string} text
 * @property {ChartEmptyPieFontOptions} font
 */

export const ChartEmptyPiePlugin = {
  id: 'emptyPie',
  afterDraw(chart) {
    const { datasets } = chart.data;
    const { emptyPie = {} } = chart.options.plugins || {};
    const { width, color, text, font } = emptyPie;

    const hasValue = datasets.some(({ data }) => data.some(Boolean));

    if (!hasValue) {
      const { chartArea, ctx } = chart;
      const { left, top, right, bottom } = chartArea;

      const centerX = (left + right) / 2;
      const centerY = (top + bottom) / 2;

      if (width) {
        const radius = Math.min(right - left, bottom - top) / 2;

        ctx.beginPath();
        ctx.lineWidth = width;
        ctx.strokeStyle = color || 'rgba(0, 0, 0, 0.5)';
        ctx.arc(centerX, centerY, (radius - width || 0), 0, 2 * Math.PI);
        ctx.stroke();
      }

      if (text) {
        if (font) {
          ctx.font = `${font.size}px ${font.family}`;
          ctx.fillStyle = font.color;
        }

        ctx.textAlign = 'center';
        ctx.fillText(text, centerX, centerY);
      }
    }
  },
};
