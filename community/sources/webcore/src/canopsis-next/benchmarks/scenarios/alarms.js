const { benchmark } = require('../utils/runner');
const { Application } = require('../modules/application');
const { ViewPage } = require('../modules/view');
const { LoginPage } = require('../modules/login');
const { AlarmsListModule } = require('../modules/alarms');

const slowdowns = [1, 2, 4, 6];
const itemsPerPageOptions = [50, 100];

const changeItemsPerPage = async (itemsPerPage, { url, viewId, tabId }) => {
  const application = new Application(url);
  const viewPage = new ViewPage(application);
  const loginPage = new LoginPage(application);
  const alarmsListModule = new AlarmsListModule(application);

  try {
    await application.openBrowser();
    await loginPage.login();
    await viewPage.openById(viewId, { tabId });
    await alarmsListModule.waitTableRow();
    await alarmsListModule.updateItemsPerPage(itemsPerPage);
    await application.page.waitForTimeout(1000);
  } catch (err) {
    console.error(err);

    throw err;
  } finally {
    await application.closeBrowser();
  }
};

benchmark.each(
  itemsPerPageOptions,
  itemsPerPage => `${itemsPerPage} alarms`,
  async (itemsPerPage, measure, { url, viewId, tabId }) => {
    const application = new Application(url);
    const viewPage = new ViewPage(application);
    const loginPage = new LoginPage(application);
    const alarmsListModule = new AlarmsListModule(application);

    await changeItemsPerPage(itemsPerPage, { url, viewId, tabId });

    measure.each(
      slowdowns,
      slowdown => (slowdown === 1 ? 'Render without slowdown' : `Render with ${slowdown}x slowdown`),
      async (slowdown, report) => {
        try {
          await application.openBrowser();
          await loginPage.login();
          await viewPage.openById(viewId, { tabId });

          await application.emulateCPUThrottling(slowdown);

          await application.startMeasurePerformance();

          await application.page.reload();

          await alarmsListModule.waitTableRow();

          const performanceMetrics = await application.stopMeasurePerformance();
          const { JSHeapUsedSize, JSHeapTotalSize } = await application.getPageMetrics();
          const longAnimationFrame = performanceMetrics.findLongestPerformanceTask();
          const { duration, blockingDuration } = longAnimationFrame.args.data;

          report({ duration, blockingDuration, JSHeapUsedSize, JSHeapTotalSize });
        } catch (err) {
          console.error(err);
        } finally {
          await application.closeBrowser();
        }
      },
    );

    measure.each(
      slowdowns,
      slowdown => (slowdown === 1 ? 'Reload without slowdown' : `Reload with ${slowdown}x slowdown`),
      async (slowdown, report) => {
        try {
          await application.openBrowser();
          await loginPage.login();
          await viewPage.openById(viewId, { tabId });

          await application.emulateCPUThrottling(slowdown);

          await alarmsListModule.waitTableRow();

          await application.startMeasurePerformance();

          await Promise.all([
            alarmsListModule.waitProgress(),
            viewPage.clickReload(),
          ]);
          await alarmsListModule.waitProgress(false);

          const performanceMetrics = await application.stopMeasurePerformance();
          const { JSHeapUsedSize, JSHeapTotalSize } = await application.getPageMetrics();
          const longAnimationFrame = performanceMetrics.findLongestPerformanceTask();
          const { duration, blockingDuration } = longAnimationFrame.args.data;

          report({ duration, blockingDuration, JSHeapUsedSize, JSHeapTotalSize });
        } catch (err) {
          console.error(err);
        } finally {
          await application.closeBrowser();
        }
      },
    );

    measure.each(
      slowdowns,
      slowdown => (slowdown === 1 ? 'Open expand panel without slowdown' : `Open expand panel with ${slowdown}x slowdown`),
      async (slowdown, report) => {
        try {
          await application.openBrowser();
          await loginPage.login();
          await viewPage.openById(viewId, { tabId });

          await application.emulateCPUThrottling(slowdown);

          await alarmsListModule.waitTableRow();

          await application.startMeasurePerformance();

          await alarmsListModule.openFirstAlarmRow();
          await alarmsListModule.waitFirstAlarmRowExpandPanel();

          const performanceMetrics = await application.stopMeasurePerformance();
          const { JSHeapUsedSize, JSHeapTotalSize } = await application.getPageMetrics();
          const longAnimationFrame = performanceMetrics.findLongestPerformanceTask();
          const { duration, blockingDuration } = longAnimationFrame.args.data;

          report({ duration, blockingDuration, JSHeapUsedSize, JSHeapTotalSize });
        } catch (err) {
          console.error(err);
        } finally {
          await application.closeBrowser();
        }
      },
    );
  },
);
