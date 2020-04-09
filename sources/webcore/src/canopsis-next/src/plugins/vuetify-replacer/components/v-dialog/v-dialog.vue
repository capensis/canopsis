<script>
import VDialog from 'vuetify/es5/components/VDialog';

import ThemeProvider from 'vuetify/es5/util/ThemeProvider';
import * as helpers from 'vuetify/es5/util/helpers';

import overlayableMixin from '../../mixins/overlayable';

export default {
  extends: VDialog,
  mixins: [overlayableMixin],
  props: {
    onClickOutside: {
      type: Function,
      default: null,
    },
  },
  /**
   * We've replaced render method because we've put changes for the `click-outside` directive value function
   *
   * @param h
   * @returns {*}
   */
  render: function render(h) {
    const children = [];
    const data = {
      class: this.classes,
      ref: 'dialog',
      directives: [{
        name: 'click-outside',
        value: () => {
          if (!this.onClickOutside || (this.onClickOutside && this.onClickOutside() !== false)) {
            this.isActive = false;
          }
        },
        args: {
          closeConditional: this.closeConditional,
          include: this.getOpenDependentElements,
        },
      }, { name: 'show', value: this.isActive }],
      on: {
        click: (e) => {
          e.stopPropagation();
        },
      },
    };

    if (!this.fullscreen) {
      data.style = {
        maxWidth: this.maxWidth === 'none' ? undefined : helpers.convertToUnit(this.maxWidth),
        width: this.width === 'auto' ? undefined : helpers.convertToUnit(this.width),
      };
    }

    children.push(this.genActivator());

    let dialog = h('div', data, this.showLazyContent(this.$slots.default));

    if (this.transition) {
      dialog = h('transition', {
        props: {
          name: this.transition,
          origin: this.origin,
        },
      }, [dialog]);
    }

    children.push(h('div', {
      class: this.contentClasses,
      attrs: { tabIndex: '-1', ...this.getScopeIdAttrs() },
      on: {
        keydown: this.onKeydown,
      },
      style: { zIndex: this.activeZIndex },
      ref: 'content',
    }, [this.$createElement(ThemeProvider, {
      props: {
        root: true,
        light: this.light,
        dark: this.dark,
      },
    }, [dialog])]));

    return h('div', {
      staticClass: 'v-dialog__container',
      style: {
        display: !this.hasActivator || this.fullWidth ? 'block' : 'inline-block',
      },
    }, children);
  },
};
</script>
