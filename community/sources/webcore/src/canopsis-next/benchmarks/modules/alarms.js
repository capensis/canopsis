class AlarmsListModule {
  constructor(application) {
    this.application = application;
  }

  waitTableRow() {
    return this.application.waitElement('.alarm-list-row');
  }

  openFirstAlarmRow() {
    return this.application.page.click('.alarm-list-row .alarms-expand-panel-btn');
  }

  waitFirstAlarmRowExpandPanel() {
    return this.application.waitElement('.v-datatable__expand-content');
  }

  waitProgress(visible = true) {
    return this.application.page.waitForSelector('.v-datatable__progress [role=progressbar]', {
      hidden: !visible,
    });
  }

  async clickItemsPerPageSelector() {
    const selectSelector = '//input[@name="itemsPerPage"]/ancestor::*[contains(@class, \'v-select\') and contains(@class, \'v-input\')]';
    await this.application.page.waitForXPath(selectSelector);
    const [selectElement] = await this.application.page.$x(
      '//input[@name="itemsPerPage"]/ancestor::*[contains(@class, \'v-select\') and contains(@class, \'v-input\')]',
    );

    return selectElement.click();
  }

  async updateItemsPerPage(itemsPerPage) {
    await this.clickItemsPerPageSelector();

    await this.application.page.waitForTimeout(1000);

    await this.application.clickListItemByContent(itemsPerPage);
  }
}

module.exports = {
  AlarmsListModule,
};
