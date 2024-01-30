class ViewPage {
  constructor(application) {
    this.application = application;
  }

  async openById(id, { tabId } = {}) {
    await this.application.navigate(`/view/${id}${tabId ? `?tabId=${tabId}` : ''}`);
  }

  clickReload() {
    return this.application.page.click('.view-periodic-refresh-btn');
  }
}

module.exports = {
  ViewPage,
};
