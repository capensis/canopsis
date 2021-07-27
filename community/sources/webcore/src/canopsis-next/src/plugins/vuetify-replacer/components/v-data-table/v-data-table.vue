<script>
import { VDataTable } from 'vuetify/es5/components/VDataTable';
import { VIcon } from 'vuetify/es5/components/VIcon';
import { consoleWarn } from 'vuetify/es5/util/console';

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
  },
  computed: {
    activeItems() {
      return this.filteredItems.filter(item => !this.isDisabledItem(item));
    },

    everyItem() {
      return this.activeItems.length && this.activeItems.every(this.isSelected);
    },
  },
  methods: {
    /* eslint-disable no-param-reassign */
    genHeaderSortingData(header, children, data, classes) {
      if (!('value' in header)) {
        consoleWarn('Headers must have a value property that corresponds to a value in the v-model array', this);
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

      const pagination = this.computedPagination;

      const addDataAttributes = (sortBy, descending) => {
        const beingSorted = sortBy === header.value;

        if (beingSorted) {
          classes.push('active');
          if (descending) {
            classes.push('desc');
            data.attrs['aria-sort'] = 'descending';
            data.attrs['aria-label'] += ': Sorted descending. Activate to remove sorting.'; // TODO: Localization
          } else {
            classes.push('asc');
            data.attrs['aria-sort'] = 'ascending';
            data.attrs['aria-label'] += ': Sorted ascending. Activate to sort descending.'; // TODO: Localization
          }
        } else {
          data.attrs['aria-label'] += ': Not sorted. Activate to sort ascending.'; // TODO: Localization
        }
      };

      if (this.multiSort) {
        const { multiSortBy = [] } = pagination;
        const sortItemIndex = multiSortBy.findIndex(item => item.sortBy === header.value);
        const sortItem = multiSortBy[sortItemIndex];

        if (sortItem) {
          const sortPriority = this.$createElement('span', {
            class: 'mx-1 caption',
          }, `${sortItemIndex + 1}`);

          children.push(sortPriority, children.pop());

          addDataAttributes(sortItem.sortBy, sortItem.descending);
        }

        return;
      }

      addDataAttributes(pagination.sortBy, pagination.descending);
    },
    /* eslint-enable no-param-reassign */

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
