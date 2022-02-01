// http://nightwatchjs.org/guide#usage

const { ROW_SIZE_KEYS } = require('../../../constants');

module.exports.command = function setCommonFields({
  size,
  row,
  title,
  parameters: {
    filter,
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
    dateInterval,
    infoPopups,
    moreInfos,
    filters,
    openedResolvedFilter,
    statsSelector,
    statSelector,
    statsPointsStyles,
    annotationLine,
    statsColors,
    newColumnNames,
    editColumnNames,
    moveColumnNames,
    deleteColumnNames,
  } = {},
  periodicRefresh,
}) {
  const addInfoPopupModal = this.page.modals.common.addInfoPopup();
  const textEditorModal = this.page.modals.common.textEditor();
  const infoPopupModal = this.page.modals.common.infoPopupSetting();
  const createFilterModal = this.page.modals.common.createFilter();
  const addStatModal = this.page.modals.stats.addStat();
  const statsDateIntervalModal = this.page.modals.stats.statsDateInterval();
  const dateIntervalField = this.page.fields.dateInterval();
  const colorPickerModal = this.page.modals.common.colorPicker();
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
      .setRowSize(ROW_SIZE_KEYS.SMARTPHONE, size.sm)
      .setRowSize(ROW_SIZE_KEYS.TABLET, size.md)
      .setRowSize(ROW_SIZE_KEYS.DESKTOP, size.lg);
  }

  if (title) {
    common
      .clickWidgetTitle()
      .clearWidgetTitle()
      .setWidgetTitle(title);
  }

  if (advanced) {
    common.clickAdvancedSettings();
  }

  if (alarmsList) {
    common.clickAlarmList();
  }

  if (filter) {
    common.clickCreateFilter();

    createFilterModal
      .verifyModalOpened()
      .fillFilterGroups(filter.groups)
      .clickSubmitButton()
      .verifyModalClosed();
  }

  if (dateInterval) {
    common.clickEditDateInterval();

    statsDateIntervalModal.verifyModalOpened()
      .selectPeriodUnit(dateInterval.period);

    if (dateInterval.periodValue !== undefined) {
      statsDateIntervalModal
        .clearPeriodValue()
        .clickPeriodValue()
        .setPeriodValue(dateInterval.periodValue);
    }

    if (dateInterval.range) {
      dateIntervalField.selectRange(dateInterval.range);
    }

    if (dateInterval.calendarStartDate) {
      dateIntervalField
        .clickStartDateButton()
        .clickDatePickerDayTab()
        .selectCalendarDay(dateInterval.calendarStartDate.day)
        .clickDatePickerHoursTab()
        .selectCalendarHour(dateInterval.calendarStartDate.hour)
        .clickDatePickerMinutesTab()
        .selectCalendarMinute(dateInterval.calendarStartDate.minute);
    }

    if (dateInterval.calendarEndDate) {
      dateIntervalField
        .clickEndDateButton()
        .selectCalendarDay(dateInterval.calendarEndDate.day)
        .clickDatePickerHoursTab()
        .selectCalendarHour(dateInterval.calendarEndDate.hour)
        .clickDatePickerMinutesTab()
        .selectCalendarMinute(dateInterval.calendarEndDate.minute);
    }

    if (dateInterval.endDate) {
      dateIntervalField
        .clearEndDate()
        .clickEndDate()
        .setEndDate(dateInterval.endDate);
    }

    if (dateInterval.startDate) {
      dateIntervalField
        .clearStartDate()
        .clickStartDate()
        .setStartDate(dateInterval.startDate);
    }

    statsDateIntervalModal
      .clickSubmitButton()
      .verifyModalClosed();
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

  if (statsSelector) {
    common.clickStatsSelect();

    if (statsSelector.newStats) {
      statsSelector.newStats.forEach((stat) => {
        common.clickAddStat();

        addStatModal
          .verifyModalOpened()
          .selectStatType(stat.type)
          .clickStatTitle()
          .clearStatTitle()
          .setStatTitle(stat.title);

        if (typeof stat.trend === 'boolean') {
          addStatModal.setStatTrend(stat.trend);
        }

        if (typeof stat.recursive === 'boolean') {
          addStatModal.setStatRecursive(stat.recursive);
        }

        if (stat.states) {
          addStatModal
            .clickStatStates()
            .setStatStates(stat.states)
            .clickParameters();
        }

        if (stat.authors) {
          addStatModal
            .clickStatAuthors()
            .clearStatAuthors()
            .setStatAuthors(stat.authors)
            .clickParameters();
        }

        if (stat.sla) {
          addStatModal
            .clickStatSla()
            .clearStatSla()
            .setStatSla(stat.sla);
        }

        addStatModal
          .clickSubmitButton()
          .verifyModalClosed();
      });
    }
  }

  if (statSelector) {
    common.clickStatSelectButton();

    addStatModal
      .verifyModalOpened()
      .selectStatType(statSelector.type)
      .clickStatTitle()
      .clearStatTitle()
      .setStatTitle(statSelector.title);

    if (typeof statSelector.trend === 'boolean') {
      addStatModal.setStatTrend(statSelector.trend);
    }

    if (typeof statSelector.recursive === 'boolean') {
      addStatModal.setStatRecursive(statSelector.recursive);
    }

    if (statSelector.states) {
      addStatModal
        .clickStatStates()
        .setStatStates(statSelector.states)
        .clickParameters();
    }

    if (statSelector.authors) {
      addStatModal
        .clickStatAuthors()
        .clearStatAuthors()
        .setStatAuthors(statSelector.authors)
        .clickParameters();
    }

    if (statSelector.sla) {
      addStatModal
        .clickStatSla()
        .clearStatSla()
        .setStatSla(statSelector.sla);
    }

    addStatModal
      .clickSubmitButton()
      .verifyModalClosed();
  }

  if (statsColors) {
    statsColors.forEach((statColor) => {
      common
        .clickStatsColor()
        .clickStatsColorItem(statColor.title);

      colorPickerModal
        .verifyModalOpened()
        .clickColorField()
        .clearColorField()
        .setColorField(statColor.color)
        .clickSubmitButton()
        .verifyModalClosed();
    });
  }

  if (statsPointsStyles) {
    statsPointsStyles.forEach((statColor) => {
      common
        .clickStatsPointsStyles()
        .selectStatsPointsStylesType(statColor.title, statColor.type);
    });
  }

  if (annotationLine) {
    common
      .clickAnnotationLine()
      .setAnnotationLineEnabled(annotationLine.isEnabled);

    if (annotationLine.isEnabled) {
      common
        .clickAnnotationValue()
        .clearAnnotationValue()
        .setAnnotationValue(annotationLine.value)
        .clickAnnotationLabel()
        .clearAnnotationLabel()
        .setAnnotationLabel(annotationLine.label)
        .clickAnnotationLineColor();

      colorPickerModal
        .verifyModalOpened()
        .clickColorField()
        .clearColorField()
        .setColorField(annotationLine.lineColor)
        .clickSubmitButton()
        .verifyModalClosed();

      common.clickAnnotationLabelColor();

      colorPickerModal
        .verifyModalOpened()
        .clickColorField()
        .clearColorField()
        .setColorField(annotationLine.labelColor)
        .clickSubmitButton()
        .verifyModalClosed();
    }
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

  if (infoPopups) {
    common.clickCreateOrEditInfoPopup();

    infoPopupModal.verifyModalOpened();

    infoPopups.forEach(({ field, template }) => {
      infoPopupModal.clickAddPopup();

      addInfoPopupModal
        .verifyModalOpened()
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

    textEditorModal.verifyModalOpened()
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

  if (filters) {
    common.clickFilters()
      .setMixFilters(filters.isMix);

    if (filters.isMix) {
      common.setFiltersType(filters.type);
    }

    if (filters.groups) {
      common.clickAddFilter();
      createFilterModal.verifyModalOpened()
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

  return this;
};
