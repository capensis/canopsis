import dedupeModelListeners from 'vuetify/es5/util/dedupeModelListeners';
import rebuildSlots from 'vuetify/es5/util/rebuildFunctionalSlots';
import VCombobox from 'vuetify/es5/components/VCombobox';
import VAutocomplete from 'vuetify/es5/components/VAutocomplete';
import VOverflowBtn from 'vuetify/es5/components/VOverflowBtn';
import { deprecate } from 'vuetify/es5/util/console';

import VSelectOriginal from 'vuetify/es5/components/VSelect/VSelect';

const VSelect = VSelectOriginal.extend({
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
        { staticClass: 'v-select__selections' },
        this.genSelectionsChildren(),
      );
    },
  },
});

const wrapper = {
  functional: true,

  $_wrapperFor: VSelect,

  props: {
    // VAutoComplete
    /** @deprecated */
    autocomplete: Boolean,
    /** @deprecated */
    combobox: Boolean,
    multiple: Boolean,
    /** @deprecated */
    tags: Boolean,
    // VOverflowBtn
    /** @deprecated */
    editable: Boolean,
    /** @deprecated */
    overflow: Boolean,
    /** @deprecated */
    segmented: Boolean,
  },

  render(h, { props, data, slots, parent }) { // eslint-disable-line max-statements
    dedupeModelListeners(data);
    const children = rebuildSlots(slots(), h);

    if (props.autocomplete) {
      deprecate('<v-select autocomplete>', '<v-autocomplete>', wrapper, parent);
    }
    if (props.combobox) {
      deprecate('<v-select combobox>', '<v-combobox>', wrapper, parent);
    }
    if (props.tags) {
      deprecate('<v-select tags>', '<v-combobox multiple>', wrapper, parent);
    }

    if (props.overflow) {
      deprecate('<v-select overflow>', '<v-overflow-btn>', wrapper, parent);
    }
    if (props.segmented) {
      deprecate('<v-select segmented>', '<v-overflow-btn segmented>', wrapper, parent);
    }
    if (props.editable) {
      deprecate('<v-select editable>', '<v-overflow-btn editable>', wrapper, parent);
    }

    /* eslint-disable no-param-reassign */
    data.attrs = data.attrs || {};

    if (props.combobox || props.tags) {
      data.attrs.multiple = props.tags;
      return h(VCombobox, data, children);
    } if (props.autocomplete) {
      data.attrs.multiple = props.multiple;
      return h(VAutocomplete, data, children);
    } if (props.overflow || props.segmented || props.editable) {
      data.attrs.segmented = props.segmented;
      data.attrs.editable = props.editable;
      return h(VOverflowBtn, data, children);
    }
    data.attrs.multiple = props.multiple;
    /* eslint-enable */
    return h(VSelect, data, children);
  },
};

export default wrapper;
