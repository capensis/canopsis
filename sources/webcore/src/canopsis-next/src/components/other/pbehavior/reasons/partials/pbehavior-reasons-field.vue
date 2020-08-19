<template lang="pug">
  v-select(
    v-field="reason",
    v-validate="'required'",
    :label="$t('modals.createPbehavior.steps.general.fields.reason')",
    :loading="pbehaviorReasonsPending",
    :items="pbehaviorReasons",
    :error-messages="errors.collect('reason')",
    return-object,
    name="reason",
    item-text="name",
    item-value="_id"
  )
</template>

<script>
import pbehaviorReasonsMixin from '@/mixins/entities/pbehavior/reasons';

export default {
  mixins: [pbehaviorReasonsMixin],
  model: {
    prop: 'reason',
    event: 'input',
  },
  inject: ['$validator'],
  props: {
    reason: {
      type: [Object, String],
      default: '',
    },
  },
};
</script>
