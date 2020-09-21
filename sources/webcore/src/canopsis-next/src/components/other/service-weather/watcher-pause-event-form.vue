<template lang="pug">
  div
    v-textarea(
      v-field="form.comment",
      v-validate="'required'",
      :label="$t('modals.createPause.comment')",
      :error-messages="errors.collect('comment')",
      name="comment"
    )
    v-select(
      v-field="form.reason",
      v-validate="'required'",
      :label="$t('modals.createPause.reason')",
      :items="reasons",
      :error-messages="errors.collect('reason')",
      name="reason"
    )
</template>

<script>
import { PAUSE_REASONS } from '@/constants';

export default {
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
  computed: {
    reasons() {
      return Object.values(PAUSE_REASONS);
    },
  },
};
</script>
