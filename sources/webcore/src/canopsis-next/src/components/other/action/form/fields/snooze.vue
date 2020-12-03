<template lang="pug">
  div
    v-textarea(
      v-field="value.message",
      :label="$t('modals.createAction.fields.message')"
    )
    old-duration-field(v-field="value.duration")
</template>

<script>
import OldDurationField from '@/components/forms/fields/old-duration.vue';

export default {
  inject: ['$validator'],
  components: {
    OldDurationField,
  },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
};
</script>
