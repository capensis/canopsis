class ViewPage {
  constructor(application) {
    this.application = application;
  }

  async openById(id, { tabId } = {}) {
    await this.application.navigate(`/view/${id}${tabId ? `?tabId=${tabId}` : ''}`);
  }

  clickReload() {
    return this.application.page.click('.view-fab-btns > .layout > .v-btn');
  }
}

module.exports = {
  ViewPage,
};
