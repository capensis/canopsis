import { mapActions, mapState, mapStores } from 'pinia';

import { useNavigationStore } from '@/stores/modules/navigation';

export const layoutNavigationEditingModeMixin = {
  computed: {
    ...mapStores(useNavigationStore),
    ...mapState(useNavigationStore, {
      isNavigationEditingMode: 'isEditingMode',
    }),
  },
  methods: {
    ...mapActions(useNavigationStore, {
      toggleNavigationEditingMode: 'toggleEditingMode',
    }),
  },
};
