// eslint-disable-next-line import/no-extraneous-dependencies
const puppeteer = require('puppeteer');

class Application {
  static findLongestPerformanceTask(performanceMetrics) {
    const allAnimationTasks = performanceMetrics.traceEvents.filter(({ name, ph }) => name === 'LongAnimationFrame'
      && ph === 'b');

    return allAnimationTasks.reduce((acc, task) => {
      if (acc.args.data.duration < task.args.data.duration) {
        return task;
      }

      return acc;
    });
  }

  browser;

  page;

  constructor(url) {
    this.url = url;
  }

  async openBrowser() {
    this.browser = await puppeteer.launch({
      ignoreHTTPSErrors: true,
      // headless: 'new',
      headless: false,
      args: [
        '--no-sandbox',
        '--window-size=1920,1040',
      ],
    });

    this.page = await this.browser.newPage();

    await this.page.setViewport({
      width: 1920,
      height: 1040,
    });
  }

  closeBrowser() {
    return this.browser.close();
  }

  navigate(url) {
    return this.page.goto(`${this.url}${url}`, { waitUntil: 'load', timeout: 120000 });
  }

  waitElement(selector) {
    return this.page.waitForSelector(selector);
  }

  getPageMetrics() {
    return this.page.metrics();
  }

  startMeasurePerformance(options) {
    return this.page.tracing.start(options);
  }

  async stopMeasurePerformance() {
    /**
     * https://docs.google.com/document/d/1CvAClvFfyA5R-PhYUmn5OOQtYMH4h6I0nSsKchNAySU/preview#heading=h.uxpopqvbjezh
     */
    const performanceMetrics = await this.page.tracing.stop();

    return JSON.parse(performanceMetrics.toString());
  }

  emulateCPUThrottling(rate) {
    return this.page.emulateCPUThrottling(rate);
  }
}

module.exports = {
  Application,
};
