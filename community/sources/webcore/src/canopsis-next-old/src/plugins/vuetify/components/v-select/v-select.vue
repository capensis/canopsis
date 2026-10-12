<script>
import VSelect from 'vuetify/lib/components/VSelect/VSelect';

export default {
  extends: VSelect,
  props: {
    alwaysDirty: {
      type: Boolean,
      default: false,
    },
    hideInput: {
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

    /**
     * We've refactored this function to be able to hide the input field
     */
    genDefaultSlot() {
      const selections = this.genSelections();
      const input = this.genInput(); // If the return is an empty array
      // push the input

      if (Array.isArray(selections)) {
        if (!this.hideInput) {
          selections.push(input); // Otherwise push it into children
        }
      } else {
        selections.children = selections.children || [];

        if (!this.hideInput) {
          selections.children.push(input);
        }
      }

      return [
        this.genFieldset(),
        this.$createElement('div', {
          staticClass: 'v-select__slot',
          directives: this.directives,
        }, [
          this.genLabel(),
          this.prefix ? this.genAffix('prefix') : null,
          selections,
          this.suffix ? this.genAffix('suffix') : null,
          this.genClearIcon(),
          this.genIconSlot(),
          this.genHiddenInput(),
        ]),
        this.genMenu(),
        this.genProgress(),
      ];
    },

  },
};
</script>
