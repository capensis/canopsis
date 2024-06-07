<script>
import { omit } from 'lodash';
import VTab from 'vuetify/es5/components/VTabs/VTab';

export default {
  extends: VTab,
  render(h) {
    const link = this.generateRouteLink(this.classes);

    /**
     * We've added this code to avoid warning error when `nativeOn` listeners added on `div` element
     */
    const [tag, data] = this.disabled
      ? ['div', omit(link.data, ['on', 'nativeOn'])]
      : [link.tag, link.data];

    data.ref = 'link';

    return h('div', {
      staticClass: 'v-tabs__div',
    }, [h(tag, data, this.$slots.default)]);
  },
};
</script>
