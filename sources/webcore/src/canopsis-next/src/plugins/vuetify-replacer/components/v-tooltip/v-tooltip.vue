<script>
import VTooltip from 'vuetify/es5/components/VTooltip';
import { getSlotType } from 'vuetify/lib/util/helpers';

export default {
  extends: VTooltip,
  methods: {
    mouseEnterHandler(e) {
      this.getActivator(e);
      this.runDelay('open');
    },

    /**
     * We've updated here point is user will mouse leave from activator to default slot we will not hide the tooltip
     *
     * @param {MouseEvent} e
     */
    mouseLeaveHandler(e) {
      if (this.$refs.activator.contains(e.relatedTarget) || this.$refs.content.contains(e.relatedTarget)) {
        return;
      }

      this.getActivator(e);
      this.runDelay('close');
    },

    genActivator() {
      const listeners = this.disabled ? {} : {
        mouseenter: this.mouseEnterHandler,
        mouseleave: this.mouseLeaveHandler,
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

  /**
   * We've added mouseleave listener for tooltip default slot for resolving the problem which was described above
   *
   * @param {Function} h
   * @returns {*}
   */
  render: function render(h) {
    const listeners = {
      mouseleave: this.mouseLeaveHandler,
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
