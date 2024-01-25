// eslint-disable-next-line import/no-extraneous-dependencies
const puppeteer = require('puppeteer');

const { PerformanceMetrics } = require('../utils/PerformanceMetrics');
const { logInfo } = require('../utils/logger');

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
      args: [
        '--no-sandbox',
        '--window-size=1920,1040',
      ],
    });
  }

  closeBrowser() {
    return this.browser.close();
  }

  async openPage() {
    this.page = await this.browser.newPage();

    await this.page.setViewport({
      width: 1920,
      height: 1040,
    });
  }

  closePage() {
    return this.page.close();
  }

  reloadPage() {
    return this.page.reload();
  }

  async navigate(url) {
    const resultUrl = `${this.url}${url}`;

    await this.page.goto(resultUrl, { timeout: 120_000 });

    logInfo(`Navigate to ${resultUrl}`);
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
    return new PerformanceMetrics(await this.page.tracing.stop());
  }

  emulateCPUThrottling(rate) {
    return this.page.emulateCPUThrottling(rate);
  }

  async clickListItemByContent(content) {
    const listItemElementSelector = `//*[contains(@class, 'v-menu__content') and contains(@class, 'menuable__content__active')]//*[contains(@class, 'v-list-item__title') and contains(text(), "${content}")]`;
    await this.page.waitForXPath(listItemElementSelector);
    const [listItemElement] = await this.page.$x(listItemElementSelector);

    return listItemElement.click();
  }
}

module.exports = {
  Application,
};
