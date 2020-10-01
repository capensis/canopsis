<template lang="pug">
  v-btn(
    v-if="hasCreateAnyPbehaviorAccess",
    icon,
    small,
    @click.stop="showCreatePbehaviorModal"
  )
    v-icon pause
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

