<template lang="pug">
  v-btn.ml-5.primary(@click="addStep") {{ $t('remediationInstructions.addStep') }}
</template>

<script>
import formMixin from '@/mixins/form';
import formArrayMixin from '@/mixins/form/array';

export default {
  inject: ['$validator'],
  mixins: [formMixin, formArrayMixin],
  model: {
    prop: 'steps',
    event: 'input',
  },
  props: {
    steps: {
      type: Array,
      default: () => ([]),
    },
  },
  methods: {
    addStep() {},
  },
};
</script>
