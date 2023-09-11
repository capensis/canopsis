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
      createAlarmTag: 'create',
      updateAlarmTag: 'update',
      removeAlarmTag: 'remove',
      bulkRemoveAlarmTag: 'bulkRemove',
    }),

    getTagColor(tag) {
      const alarmTag = this.alarmTags.find(({ value }) => tag === value);

      return alarmTag?.color;
    },
  },
};
