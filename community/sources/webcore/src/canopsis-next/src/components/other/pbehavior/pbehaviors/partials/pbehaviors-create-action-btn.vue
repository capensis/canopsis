<template>
  <c-action-btn
    v-if="hasCreateAnyPbehaviorAccess"
    :tooltip="$t('modals.createPbehavior.create.title')"
    icon="pause"
    @click="showCreatePbehaviorModal"
  />
</template>

<script>
import { MODALS } from '@/constants';

import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';

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
    showCreatePbehaviorModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: createEntityIdPatternByValue(this.entityId),
        },
      });
    },
  },
};
</script>
