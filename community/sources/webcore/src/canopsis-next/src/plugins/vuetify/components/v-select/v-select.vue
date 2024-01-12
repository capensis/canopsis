<script>
import VSelect from 'vuetify/lib/components/VSelect/VSelect';

export default {
  extends: VSelect,
  props: {
    alwaysDirty: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isDirty() {
      return this.selectedItems.length > 0 || this.alwaysDirty;
    },
  },
  methods: {
    genSelectionsChildren() {
      if (this.$scopedSlots.selections) {
        return this.$scopedSlots.selections({ items: this.selectedItems });
      }

      let { length } = this.selectedItems;
      const children = new Array(length);
      let genSelection;

      if (this.$scopedSlots.selection) {
        genSelection = this.genSlotSelection;
      } else if (this.hasChips) {
        genSelection = this.genChipSelection;
      } else {
        genSelection = this.genCommaSelection;
      }

      // eslint-disable-next-line no-plusplus
      while (length--) {
        children[length] = genSelection(
          this.selectedItems[length],
          length,
          length === children.length - 1,
        );
      }

      return children;
    },

    genSelections() {
      return this.$createElement(
        'div',
        {
          staticClass: 'v-select__selections',
        },
        this.genSelectionsChildren(),
      );
    },
  },
};
</script>
