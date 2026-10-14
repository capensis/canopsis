<template>
  <v-select
    :value="value"
    :label="label || $t('common.type')"
    :items="fieldPbehaviorTypes"
    :loading="fieldPbehaviorTypesPending"
    :multiple="multiple"
    :chips="chips"
    :deletable-chips="chips"
    :small-chips="chips"
    item-text="name"
    item-value="_id"
    disabled
  />
</template>

<script>
import { entitiesFieldPbehaviorFieldTypeMixin } from '@/mixins/entities/pbehavior/types-field';

export default {
  mixins: [entitiesFieldPbehaviorFieldTypeMixin],
  props: {
    value: {
      type: [Object, String, Array],
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    chips: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
