<template lang="pug">
  action-btn(
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

import ActionBtn from '@/components/common/buttons/action-btn.vue';

export default {
  components: { ActionBtn },
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

