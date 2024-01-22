class ViewPage {
  constructor(application) {
    this.application = application;
  }

  async openById(id, { tabId } = {}) {
    await this.application.navigate(`/view/${id}${tabId ? `?tabId=${tabId}` : ''}`);
  }

  async clickReload() {
    await Promise.all([
      this.application.page.waitForSelector('.v-data-table__progress [role=progressbar]', {
        hidden: true,
      }),
      this.application.page.click('.view-fab-btns > .layout > .flex:nth-of-type(2) > .v-btn'),
    ]);
  }
}

module.exports = {
  ViewPage,
};
