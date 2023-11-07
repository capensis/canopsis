<script>
import { VDataTable } from 'vuetify/lib/components/VDataTable';
import { getObjectValueByPath } from 'vuetify/lib/util/helpers';
import ExpandTransitionGenerator from 'vuetify/lib/components/transitions/expand-transition';

export default {
  extends: VDataTable,
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

        rows.push(this.$createElement('tr', { class: 'v-data-table__expanded__row' }, [transition]));
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
  },
};
</script>
