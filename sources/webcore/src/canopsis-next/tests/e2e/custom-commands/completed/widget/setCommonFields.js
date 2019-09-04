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
    modalType,
    alarmsList,
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
      .setRowSize('sm', size.sm)
      .setRowSize('md', size.md)
      .setRowSize('lg', size.lg);
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
      .setPeriodicRefreshSwitch(true)
      .clearPeriodicRefreshField()
      .setPeriodicRefreshField(periodicRefresh);
  }

  if (alarmsList) {
    common
      .clickAlarmList()
      .clickElementsPerPage()
      .selectElementsPerPage(alarmsList.perPage)
      .clickElementsPerPage();
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
      .selectSortOrderBy(sort.orderBy)
      .selectSortOrders(sort.order);
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

  if (modalType) {
    common
      .clickModalType()
      .clickModalTypeField(modalType);
  }

  return this;
};
