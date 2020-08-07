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

export default {
  mixins: [entitiesPbehaviorTypesMixin],
  props: {
    pbehaviorType: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
