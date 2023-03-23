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
/*  beforeDatasetDraw(chart, dataset) {
    dataset.meta.data.forEach((bar) => {
      bar.x = bar.x - 96.00582885742182;
    });
  },
  afterDatasetDraw(chart, dataset) {
    console.log(dataset, chart.canvas.getContext('2d'));

    dataset.meta.data.forEach((bar) => {
      bar.x = bar.x + 96.00582885742182;
    });
  }, */
};
