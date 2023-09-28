import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('theme');

export const entitiesThemesMixin = {
  computed: {
    ...mapGetters({
      themes: 'items',
      themesMeta: 'meta',
      themesPending: 'pending',
    }),
  },
  methods: {
    ...mapActions({
      fetchThemesList: 'fetchList',
      fetchThemesListWithPreviousParams: 'fetchListWithPreviousParams',
      createTheme: 'create',
      updateTheme: 'update',
      removeTheme: 'remove',
      bulkRemoveThemes: 'bulkRemove',
    }),
  },
};
