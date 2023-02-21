import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('infos');

export const entitiesInfosMixin = {
  computed: {
    ...mapGetters({
      alarmInfos: 'alarmInfos',
      alarmInfosRules: 'alarmInfosRules',
      entityInfos: 'entityInfos',
      infosPending: 'pending',
    }),
  },
  methods: {
    ...mapActions({
      fetchInfos: 'fetch',
    }),
  },
};
