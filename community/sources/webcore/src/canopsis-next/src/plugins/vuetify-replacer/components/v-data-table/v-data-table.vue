<script>
import { VDataTable } from 'vuetify/es5/components/VDataTable';
import { VCheckbox } from 'vuetify/es5/components/VCheckbox';

export default {
  extends: VDataTable,
  props: {
    isDisabledItem: {
      type: Function,
      default: item => !item,
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
  },
};
</script>
