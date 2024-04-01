<script>
import { VDataTable } from 'vuetify/lib/components/VDataTable';
import { getObjectValueByPath, getSlot, getPrefixedScopedSlots } from 'vuetify/lib/util/helpers';
import ExpandTransitionGenerator from 'vuetify/lib/components/transitions/expand-transition';

import VSimpleTable from './v-simple-table.vue';
import VDataTableHeader from './v-data-table-header.vue';

export default {
  components: { VSimpleTable, VDataTableHeader },
  extends: VDataTable,
  props: {
    ultraDense: {
      type: Boolean,
      default: false,
    },
    itemSelectable: {
      type: Function,
      required: false,
    },
    ellipsisHeaders: {
      type: Boolean,
      default: false,
    },
  },
  watch: {
    /**
     * We've added watcher for resetting page if page is out of range
     */
    serverItemsLength(serverItemsLength) {
      const page = this.options?.page ?? 1;
      const pageCount = Math.ceil(serverItemsLength / this.computedItemsPerPage) || 1;

      if (page > pageCount) {
        this.$emit('update:options', {
          ...this.options,

          page: pageCount,
        });
      }
    },
  },
  methods: {
    /**
     * We've added expand transition here
     */
    genScopedRows(items) {
      const rows = [];

      for (let i = 0; i < items.length; i += 1) {
        const item = items[i];
        let children = [];

        rows.push(this.$scopedSlots.item({
          ...this.createItemProps(item, i),
          isMobile: this.isMobile,
        }));

        if (this.isExpanded(item)) {
          children = this.$createElement('div', {
            class: 'v-data-table__expanded__content',
            key: `expand-${getObjectValueByPath(item, this.itemKey)}`,
          }, this.$scopedSlots['expanded-item']({
            headers: this.computedHeaders,
            isMobile: this.isMobile,
            index: i,
            item,
          }));
        }

        const transition = this.$createElement('transition-group', {
          class: 'v-data-table__expanded__col',
          attrs: { colspan: this.computedHeaders.length },
          props: {
            tag: 'td',
          },
          on: ExpandTransitionGenerator('v-data-table__expanded__col'),
        }, [children]);

        rows.push(this.$createElement('tr', { class: 'v-data-table__expanded v-data-table__expanded__row' }, [transition]));
      }

      return rows;
    },

    genDefaultExpandedRow(item, index) {
      const isExpanded = this.isExpanded(item);
      const classes = {
        'v-data-table__expanded v-data-table__expanded__row': isExpanded,
      };
      let children = [];

      const headerRow = this.genDefaultSimpleRow(item, index, classes);

      if (isExpanded) {
        children = this.$createElement('div', {
          class: 'v-data-table__expanded__content',
          key: `expand-${getObjectValueByPath(item, this.itemKey)}`,
        }, this.$scopedSlots['expanded-item']({
          headers: this.computedHeaders,
          isMobile: this.isMobile,
          index,
          item,
        }));
      }

      const transition = this.$createElement('transition-group', {
        class: 'v-data-table__expanded__col',
        attrs: { colspan: this.computedHeaders.length },
        props: {
          tag: 'td',
        },
        on: ExpandTransitionGenerator('v-data-table__expanded__col'),
      }, [children]);

      const expandedRow = this.$createElement('tr', {
        staticClass: 'v-data-table__expanded v-data-table__expanded__content',
      }, [transition]);

      return [headerRow, expandedRow];
    },

    genDefaultScopedSlot(props) {
      const simpleProps = {
        height: this.height,
        fixedHeader: this.fixedHeader,
        dense: this.dense,
        ultraDense: this.ultraDense,
      };
      // if (this.virtualRows) {
      //   return this.$createElement(VVirtualTable, {
      //     props: Object.assign(simpleProps, {
      //       items: props.items,
      //       height: this.height,
      //       rowHeight: this.dense ? 24 : 48,
      //       headerHeight: this.dense ? 32 : 48,
      //       // TODO: expose rest of props from virtual table?
      //     }),
      //     scopedSlots: {
      //       items: ({ items }) => this.genItems(items, props) as any,
      //     },
      //   }, [
      //     this.proxySlot('body.before', [this.genCaption(props), this.genHeaders(props)]),
      //     this.proxySlot('bottom', this.genFooters(props)),
      //   ])
      // }

      return this.$createElement(VSimpleTable, {
        props: simpleProps,
        class: {
          'v-data-table--expand': this.showExpand,
          'v-data-table--mobile': this.isMobile,
          'v-data-table--selectable': this.showSelect,
        },
      }, [
        this.proxySlot('top', getSlot(this, 'top', {
          ...props,
          isMobile: this.isMobile,
        }, true)),
        this.genCaption(props),
        this.genColgroup(props),
        this.genHeaders(props),
        this.genBody(props),
        this.genFoot(props),
        this.proxySlot('bottom', this.genFooters(props)),
      ]);
    },
    genHeaders(props) {
      const data = {
        props: {
          ...this.sanitizedHeaderProps,

          headers: this.computedHeaders,
          options: props.options,
          mobile: this.isMobile,
          showGroupBy: this.showGroupBy,
          checkboxColor: this.checkboxColor,
          someItems: this.someItems,
          everyItem: this.everyItem,
          singleSelect: this.singleSelect,
          ellipsisHeaders: this.ellipsisHeaders,
          disableSort: this.disableSort,
          disableSelect: this.selectableItems.length === 0,
        },
        on: {
          sort: props.sort,
          group: props.group,
          'toggle-select-all': this.toggleSelectAll,
        },
      }; // TODO: rename to 'head'? (thead, tbody, tfoot)

      const children = [getSlot(this, 'header', { ...data, isMobile: this.isMobile })];

      if (!this.hideDefaultHeader) {
        const scopedSlots = getPrefixedScopedSlots('header.', this.$scopedSlots);
        children.push(this.$createElement(VDataTableHeader, { ...data, scopedSlots }));
      }

      if (this.loading) children.push(this.genLoading());
      return children;
    },

    isSelectable(item) {
      if (this.itemSelectable) {
        return this.itemSelectable(item);
      }

      return getObjectValueByPath(item, this.selectableKey) !== false;
    },
  },
};
</script>
