<script>
import VSpeedDial from 'vuetify/lib/components/VSpeedDial';

export default {
  extends: VSpeedDial,
  props: {
    hideOnClickOutside: {
      type: Boolean,
      default: false,
    },
  },
  /**
   * We've made click-outside optional
   */
  render(h) {
    let children = [];
    const directives = this.hideOnClickOutside ? [{
      name: 'click-outside',
      value: () => this.isActive = false,
    }] : [];

    const data = {
      directives,

      class: this.classes,
      on: {
        click: () => this.isActive = !this.isActive,
      },
    };

    if (this.openOnHover) {
      data.on.mouseenter = () => this.isActive = true;

      data.on.mouseleave = () => this.isActive = false;
    }

    if (this.isActive) {
      let btnCount = 0;
      children = (this.$slots.default || []).map((slot, i) => {
        if (slot.tag && typeof slot.componentOptions !== 'undefined') {
          btnCount += 1;

          return h('div', {
            style: {
              transitionDelay: `${btnCount * 0.05}s`,
            },
            key: i,
          }, [slot]);
        }

        return { ...slot, key: i };
      });
    }

    const list = h('transition-group', {
      class: 'v-speed-dial__list',
      props: {
        name: this.transition,
        mode: this.mode,
        origin: this.origin,
        tag: 'div',
      },
    }, children);

    return h('div', data, [this.$slots.activator, list]);
  },
};
</script>
