<script>
import VTooltip from 'vuetify/es5/components/VTooltip';
import { getSlotType } from 'vuetify/lib/util/helpers';

export default {
  extends: VTooltip,
  methods: {
    genActivator: function genActivator() {
      const listeners = this.disabled ? {} : {
        mouseenter: (e) => {
          this.getActivator(e);
          this.runDelay('open');
        },
        mouseleave: (e) => {
          if (e.toElement && this.$refs.content.contains(e.toElement)) {
            return;
          }

          this.getActivator(e);
          this.runDelay('close');
        },
      };

      if (getSlotType(this, 'activator') === 'scoped') {
        const activator = this.$scopedSlots.activator({ on: listeners });
        this.activatorNode = activator;
        return activator;
      }
      return this.$createElement('span', {
        on: listeners,
        ref: 'activator',
      }, this.$slots.activator);
    },
  },
  render: function render(h) {
    const listeners = {
      mouseleave: (e) => {
        if (e.toElement && this.$refs.activator.contains(e.toElement)) {
          return;
        }

        this.getActivator(e);
        this.runDelay('close');
      },
    };

    const tooltip = h('div', this.setBackgroundColor(this.color, {
      on: listeners,
      staticClass: 'v-tooltip__content',
      class: {
        [this.contentClass]: true,
        menuable__content__active: this.isActive,
        'v-tooltip__content--fixed': this.activatorFixed,
      },
      style: this.styles,
      attrs: this.getScopeIdAttrs(),
      directives: [{
        name: 'show',
        value: this.isContentActive,
      }],
      ref: 'content',
    }), this.showLazyContent(this.$slots.default));

    return h(this.tag, {
      staticClass: 'v-tooltip',
      class: this.classes,
    }, [h('transition', {
      props: {
        name: this.computedTransition,
      },
    }, [tooltip]), this.genActivator()]);
  },
};
</script>
