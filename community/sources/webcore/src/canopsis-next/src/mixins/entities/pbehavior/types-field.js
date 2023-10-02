import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('pbehaviorTypes');

export const entitiesFieldPbehaviorFieldTypeMixin = {
  computed: {
    ...mapGetters({
      fieldPbehaviorTypes: 'fieldItems',
      fieldPbehaviorTypesPending: 'fieldPending',
    }),
  },
  methods: {
    ...mapActions({
      fetchFieldPbehaviorTypesList: 'fetchFieldList',
    }),
  },
};
