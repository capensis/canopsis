<script>
import { VDataTable } from 'vuetify/es5/components/VDataTable';
import { VIcon } from 'vuetify/es5/components/VIcon';
import { VCheckbox } from 'vuetify/es5/components/VCheckbox';
import { VLayout } from 'vuetify/es5/components/VGrid';
import { consoleWarn } from 'vuetify/es5/util/console';
import { getObjectValueByPath } from 'vuetify/es5/util/helpers';

import { DEFAULT_MAX_MULTI_SORT_COLUMNS_COUNT } from '@/config';

import { isDarkColor } from '@/helpers/color';

import ExpandTransitionGenerator from '../transitions/expand-transition';

export default {
  extends: VDataTable,
  props: {
    isDisabledItem: {
      type: Function,
      default: item => !item,
    },
    multiSort: {
      type: Boolean,
      default: false,
    },
    tableClass: {
      type: String,
      required: false,
    },
    dense: {
      type: Boolean,
      default: false,
    },
    ultraDense: {
      type: Boolean,
      default: false,
    },
    ellipsisHeaders: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isDark() {
      return isDarkColor(this.$vuetify.theme['table-background']);
    },

    activeItems() {
      return this.filteredItems.filter(item => !this.isDisabledItem(item));
    },

    everyItem() {
      return this.activeItems.length && this.activeItems.every(this.isSelected);
    },

    classes() {
      return {
        'v-datatable v-table': true,
        'v-datatable--expand': this.expand,
        'v-datatable--select-all': this.selectAll !== false,
        'v-datatable--dense': this.dense,
        'v-datatable--ultra-dense': this.ultraDense,
        [this.tableClass]: !!this.tableClass,
        ...this.themeClasses,
      };
    },
  },
  methods: {
    genExpandedRow(props) {
      const children = [];
      if (this.isExpanded(props.item)) {
        const expand = this.$createElement('div', {
          class: 'v-datatable__expand-content',
          key: getObjectValueByPath(props.item, this.itemKey),
        }, [this.$scopedSlots.expand(props)]);
        children.push(expand);
      }

      const transition = this.$createElement('transition-group', {
        class: 'v-datatable__expand-col',
        attrs: { colspan: this.headerColumns },
        props: {
          tag: 'td',
        },
        on: ExpandTransitionGenerator('v-datatable__expand-col--expanded'),
      }, children);

      return this.genTR([transition], {
        class: 'v-datatable__expand-row',
        key: `${getObjectValueByPath(props.item, this.itemKey) || props.index}-row`,
      });
    },

    /* eslint-disable no-param-reassign */
    /**
     * Get header sorting data.
     *
     * @note Was replaced for multi sort support
     *
     * @param {Object} header
     * @param {Object[]} children
     * @param {Object} data
     * @param {string[]} classes
     */
    genHeaderSortingData(header, children, data, classes) {
      if (!('value' in header)) {
        consoleWarn('Headers must have a value property that corresponds to a value in the v-model array', this);
      }

      /**
       * Add data attributes into header item
       *
       * @param {string | number} sortBy
       * @param {boolean} descending
       */
      const addDataAttributes = (sortBy, descending) => {
        const beingSorted = sortBy === header.value;

        if (beingSorted) {
          classes.push('active');
          if (descending) {
            classes.push('desc');
            data.attrs['aria-sort'] = 'descending';
            data.attrs['aria-label'] += ': Sorted descending. Activate to remove sorting.';
          } else {
            classes.push('asc');
            data.attrs['aria-sort'] = 'ascending';
            data.attrs['aria-label'] += ': Sorted ascending. Activate to sort descending.';
          }
        } else {
          data.attrs['aria-label'] += ': Not sorted. Activate to sort ascending.';
        }
      };

      const pagination = this.computedPagination;

      /**
       * Added multi sort support
       */
      if (this.multiSort) {
        const { multiSortBy = [] } = pagination;
        const sortItemIndex = multiSortBy.findIndex(item => item.sortBy === header.value);
        const sortItem = multiSortBy[sortItemIndex];

        if (sortItem) {
          const sortPriority = this.$createElement('span', {
            class: 'v-datatable-header__sort-badge',
          }, `${sortItemIndex + 1}`);

          children.push(sortPriority);

          addDataAttributes(sortItem.sortBy, sortItem.descending);
        } else if (multiSortBy.length >= DEFAULT_MAX_MULTI_SORT_COLUMNS_COUNT) {
          /**
           * This condition was added for compliance with max multiSort limit
           */
          return;
        }
      } else {
        addDataAttributes(pagination.sortBy, pagination.descending);
      }

      data.attrs.tabIndex = 0;
      data.on = {
        click: () => {
          this.expanded = {};
          this.sort(header.value);
        },
        keydown: (e) => {
          // check for space
          if (e.keyCode === 32) {
            e.preventDefault();
            this.sort(header.value);
          }
        },
      };

      classes.push('sortable');

      const icon = this.$createElement('span', [this.$createElement(VIcon, {
        props: {
          small: true,
        },
      }, this.sortIcon)]);

      if (!header.align || header.align === 'left') {
        children.push(icon);
      } else {
        children.unshift(icon);
      }
    },
    /* eslint-enable no-param-reassign */

    genHeaderData(header, children, key) {
      const classes = ['column'];
      const data = {
        key,
        attrs: {
          role: 'columnheader',
          scope: 'col',
          width: header.width || null,
          'aria-label': header[this.headerText] || '',
          'data-value': header.value || '',
          'aria-sort': 'none',
        },
      };

      if (header.sortable == null || header.sortable) {
        this.genHeaderSortingData(header, children, data, classes);
      } else {
        data.attrs['aria-label'] += ': Not sorted.'; // TODO: Localization
      }

      classes.push(`text-xs-${header.align || 'left'}`);
      if (Array.isArray(header.class)) {
        classes.push(...header.class);
      } else if (header.class) {
        classes.push(header.class);
      }
      data.class = classes;

      return [
        data,
        children.length > 1
          ? [this.$createElement(VLayout, { class: 'align-center' }, children)]
          : children,
      ];
    },

    genHeader: function genHeader(header, key) {
      const array = [
        this.$scopedSlots.headerCell
          ? this.$scopedSlots.headerCell({ header })
          : this.$createElement(
            'span',
            {
              attrs: this.ellipsisHeaders ? { title: header[this.headerText] } : {},
              class: { 'v-datatable-header-span--ellipsis': this.ellipsisHeaders },
            },
            header[this.headerText],
          ),
      ];

      return this.$createElement.apply(this, ['th', ...this.genHeaderData(header, array, key)]);
    },

    /**
     * Get thead element for a table
     *
     * @note Was replaced for disabling select all checkbox if all items is disabled
     */
    genTHead() {
      if (this.hideHeaders) {
        return null;
      }

      let children = [];

      if (this.$scopedSlots.headers) {
        const row = this.$scopedSlots.headers({
          headers: this.headers,
          indeterminate: this.indeterminate,
          all: this.everyItem,
        });

        children = [this.hasTag(row, 'th') ? this.genTR(row) : row, this.genTProgress()];
      } else {
        const row = this.headers.map((o, i) => this.genHeader(o, this.headerKey ? o[this.headerKey] : i));
        const checkbox = this.$createElement(VCheckbox, {
          props: {
            dark: this.dark,
            light: this.light,
            color: this.selectAll === true ? '' : this.selectAll,
            hideDetails: true,
            inputValue: this.everyItem,
            indeterminate: this.indeterminate,

            /**
             * disabled added for case with all inactive items
             */
            disabled: !this.activeItems.length,
          },
          on: { change: this.toggle },
        });

        if (this.hasSelectAll) {
          row.unshift(this.$createElement('th', [checkbox]));
        }

        children = [this.genTR(row), this.genTProgress()];
      }

      return this.$createElement('thead', [children]);
    },

    /**
     * Update pagination parameters for the sorting
     *
     * @note Was replaced for multi sort support
     *
     * @param {string | number} index
     */
    sort(index) {
      if (this.multiSort) {
        this.updateMultiSort(index);

        return;
      }

      const { sortBy, descending } = this.computedPagination;

      if (sortBy === null) {
        this.updatePagination({ sortBy: index, descending: false });
      } else if (sortBy === index && !descending) {
        this.updatePagination({ descending: true });
      } else if (sortBy !== index) {
        this.updatePagination({ sortBy: index, descending: false });
      } else if (!this.mustSort) {
        this.updatePagination({ sortBy: null, descending: null });
      } else {
        this.updatePagination({ sortBy: index, descending: false });
      }
    },

    /**
     * Update pagination parameters for the multi sorting
     *
     * @note New method
     *
     * @param {string | number} index
     */
    updateMultiSort(index) {
      const { multiSortBy = [] } = this.computedPagination;
      let newMultiSortBy = [...multiSortBy];

      const sortItemIndex = multiSortBy.findIndex(item => item.sortBy === index);
      const sortItem = multiSortBy[sortItemIndex];

      if (sortItem) {
        if (!sortItem.descending) {
          newMultiSortBy[sortItemIndex] = { ...sortItem, descending: true };
        } else {
          newMultiSortBy = newMultiSortBy.filter(item => item.sortBy !== index);
        }
      } else {
        newMultiSortBy.push({ sortBy: index, descending: false });
      }

      this.updatePagination({ multiSortBy: newMultiSortBy });
    },
  },
};
</script>
