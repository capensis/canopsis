class ViewPage {
  constructor(application) {
    this.application = application;
  }

  async openById(id, { tabId } = {}) {
    await this.application.navigate(`/view/${id}${tabId ? `?tabId=${tabId}` : ''}`);
  }

  async clickReload() {
    await Promise.all([
      this.application.page.waitForSelector('.v-datatable__progress [role=progressbar]', {
        hidden: true,
      }),
      this.application.page.click('.view-fab-btns > .layout > .v-btn'),
    ]);
  }
}

module.exports = {
  ViewPage,
};
