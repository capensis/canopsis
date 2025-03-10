<script>
import VChip from 'vuetify/lib/components/VChip/VChip';

/**
 * Checks if all nodes in the given array have empty text content.
 *
 * @param {Array} nodes - An array of nodes to check. Each node can have a `text` property and a `children` array.
 * @returns {boolean} - Returns `true` if all nodes and their children have empty text, otherwise `false`.
 */
const isEmptyText = (nodes = []) => nodes.every(node => !node?.text && isEmptyText(node?.children));

export default {
  extends: VChip,
  computed: {
    hasFillBorderIcon() {
      return !!this.$el?.querySelector?.('.v-icon--fill-border');
    },
  },
  render(h) {
    const children = [this.genContent()];
    const routeLink = this.generateRouteLink();
    let { data } = routeLink;

    data.class = {
      ...data.class,
      'v-chip--without-text': isEmptyText(this.$slots.default),
      'v-chip--with-fill-border-icon': this.hasFillBorderIcon,
    };

    data.directives.push({
      name: 'show',
      value: this.active,
    });

    data = this.setBackgroundColor(this.color, data);
    const color = this.textColor || (this.outlined && this.color);
    return h(routeLink.tag, this.setTextColor(color, data), children);
  },
};
</script>

<style lang="scss">
.v-chip.v-chip--without-text:has(.v-icon--fill-border) .v-chip__content {
  position: absolute;
  top: 0;
  left: 0;
  padding: 0;

  * {
    width: 100% !important;
    height: 100% !important;
  }
}
</style>
