<template>
  <c-action-btn
    v-if="hasReadAnyPbehaviorAccess"
    :tooltip="$t('common.pbehaviorList')"
    icon="list"
    @click="showPbehaviorsListModal"
  />
</template>

<script>
import { CRUD_ACTIONS, MODALS } from '@/constants';

import { permissionsTechnicalExploitationPbehaviorMixin } from '@/mixins/permissions/technical/exploitation/pbehavior';

export default {
  mixins: [
    permissionsTechnicalExploitationPbehaviorMixin,
  ],
  props: {
    entityId: {
      type: [Number, String],
      required: true,
    },
  },
  methods: {
    showPbehaviorsListModal() {
      this.$modals.show({
        name: MODALS.pbehaviorList,
        config: {
          entityId: this.entityId,
          availableActions: [CRUD_ACTIONS.delete, CRUD_ACTIONS.update],
        },
      });
    },
  },
};
</script>
