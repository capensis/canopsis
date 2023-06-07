import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('tag');

export const entitiesTagMixin = {
  computed: {
    ...mapGetters({
      tags: 'items',
      tagsPending: 'pending',
      tagsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchTagsList: 'fetchList',
      createTag: 'create',
      updateTag: 'update',
      removeTag: 'remove',
      bulkRemoveTags: 'bulkRemove',
    }),
  },
};
