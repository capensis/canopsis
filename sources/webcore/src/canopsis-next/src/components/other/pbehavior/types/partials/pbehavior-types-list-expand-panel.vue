<template lang="pug">
  v-layout.pa-3.secondary.lighten-1(column)
    v-layout(row)
      v-text-field.theme--dark(
        :value="pbehaviorType.description",
        :label="$t('common.description')",
        disabled
      )
    v-layout(row)
      v-text-field.theme--dark(
        :value="pbehaviorType.type",
        :label="$t('modals.createPbehaviorType.fields.type')",
        disabled
      )
</template>

<script>
import entitiesPbehaviorTypesMixin from '@/mixins/entities/pbehavior/types';
import { MODALS } from '@/constants';

export default {
  mixins: [entitiesPbehaviorTypesMixin],
  props: {
    pbehaviorType: {
      type: Object,
      default: () => ({}),
    },
  },
  methods: {
    showEditPbehaviorTypeModal(pbehaviorType) {
      this.$modals.show({
        name: MODALS.createPbehaviorType,
        config: {
          pbehaviorType,
          action: async (newPbehaviorType) => {
            await this.updatePbehaviorType({
              data: newPbehaviorType,
              id: pbehaviorType._id,
            });
            this.$emit('refresh');
          },
        },
      });
    },

    showRemovePbehaviorTypeModal(pbehaviorTypeId) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removePbehaviorType({ id: pbehaviorTypeId });
            this.$emit('refresh');
          },
        },
      });
    },
  },
};
</script>
