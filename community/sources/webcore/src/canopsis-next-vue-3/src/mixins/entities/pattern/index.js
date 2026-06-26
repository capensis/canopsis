import { usePattern } from '@/hooks/store/modules/pattern';

const patternComputed = [
  'patternsMeta',
  'patternsPending',
  'patterns',
].reduce((acc, key) => {
  acc[key] = () => usePattern()[key].value;

  return acc;
}, {});

export const entitiesPatternsMixin = {
  computed: patternComputed,
  methods: {
    fetchPatternsList(...args) {
      return usePattern().fetchPatternsList(...args);
    },

    createPattern(...args) {
      return usePattern().createPattern(...args);
    },

    updatePattern(...args) {
      return usePattern().updatePattern(...args);
    },

    removePattern(...args) {
      return usePattern().removePattern(...args);
    },

    bulkRemovePatterns(...args) {
      return usePattern().bulkRemovePatterns(...args);
    },

    fetchPatternsListWithPreviousParams(...args) {
      return usePattern().fetchPatternsListWithPreviousParams(...args);
    },

    fetchPatternsListWithoutStore(...args) {
      return usePattern().fetchPatternsListWithoutStore(...args);
    },
  },
};
