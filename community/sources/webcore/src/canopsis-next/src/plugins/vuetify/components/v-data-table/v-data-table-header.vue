<script>
import { VDataTableHeader } from 'vuetify/lib/components/VDataTable';
import mergeData from 'vuetify/lib/util/mergeData';
import dedupeModelListeners from 'vuetify/lib/util/dedupeModelListeners';
import rebuildSlots from 'vuetify/lib/util/rebuildFunctionalSlots';

import VDataTableHeaderDesktop from './v-data-table-header-desktop.vue';

export default {
  components: { VDataTableHeaderDesktop },
  extends: VDataTableHeader,
  /**
   * We disable mobile version and added customized VDataTableHeaderDesktop
   */
  render(h, { props, data, slots }) {
    dedupeModelListeners(data);
    const children = rebuildSlots(slots(), h);
    const mergedData = mergeData(data, {
      props,
    });

    return h(VDataTableHeaderDesktop, mergedData, children);
  },
};
</script>
