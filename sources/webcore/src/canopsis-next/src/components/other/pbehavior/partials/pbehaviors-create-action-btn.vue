<template lang="pug">
  c-action-btn(
    v-if="hasCreateAnyPbehaviorAccess",
    :tooltip="$t('modals.createPbehavior.create.title')",
    icon="pause",
    @click="showCreatePbehaviorModal"
  )
</template>

<script>
import { MODALS } from '@/constants';

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
    showCreatePbehaviorModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          filter: {
            _id: { $in: [this.entityId] },
          },
        },
      });
    },
  },
};
</script>

