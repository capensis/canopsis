import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('dynamicInfo');

export default {
  computed: {
    ...mapGetters({
      dynamicInfosPending: 'pending',
      dynamicInfos: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchDynamicInfosList: 'fetchList',
      refreshDynamicInfosList: 'fetchListWithPreviousParams',
      createDynamicInfo: 'create',
      removeDynamicInfo: 'remove',
    }),
  },
};
