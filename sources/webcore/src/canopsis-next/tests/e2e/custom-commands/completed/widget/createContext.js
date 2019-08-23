// http://nightwatchjs.org/guide#usage

module.exports.command = function createContext({
  common,
  advanced,
}) {
  const createFilter = this.page.modals.common.createFilter();
  const contextWidget = this.page.widget.context();

  if (common) {
    this.completed.widget.setCommonField(common);
  }

  if (common && advanced) {
    const {
      sort,
      columnNames,
      filters,
      typeOfEntities,
    } = advanced;
    contextWidget.clickAdvancedSettings();

    if (sort) {
      const { name, order } = sort;
      contextWidget.clickDefaultSortColumn();

      if (name) {
        contextWidget.selectSortColumn(name);
      }
      if (order) {
        contextWidget.selectSortOrder(order);
      }
    }

    if (columnNames) {
      const { add } = columnNames;
      contextWidget.clickColumnNames();

      if (add) {
        const {
          position,
          label,
          value,
        } = add;
        contextWidget.clickColumnAdd();

        if (position) {
          if (label) {
            contextWidget.setColumnLabel(position, label);
          }
          if (value) {
            contextWidget.setColumnValue(position, value);
          }
        }
      }
    }

    if (filters) {
      const { add } = filters;
      contextWidget.clickFilters();

      if (add) {
        const {
          title, or, and, rule,
        } = add;

        contextWidget.clickAddFilter();

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
    if (typeOfEntities) {
      const {
        component,
        connector,
        resource,
        watcher,
      } = typeOfEntities;
      contextWidget.clickContextTypeOfEntities();

      if (component) {
        contextWidget.clickEntitiesTypeCheckbox(1);
      }

      if (connector) {
        contextWidget.clickEntitiesTypeCheckbox(2);
      }

      if (resource) {
        contextWidget.clickEntitiesTypeCheckbox(3);
      }

      if (watcher) {
        contextWidget.clickEntitiesTypeCheckbox(4);
      }
    }
  }

  if (common) {
    contextWidget.clickSubmitContext();
  }
};
