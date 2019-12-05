<template lang="pug">
  v-stepper(v-model="stepper", non-linear)
    v-stepper-header
      v-stepper-step.py-0(
        :complete="stepper > 1",
        step="1",
        editable
      ) {{ $t('modals.createDynamicInfo.steps.general.title') }}
      v-divider
      v-stepper-step.py-0(
        :complete="stepper > 2",
        step="2",
        editable
      ) {{ $t('modals.createDynamicInfo.steps.infos.title') }}
      v-divider
      v-stepper-step.py-0(
        :complete="stepper > 3",
        step="3",
        editable
      ) {{ $t('modals.createDynamicInfo.steps.patterns.title') }}
    v-stepper-items
      v-stepper-content(step="1")
        v-card
          v-card-text
            div General form
      v-stepper-content(step="2")
        v-card
          v-card-text
            div Info form
      v-stepper-content(step="3")
        v-card
          v-card-text
            div Patterns form
</template>

<script>
export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  inject: ['$validator'],
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      stepper: 1,
    };
  },
};
</script>
