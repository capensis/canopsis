<template lang="pug">
  c-action-btn(
    v-if="hasReadAnyPbehaviorAccess",
    :tooltip="$t('alarmList.actions.titles.pbehaviorList')",
    icon="list",
    @click="showPbehaviorsListModal"
  )
</template>

<script>
import { CRUD_ACTIONS, MODALS } from '@/constants';

import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';
import rightsTechnicalExploitationPbehaviorMixin from '@/mixins/rights/technical/exploitation/pbehavior';

export default {
  mixins: [
    entitiesPbehaviorMixin,
    rightsTechnicalExploitationPbehaviorMixin,
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

