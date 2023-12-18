// eslint-disable-next-line import/no-extraneous-dependencies
const puppeteer = require('puppeteer');

const { PerformanceMetrics } = require('../utils/PerformanceMetrics');

class Application {
  browser;

  page;

  constructor(url) {
    this.url = url;
  }

  async openBrowser() {
    this.browser = await puppeteer.launch({
      ignoreHTTPSErrors: true,
      headless: 'new',
      // headless: false,
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
    return new PerformanceMetrics(await this.page.tracing.stop());
  }

  emulateCPUThrottling(rate) {
    return this.page.emulateCPUThrottling(rate);
  }

  async clickListItemByContent(content) {
    const listItemElementSelector = `//*[contains(@class, 'v-menu__content') and contains(@class, 'menuable__content__active')]//*[contains(@class, 'v-list__tile__title') and contains(text(), "${content}")]`;
    await this.page.waitForXPath(listItemElementSelector);
    const [listItemElement] = await this.page.$x(listItemElementSelector);

    return listItemElement.click();
  }
}

module.exports = {
  Application,
};
