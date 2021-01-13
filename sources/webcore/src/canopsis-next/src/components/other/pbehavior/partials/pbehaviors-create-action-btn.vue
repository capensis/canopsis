<template lang="pug">
  action-btn(
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

