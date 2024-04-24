<script>
import VCombobox from 'vuetify/lib/components/VCombobox';

import VMenu from '../v-menu/v-menu.vue';

/**
 * Check is attach prop enabled
 *
 * @param {string | boolean} attach
 * @returns {boolean}
 */
const isAttached = attach => attach === '' // If used as a boolean prop (<v-menu attach>)
  || attach === true // If bound to a boolean (<v-menu :attach="true">)
  || attach === 'attach'; // If bound as boolean prop in pug (v-menu(attach))

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
    genCustomList() {
      return this.$scopedSlots.list?.(this.listData);
    },
    genMenu() {
      const props = this.$_menuProps;
      props.activator = this.$refs['input-slot'];

      if (!('attach' in props)) {
        if (isAttached(this.attach)) {
          // Attach to root el so that
          // menu covers prepend/append icons
          props.attach = this.$el;
        } else {
          props.attach = this.attach;
        }
      }

      return this.$createElement(VMenu, {
        attrs: {
          role: undefined,
        },
        props,
        on: {
          input: (val) => {
            this.isMenuActive = val;
            this.isFocused = val;
          },
          scroll: this.onScroll,
        },
        ref: 'menu',
      }, [
        this.$scopedSlots.list
          ? this.genCustomList()
          : this.genList(),
      ]);
    },

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
