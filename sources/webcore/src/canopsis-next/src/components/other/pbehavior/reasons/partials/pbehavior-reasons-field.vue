<template lang="pug">
  v-select(
    v-field="reason",
    :label="$t('modals.createPbehavior.steps.general.fields.reason')",
    :loading="pbehaviorReasonsPending",
    :items="pbehaviorReasons",
    :error-messages="errors.collect(name)",
    :name="name",
    item-text="name",
    item-value="_id",
    return-object
  )
</template>

<script>
import { MAX_LIMIT } from '@/constants';

import pbehaviorReasonsMixin from '@/mixins/entities/pbehavior/reasons';

export default {
  $_veeValidate: {
    value() {
      return this.value;
    },

    name() {
      return this.name;
    },
  },
  inject: ['$validator'],
  mixins: [pbehaviorReasonsMixin],
  model: {
    prop: 'reason',
    event: 'input',
  },
  props: {
    reason: {
      type: [Object, String],
      default: '',
    },
    name: {
      type: String,
      default: 'reason',
    },
  },
  mounted() {
    this.fetchPbehaviorReasonsList({
      params: { limit: MAX_LIMIT },
    });
  },
};
</script>
