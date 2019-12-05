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
            dynamic-info-general-form(v-field="form.general")
      v-stepper-content(step="2")
        v-card
          v-card-text
            dynamic-info-infos-form(v-field="form.infos")
      v-stepper-content(step="3")
        v-card
          v-card-text
            dynamic-info-patterns-form(v-field="form.patterns")
</template>

<script>
import DynamicInfoGeneralForm from './partials/dynamic-info-general-form.vue';
import DynamicInfoInfosForm from './partials/dynamic-info-infos-form.vue';
import DynamicInfoPatternsForm from './partials/dynamic-info-patterns-form.vue';

export default {
  components: {
    DynamicInfoGeneralForm,
    DynamicInfoInfosForm,
    DynamicInfoPatternsForm,
  },
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
