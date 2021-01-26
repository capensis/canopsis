<template lang="pug">
  div
    v-textarea(
      v-field="form.comment",
      v-validate="'required'",
      :label="$t('modals.createPause.comment')",
      :error-messages="errors.collect('comment')",
      name="comment"
    )
    pbehavior-reasons-field(v-field="form.reason")
</template>

<script>
import PbehaviorReasonsField from '@/components/other/pbehavior/reasons/partials/pbehavior-reasons-field.vue';

export default {
  components: { PbehaviorReasonsField },
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
