<script>
import VCombobox from 'vuetify/lib/components/VCombobox';

export default {
  extends: VCombobox,
  props: {
    blurOnCreate: {
      type: Boolean,
      default: false,
    },
    forceSearching: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isSearching() {
      if (this.multiple) {
        return this.searchIsDirty;
      }

      return this.searchIsDirty && (this.forceSearching || this.internalSearch !== this.getText(this.selectedItem));
    },
  },
  methods: {
    onEnterDown(e) {
      e.preventDefault();

      /**
       * If has menu index, let v-select-list handle
       */
      if (this.getMenuIndex() > -1) {
        return;
      }

      this.$nextTick(() => {
        this.updateSelf();

        /**
         * We've added this condition for the closing menu when we've pressed enter to create new one
         */
        if (this.blurOnCreate) {
          this.blur();
        }
      });
    },
  },
};
</script>
