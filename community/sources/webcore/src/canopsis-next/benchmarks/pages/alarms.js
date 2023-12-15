class AlarmsListModule {
  constructor(application) {
    this.application = application;
  }

  waitTableRow() {
    return this.application.waitElement('.alarm-list-row');
  }

  waitProgress(visible = true) {
    return this.application.page.waitForSelector('.v-datatable__progress [role=progressbar]', {
      hidden: !visible,
    });
  }

  async updateItemsPerPage(itemsPerPage) {
    const [selectElement] = await this.application.page.$x(
      '//input[@name="itemsPerPage"]/ancestor::*[contains(@class, \'v-select\') and contains(@class, \'v-input\')]',
    );
    selectElement.click();

    const listItemElementSelector = `//*[contains(@class, 'v-menu__content') and contains(@class, 'menuable__content__active')]//*[contains(@class, 'v-list__tile__title') and contains(text(), "${itemsPerPage}")]`;
    await this.application.page.waitForXPath(listItemElementSelector);
    const [listItemElement] = await this.application.page.$x(listItemElementSelector);
    await listItemElement.click();
  }
}

module.exports = {
  AlarmsListModule,
};
