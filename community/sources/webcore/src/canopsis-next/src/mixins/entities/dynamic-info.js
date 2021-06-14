import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('dynamicInfo');

export const entitiesDynamicInfoMixin = {
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
      updateDynamicInfo: 'update',
      createDynamicInfo: 'create',
      removeDynamicInfo: 'remove',
    }),
  },
};
