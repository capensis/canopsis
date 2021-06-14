import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('dynamicInfo');

export default {
  computed: {
    ...mapGetters({
      dynamicInfosPending: 'pending',
      dynamicInfos: 'items',
      dynamicInfosMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchDynamicInfosList: 'fetchList',
      refreshDynamicInfosList: 'fetchListWithPreviousParams',
      updateDynamicInfo: 'update',
      createDynamicInfo: 'create',
      removeDynamicInfo: 'remove',
    }),
  },
};
