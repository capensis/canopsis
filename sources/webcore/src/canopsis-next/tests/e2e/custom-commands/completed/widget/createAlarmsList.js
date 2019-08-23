// http://nightwatchjs.org/guide#usage

module.exports.command = function createAlarmsList({
  common,
  advanced,
}) {
  const textEditor = this.page.modals.common.textEditor();
  const addInfoPopup = this.page.modals.alarm.addInfoPopup();
  const infoPopupSetting = this.page.modals.alarm.infoPopupSetting();
  const createFilter = this.page.modals.common.createFilter();
  const alarmsWidget = this.page.widget.alarms();

  if (common) {
    this.completed.widget.setCommonField(common);
  }

  if (common && advanced) {
    const {
      sort,
      columnNames,
      defaultNumberOfElementsPerPage,
      filterOnOpenResolved,
      filters,
      infoPopap,
      moreInfo,
      enableHtml,
      ackGroup,
    } = advanced;
    alarmsWidget.clickAdvancedSettings();

    if (sort) {
      const { name, order } = sort;
      alarmsWidget.clickDefaultSortColumn();

      if (name) {
        alarmsWidget.selectSortColumn(name);
      }
      if (order) {
        alarmsWidget.selectSortOrder(order);
      }
    }

    if (columnNames) {
      const { add } = columnNames;
      alarmsWidget.clickColumnNames();

      if (add) {
        const {
          position,
          label,
          value,
        } = add;
        alarmsWidget.clickColumnAdd();

        if (position) {
          if (label) {
            alarmsWidget.setColumnLabel(position, label);
          }
          if (value) {
            alarmsWidget.setColumnValue(position, value);
          }
        }
      }
    }

    if (defaultNumberOfElementsPerPage) {
      const { count } = defaultNumberOfElementsPerPage;
      alarmsWidget.clickDefaultNumberOfElementsPerPage();

      if (count) {
        alarmsWidget.selectElementsPerPage(count);
      }
    }

    if (filterOnOpenResolved) {
      const { open, resolved } = filterOnOpenResolved;
      alarmsWidget.clickFilterOnOpenResolved();

      if (open) {
        alarmsWidget.clickOpenFilter();
      }
      if (resolved) {
        alarmsWidget.clickResolvedFilter();
      }
    }

    if (filters) {
      const { add } = filters;
      alarmsWidget.clickFilters();

      if (add) {
        const {
          title, or, and, rule,
        } = add;

        alarmsWidget.clickAddFilter();

        createFilter.verifyModalOpened();

        if (title) {
          createFilter.setFilterTitle(title);
        }
        if (or) {
          createFilter.clickRadioOr();
        }
        if (and) {
          createFilter.clickRadioAnd();
        }
        if (rule) {
          const { field, operator } = rule;
          createFilter.clickAddRule();

          if (field) {
            createFilter.selectFieldRule(field);
          }
          if (operator) {
            createFilter.selectOperatorRule(operator);
          }
        }
        createFilter.clickSubmitFilter()
          .verifyModalClosed();
      }
    }

    if (infoPopap) {
      const { add } = infoPopap;
      alarmsWidget.clickInfoPopupButton();

      infoPopupSetting.verifyModalOpened();
      if (add) {
        const { column, template } = add;
        infoPopupSetting.clickAddPopup();

        addInfoPopup.verifyModalOpened();

        if (column) {
          addInfoPopup.selectSelectedColumn(column);
        }
        if (template) {
          addInfoPopup.setTemplate(template);
        }

        addInfoPopup.clickSubmitButton()
          .verifyModalClosed();
      }
      infoPopupSetting.clickSubmitButton()
        .verifyModalClosed();
    }

    if (moreInfo) {
      const { text } = moreInfo;
      alarmsWidget.clickCreateEditMore();

      textEditor.verifyModalOpened();

      if (text) {
        textEditor.setRTE(text);
      }

      textEditor.clickSubmitButton()
        .verifyModalClosed();
    }

    if (enableHtml) {
      alarmsWidget.clickEnableHtml();
    }

    if (ackGroup) {
      const {
        isAckNoteRequired,
        isMultiAckEnabled,
        fastAckOutput,
      } = ackGroup;

      alarmsWidget.clickAckGroup();

      if (isAckNoteRequired) {
        alarmsWidget.clickIsAckNoteRequired();
      }

      if (isMultiAckEnabled) {
        alarmsWidget.clickIsMultiAckEnabled();
      }

      if (fastAckOutput) {
        const {
          enabled,
          output,
        } = fastAckOutput;
        alarmsWidget.clickFastAckOutput();

        if (enabled) {
          alarmsWidget.clickFastAckOutputSwitch();
        }
        if (enabled && output) {
          alarmsWidget.setFastAckOutputText(output);
        }
      }
    }
  }

  if (common) {
    alarmsWidget.clickSubmitAlarms();
  }
};
