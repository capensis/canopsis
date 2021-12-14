export const ChartBackgroundPlugin = {
  beforeDraw: (chart) => {
    const { ctx } = chart;
    const { background = {} } = chart.options.plugins || {};

    if (background.color) {
      ctx.save();
      ctx.globalCompositeOperation = 'destination-over';
      ctx.fillStyle = background.color;
      ctx.fillRect(0, 0, chart.width, chart.height);
      ctx.restore();
    }
  },
};
