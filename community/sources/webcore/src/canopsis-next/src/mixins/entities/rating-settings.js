import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('ratingSettings');

export const entitiesRatingSettingsMixin = {
  computed: {
    ...mapGetters({
      ratingSettings: 'items',
      ratingSettingsPending: 'pending',
      ratingSettingsMeta: 'meta',
      ratingSettingsUpdatedAt: 'updatedAt',
    }),
  },
  methods: {
    ...mapActions({
      fetchRatingSettingsListWithPreviousParams: 'fetchListWithPreviousParams',
      fetchRatingSettingsListWithoutStore: 'fetchListWithoutStore',
      fetchRatingSettingsList: 'fetchList',
      updateRatingSetting: 'update',
      bulkUpdateRatingSettings: 'bulkUpdate',
    }),
  },
};
