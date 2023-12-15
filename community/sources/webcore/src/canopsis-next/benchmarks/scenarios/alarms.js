const { benchmark } = require('../utils/report');
const { Application } = require('../utils/application');
const { ViewPage } = require('../pages/view');
const { LoginPage } = require('../pages/login');
const { AlarmsListModule } = require('../pages/alarms');

const slowdowns = [1, 2, 4, 6];

[50, 100].forEach((itemsPerPage) => {
  benchmark(`Render ${itemsPerPage} alarms`, async (measure, { url, viewId, tabId }) => {
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
    } finally {
      await application.closeBrowser();
    }

    slowdowns.forEach((slowdown) => {
      measure(slowdown === 1 ? 'Without slowdown' : `With ${slowdown}x slowdown`, async (report) => {
        try {
          await application.openBrowser();
          await loginPage.login();

          await viewPage.openById(viewId, { tabId });

          await application.emulateCPUThrottling(slowdown);

          await application.startMeasurePerformance();

          await application.page.reload();

          await alarmsListModule.waitTableRow();

          const performanceMetrics = await application.stopMeasurePerformance();

          const longAnimationFrame = Application.findLongestPerformanceTask(performanceMetrics);

          const { duration } = longAnimationFrame.args.data;

          report(duration);
        } catch (err) {
          console.error(err);
        } finally {
          await application.closeBrowser();
        }
      });
    });
  });
});

benchmark('Reload 50 alarms', (measure, { url, viewId, tabId }) => {
  const application = new Application(url);
  const viewPage = new ViewPage(application);
  const loginPage = new LoginPage(application);
  const alarmsListModule = new AlarmsListModule(application);

  [].forEach((slowdown) => {
    measure(slowdown === 1 ? 'Without slowdown' : `With ${slowdown}x slowdown`, async (report) => {
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

        const longAnimationFrame = Application.findLongestPerformanceTask(performanceMetrics);

        const { duration } = longAnimationFrame.args.data;

        report(duration);
      } catch (err) {
        console.error(err);
      } finally {
        await application.closeBrowser();
      }
    });
  });
});
