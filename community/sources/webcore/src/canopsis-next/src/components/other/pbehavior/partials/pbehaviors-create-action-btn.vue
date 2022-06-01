<template lang="pug">
  c-action-btn(
    v-if="hasCreateAnyPbehaviorAccess",
    :tooltip="$t('modals.createPbehavior.create.title')",
    icon="pause",
    @click="showCreatePbehaviorModal"
  )
</template>

<script>
import { ENTITY_PATTERN_FIELDS, MODALS, PATTERN_CONDITIONS } from '@/constants';

import { entitiesPbehaviorMixin } from '@/mixins/entities/pbehavior';
import { permissionsTechnicalExploitationPbehaviorMixin } from '@/mixins/permissions/technical/exploitation/pbehavior';

export default {
  mixins: [
    entitiesPbehaviorMixin,
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
          entityPattern: [[{
            field: ENTITY_PATTERN_FIELDS.id,
            cond: {
              type: PATTERN_CONDITIONS.equal,
              value: this.entityId,
            },
          }]],
        },
      });
    },
  },
};
</script>
