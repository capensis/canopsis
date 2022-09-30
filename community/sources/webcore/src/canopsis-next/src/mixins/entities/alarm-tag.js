import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('alarmTag');

export const entitiesAlarmTagMixin = {
  computed: {
    ...mapGetters({
      alarmTags: 'items',
      getAlarmTagById: 'getItemById',
      alarmTagsPending: 'pending',
      alarmTagsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchAlarmTagsList: 'fetchList',
      removeAlarmTag: 'remove',
      createAlarmTag: 'create',
      updateAlarmTag: 'update',
    }),
  },
};
