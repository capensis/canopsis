<script>
import { VDataTable } from 'vuetify/es5/components/VDataTable';
import { VIcon } from 'vuetify/es5/components/VIcon';
import { VCheckbox } from 'vuetify/es5/components/VCheckbox';
import { consoleWarn } from 'vuetify/es5/util/console';

import { DEFAULT_MAX_MULTI_SORT_COLUMNS_COUNT } from '@/config';

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
  },
  computed: {
    activeItems() {
      return this.filteredItems.filter(item => !this.isDisabledItem(item));
    },

    everyItem() {
      return this.activeItems.length && this.activeItems.every(this.isSelected);
    },

    classes() {
      return {
        'v-datatable v-table': true,
        'v-datatable--select-all': this.selectAll !== false,
        'v-datatable--dense': this.dense,
        [this.tableClass]: !!this.tableClass,
        ...this.themeClasses,
      };
    },
  },
  methods: {
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
            data.attrs['aria-label'] += ': Sorted descending. Activate to remove sorting.'; // vuetify TODO: Localization
          } else {
            classes.push('asc');
            data.attrs['aria-sort'] = 'ascending';
            data.attrs['aria-label'] += ': Sorted ascending. Activate to sort descending.'; // vuetify TODO: Localization
          }
        } else {
          data.attrs['aria-label'] += ': Not sorted. Activate to sort ascending.'; // vuetify TODO: Localization
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

      const icon = this.$createElement(VIcon, {
        props: {
          small: true,
        },
      }, this.sortIcon);
      if (!header.align || header.align === 'left') {
        children.push(icon);
      } else {
        children.unshift(icon);
      }
    },
    /* eslint-enable no-param-reassign */

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

<style lang="scss">
$densePadding: 6px;
$denseCellHeight: 32px;
$denseColorIndicatorPadding: 1px 5px;

table.v-datatable {
  .v-datatable-header__sort-badge {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    border: 0;
    border-radius: 50%;
    min-width: 18px;
    min-height: 18px;
    height: 18px;
    width: 18px;
    background-color: rgba(0, 0, 0, .12);
    color: rgba(0, 0, 0, .87);
  }

  &--dense.v-datatable .service-dependencies {
    .v-treeview-node__root {
      min-height: $denseCellHeight;
      height: $denseCellHeight;

      .v-btn {
        width: $denseCellHeight - 4;
        height: $denseCellHeight - 4;
        margin: 2px;
      }
    }

    tbody, thead {
      td, th {
        padding: 0 $densePadding;
      }

      td:not(.v-datatable__expand-col) {
        height: $denseCellHeight;

        .v-btn {
          margin-top: 0;
          margin-bottom: 0;
        }

        .c-action-btn__button {
          margin: 0 !important;
        }

        .color-indicator {
          padding: $denseColorIndicatorPadding;
        }
      }
    }
  }
}
</style>
