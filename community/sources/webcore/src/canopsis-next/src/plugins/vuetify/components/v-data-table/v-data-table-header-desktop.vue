<script>
import VDataTableHeaderDesktop from 'vuetify/lib/components/VDataTable/VDataTableHeaderDesktop';
import VSimpleCheckbox from 'vuetify/lib/components/VCheckbox/VSimpleCheckbox';
import VLayout from 'vuetify/lib/components/VGrid/VLayout';
import { convertToUnit, wrapInArray } from 'vuetify/lib/util/helpers';

export default {
  extends: VDataTableHeaderDesktop,
  props: {
    disableSelect: {
      type: Boolean,
      default: false,
    },
    ellipsisHeaders: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    genSelectAll() {
      const data = {
        props: {
          value: this.everyItem,
          indeterminate: !this.everyItem && this.someItems,
          disabled: this.disableSelect,
          color: this.checkboxColor ?? '',
        },
        on: {
          input: v => this.$emit('toggle-select-all', v),
        },
      };

      if (this.$scopedSlots['data-table-select']) {
        return this.$scopedSlots['data-table-select'](data);
      }

      return this.$createElement(VSimpleCheckbox, {
        staticClass: 'v-data-table__checkbox',
        ...data,
      });
    },

    genHeader(header) {
      const data = {
        /**
         * Added to avoid dragging error
         */
        key: header.value,
        attrs: {
          role: 'columnheader',
          scope: 'col',

          /**
           * Added to resizing and dragging
           */
          'data-value': header.value,
          'aria-label': header.text || '',
        },
        style: {
          width: convertToUnit(header.width),
          minWidth: convertToUnit(header.width),
        },
        class: [`text-${header.align || 'start'}`, ...wrapInArray(header.class), header.divider && 'v-data-table__divider'],
        on: {},
      };
      const children = [];

      if (header.value === 'data-table-select' && !this.singleSelect) {
        return this.$createElement('th', data, [this.genSelectAll()]);
      }

      children.push(this.$scopedSlots[header.value] ? this.$scopedSlots[header.value]({
        header,
      }) : this.$createElement('span', {
        class: { 'v-data-table-header-span--ellipsis': this.ellipsisHeaders },
        attrs: this.ellipsisHeaders ? { title: header.text } : {},
      }, [header.text]));

      if (!this.disableSort && (header.sortable || !('sortable' in header))) {
        data.on.click = () => this.$emit('sort', header.value);

        const sortIndex = this.options.sortBy.findIndex(k => k === header.value);
        const beingSorted = sortIndex >= 0;
        const isDesc = this.options.sortDesc[sortIndex];
        data.class.push('sortable');
        const {
          ariaLabel,
          ariaSort,
        } = this.getAria(beingSorted, isDesc);
        data.attrs['aria-label'] += `${header.text ? ': ' : ''}${ariaLabel}`;
        data.attrs['aria-sort'] = ariaSort;

        if (beingSorted) {
          data.class.push('active');
          data.class.push(isDesc ? 'desc' : 'asc');
        }

        if (header.align === 'end') {
          children.unshift(this.genSortIcon());
        } else {
          children.push(this.genSortIcon());
        }

        if (this.options.multiSort && beingSorted) {
          children.push(this.$createElement('span', {
            class: 'v-data-table-header__sort-badge',
          }, [String(sortIndex + 1)]));
        }
      }

      if (this.showGroupBy && header.groupable !== false) {
        children.push(this.genGroupByToggle(header));
      }

      if (children.length > 1) {
        return this.$createElement('th', data, [
          this.$createElement(VLayout, { class: 'align-center' }, children),
        ]);
      }

      return this.$createElement('th', data, children);
    },
  },
};
</script>
