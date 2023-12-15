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
}

module.exports = {
  AlarmsListModule,
};
