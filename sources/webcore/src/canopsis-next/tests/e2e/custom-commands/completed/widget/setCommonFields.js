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
    advanced,
    elementsPerPage,
    infoPopups,
    moreInfos,
    openedResolvedFilter,
    newColumnNames,
    editColumnNames,
    moveColumnNames,
    deleteColumnNames,
    filters,
  } = {},
  periodicRefresh,
}) {
  const addInfoPopup = this.page.modals.common.addInfoPopup();
  const createFilter = this.page.modals.common.createFilter();
  const textEditor = this.page.modals.common.textEditor();
  const infoPopupModal = this.page.modals.common.infoPopupSetting();
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

  if (advanced) {
    common.clickAdvancedSettings();
  }

  if (alarmsList) {
    common.clickAlarmList();
  }

  if (periodicRefresh) {
    common
      .clickPeriodicRefresh()
      .setPeriodicRefreshSwitch(true)
      .clearPeriodicRefreshField()
      .setPeriodicRefreshField(periodicRefresh);
  }

  if (elementsPerPage) {
    common
      .clickElementsPerPage()
      .selectElementsPerPage(elementsPerPage)
      .clickElementsPerPage();
  }

  if (limit) {
    common
      .clickWidgetLimit()
      .clearWidgetLimitField()
      .setWidgetLimitField(limit);
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

  if (openedResolvedFilter) {
    common
      .clickFilterOnOpenResolved()
      .setOpenFilter(openedResolvedFilter.open)
      .setResolvedFilter(openedResolvedFilter.resolve);
  }

  if (filters) {
    common.clickFilters()
      .setMixFilters(filters.isMix);

    if (filters.isMix) {
      common.setFiltersType(filters.type);
    }

    if (filters.groups) {
      common.clickAddFilter();
      createFilter.verifyModalOpened()
        .clearFilterTitle()
        .setFilterTitle(filters.title)
        .fillFilterGroups(filters.groups)
        .clickSubmitButton()
        .verifyModalClosed();
    }

    if (filters.selected) {
      filters.selected.forEach((element) => {
        common.selectFilter(element);
      });
    }
  }

  if (infoPopups) {
    common.clickInfoPopup();

    infoPopupModal.verifyModalOpened();

    infoPopups.forEach(({ field, template }) => {
      infoPopupModal.clickAddPopup();

      addInfoPopup.verifyModalOpened()
        .selectSelectedColumn(field)
        .setTemplate(template)
        .clickSubmitButton()
        .verifyModalClosed();
    });

    infoPopupModal.clickSubmitButton()
      .verifyModalClosed();
  }

  if (moreInfos) {
    common.clickCreateMoreInfos();

    textEditor.verifyModalOpened()
      .clickField()
      .setField(moreInfos)
      .clickSubmitButton()
      .verifyModalClosed();
  }

  if (newColumnNames || editColumnNames || moveColumnNames || deleteColumnNames) {
    common.clickColumnNames();
  }

  if (newColumnNames) {
    newColumnNames.forEach(({ index, data }) => {
      common
        .clickAddColumnName()
        .editColumnName(index, data);
    });
  }

  if (editColumnNames) {
    editColumnNames.forEach(({ index, data }) => {
      common.editColumnName(index, data);
    });
  }

  if (moveColumnNames) {
    moveColumnNames.forEach(({ index, up, down }) => {
      if (up) {
        common.clickColumnNameUpWard(index);
      }

      if (down) {
        common.clickColumnNameDownWard(index);
      }
    });
  }

  if (deleteColumnNames) {
    deleteColumnNames.forEach((index) => {
      common.clickDeleteColumnName(index);
    });
  }

  return this;
};
