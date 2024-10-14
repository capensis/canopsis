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
import { getPbehaviorNameByEntities } from '@/helpers/entities/pbehavior/form';

import { permissionsTechnicalExploitationPbehaviorMixin } from '@/mixins/permissions/technical/exploitation/pbehavior';

export default {
  inject: ['$system'],
  mixins: [
    permissionsTechnicalExploitationPbehaviorMixin,
  ],
  props: {
    entity: {
      type: Object,
      required: true,
    },
  },
  methods: {
    showCreatePbehaviorModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: createEntityIdPatternByValue(this.entity?._id),
          defaultName: getPbehaviorNameByEntities([this.entity], this.$system.timezone),
        },
      });
    },
  },
};
</script>
