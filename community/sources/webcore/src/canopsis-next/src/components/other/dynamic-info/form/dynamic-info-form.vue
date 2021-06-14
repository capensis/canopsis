<template lang="pug">
  v-stepper(v-model="stepper", non-linear)
    v-stepper-header
      v-stepper-step.py-0(
        :complete="stepper > steps.GENERAL",
        :step="steps.GENERAL",
        :rules="[() => !hasGeneralFormAnyError]",
        editable
      ) {{ $t('common.general') }}
        small(v-if="hasGeneralFormAnyError") {{ $t('modals.createDynamicInfo.errors.invalid') }}
      v-divider
      v-stepper-step.py-0(
        :complete="stepper > steps.INFOS",
        :step="steps.INFOS",
        :rules="[() => !hasInfosFormAnyError]",
        editable
      ) {{ $t('modals.createDynamicInfo.steps.infos.title') }}
        small(v-if="hasInfosFormAnyError") {{ $t('modals.createDynamicInfo.errors.invalid') }}
      v-divider
      v-stepper-step.py-0(
        :complete="stepper > steps.PATTERNS",
        :step="steps.PATTERNS",
        :rules="[() => !hasPatternsFormAnyError]",
        editable
      ) {{ $t('modals.createDynamicInfo.steps.patterns.title') }}
        small(v-if="hasPatternsFormAnyError") {{ $t('modals.createDynamicInfo.errors.invalid') }}
    v-stepper-items
      v-stepper-content.pa-4(:step="steps.GENERAL")
        v-card
          v-card-text.pa-0
            dynamic-info-general-form(
              v-field="form",
              :is-disabled-id-field="isDisabledIdField",
              ref="generalForm"
            )
      v-stepper-content.pa-4(:step="steps.INFOS")
        v-card
          v-card-text.pa-0
            dynamic-info-infos-form(v-field="form.infos", ref="infosForm")
      v-stepper-content.pa-4(:step="steps.PATTERNS")
        v-card
          v-card-text.pa-0
            c-patterns-field(v-field="form.patterns", ref="patternsForm", alarm, entity, some-required)
</template>

<script>
import DynamicInfoGeneralForm from './partials/dynamic-info-general-form.vue';
import DynamicInfoInfosForm from './partials/dynamic-info-infos-form.vue';

export default {
  inject: ['$validator'],
  components: {
    DynamicInfoGeneralForm,
    DynamicInfoInfosForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      stepper: 1,
      hasGeneralFormAnyError: false,
      hasInfosFormAnyError: false,
      hasPatternsFormAnyError: false,
    };
  },
  computed: {
    steps() {
      return {
        GENERAL: 1,
        INFOS: 2,
        PATTERNS: 3,
      };
    },
  },
  mounted() {
    this.$watch(() => this.$refs.generalForm.hasAnyError, (value) => {
      this.hasGeneralFormAnyError = value;
    });

    this.$watch(() => this.$refs.infosForm.hasAnyError, (value) => {
      this.hasInfosFormAnyError = value;
    });

    this.$watch(() => this.$refs.patternsForm.hasAnyError, (value) => {
      this.hasPatternsFormAnyError = value;
    });
  },
};
</script>
