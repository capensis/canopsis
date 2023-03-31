<script>
import VCombobox from 'vuetify/es5/components/VCombobox';
import VSelect from 'vuetify/es5/components/VSelect/VSelect';

import { camelize } from 'vuetify/es5/util/helpers';
import { consoleWarn } from 'vuetify/es5/util/console';

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
    valuePreparer: {
      type: Function,
      default: v => v,
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
    genMenu() {
      const props = this.$_menuProps;
      props.activator = this.$refs['input-slot'];

      /**
       * Deprecate using menu props directly
       * TODO: remove (2.0)
       */
      const allProps = { ...VMenu.extends.options.props, ...VMenu.props };
      const inheritedProps = Object.keys(allProps);

      const deprecatedProps = Object.keys(this.$attrs).reduce((acc, attr) => {
        if (inheritedProps.includes(camelize(attr))) acc.push(attr);
        return acc;
      }, []);

      // eslint-disable-next-line no-restricted-syntax
      for (const prop of deprecatedProps) {
        props[camelize(prop)] = this.$attrs[prop];
      }

      if (process.env.NODE_ENV !== 'production') {
        if (deprecatedProps.length) {
          const multiple = deprecatedProps.length > 1;
          let replacement = deprecatedProps.reduce((acc, p) => {
            acc[camelize(p)] = this.$attrs[p];
            return acc;
          }, {});
          const propsJoined = deprecatedProps.map(p => `'${p}'`).join(', ');
          const separator = multiple ? '\n' : '\'';

          const onlyBools = Object.keys(replacement).every((prop) => {
            const propType = allProps[prop];
            const value = replacement[prop];
            return value === true || ((propType.type || propType) === Boolean && value === '');
          });

          if (onlyBools) {
            replacement = Object.keys(replacement).join(', ');
          } else {
            replacement = JSON.stringify(replacement, null, multiple ? 2 : 0)
              .replace(/"([^(")"]+)":/g, '$1:')
              .replace(/"/g, '\'');
          }

          consoleWarn(
            `${propsJoined} ${multiple ? 'are' : 'is'} deprecated, use `
            + `${separator}${onlyBools ? '' : ':'}menu-props="${replacement}"${separator} instead`,
            this,
          );
        }
      }

      /**
       * Attach to root el so that
       * menu covers prepend/append icons
       */
      if (isAttached(this.attach)) {
        props.attach = this.$el;
      } else {
        props.attach = this.attach;
      }

      return this.$createElement(VMenu, {
        props,
        on: {
          input: (val) => {
            this.isMenuActive = val;
            this.isFocused = val;
          },
        },
        ref: 'menu',
      }, [this.genList()]);
    },
    onEnterDown(e) {
      e.preventDefault();

      VSelect.options.methods.onEnterDown.call(this);

      /**
       * If has menu index, let v-select-list handle
       */
      if (this.getMenuIndex() > -1) {
        return;
      }

      this.updateSelf();

      /**
       * We've added this condition for the closing menu when we've pressed enter to create new one
       */
      if (this.blurOnCreate) {
        this.blur();
      }
    },

    setValue(value = this.internalSearch) {
      const oldValue = this.internalValue;
      const newValue = this.valuePreparer(value);

      this.internalValue = newValue;

      if (newValue !== oldValue) {
        this.$emit('change', newValue);
      }
    },
  },
};
</script>
