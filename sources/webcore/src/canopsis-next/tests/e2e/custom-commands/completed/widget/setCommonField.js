// http://nightwatchjs.org/guide#usage

module.exports.command = function setCommonField({
  row,
  sm,
  md,
  lg,
  title,
  periodRefresh,
}) {
  const common = this.page.widget.common();

  if (row) {
    common.clickRowGridSize()
      .setRow('row');
  } else {
    common.clickRowGridSize();
  }
  if (sm) {
    common.setSlider('sm', sm);
  }
  if (md) {
    common.setSlider('md', md);
  }
  if (lg) {
    common.setSlider('lg', lg);
  }
  if (title) {
    common.clickWidgetTitle()
      .clearWidgetTitleField()
      .setWidgetTitleField(title);
  }
  if (periodRefresh) {
    const status = common.getPeriodicRefreshSwitchStatus();

    common.clickPeriodicRefresh();

    if (!status) {
      common.clickPeriodicRefreshSwitch();
    }

    common.clearPeriodicRefreshField()
      .setPeriodicRefreshField(periodRefresh);
  }
  return this;
};
