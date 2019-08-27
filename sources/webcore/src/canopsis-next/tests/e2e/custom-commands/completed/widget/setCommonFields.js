// http://nightwatchjs.org/guide#usage

module.exports.command = function setCommonFields({
  size,
  row,
  title,
  parameters: {
    sort,
    columnSM,
    columnMD,
    columnLG,
    limit,
    margin,
    heightFactor,
  } = {},
  periodicRefresh,
  advanced = false,
}) {
  const common = this.page.widget.common();

  if (row) {
    common
      .clickRowGridSize()
      .setRow(row);
  } else {
    common.clickRowGridSize();
  }

  if (size) {
    common
      .setSlider('sm', size.sm)
      .setSlider('md', size.md)
      .setSlider('lg', size.lg);
  }

  if (title) {
    common
      .clickWidgetTitle()
      .clearWidgetTitleField()
      .setWidgetTitleField(title);
  }

  if (periodicRefresh) {
    common
      .clickPeriodicRefresh()
      .togglePeriodicRefreshSwitch(true)
      .clearPeriodicRefreshField()
      .setPeriodicRefreshField(periodicRefresh);
  }

  if (limit) {
    common
      .clickWidgetLimit()
      .clearWidgetLimitField()
      .setWidgetLimitField(limit);
  }

  if (advanced) {
    common.clickAdvancedSettings();
  }

  if (sort) {
    common
      .clickDefaultSortColumn()
      .selectSortOrderBy(2)
      .selectSortOrders(2);
  }

  if (columnSM) {
    common.setColumn('SM', columnSM);
  }

  if (columnMD) {
    common.setColumn('MD', columnMD);
  }

  if (columnLG) {
    common.setColumn('LG', columnLG);
  }

  if (margin) {
    common
      .clickMarginBlock()
      .setMargin('top', margin.top)
      .setMargin('right', margin.right)
      .setMargin('bottom', margin.bottom)
      .setMargin('left', margin.left);
  }

  if (heightFactor) {
    common
      .clickHeightFactor()
      .setHeightFactor(heightFactor);
  }

  return this;
};
