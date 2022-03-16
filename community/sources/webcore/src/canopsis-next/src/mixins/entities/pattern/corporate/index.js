import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('pattern/corporate');

export const entitiesCorporatePatternsMixin = {
  computed: {
    ...mapGetters({
      corporatePatternsMeta: 'meta',
      corporatePatternsPending: 'pending',
      corporatePatterns: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchCorporatePatternsList: 'fetchList',
      fetchCorporatePatternsListWithPreviousParams: 'fetchListWithPreviousParams',
    }),
  },
};
