<template lang="pug">
  v-stepper(v-model="stepper", non-linear)
    v-stepper-header
      v-stepper-step.py-0(
        :complete="stepper > steps.GENERAL",
        :step="steps.GENERAL",
        :rules="[() => !hasGeneralFormAnyError]",
        editable
      ) {{ $t('modals.createDynamicInfo.steps.general.title') }}
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
      v-stepper-content(:step="steps.GENERAL")
        v-card
          v-card-text
            dynamic-info-general-form(v-field="form.general", ref="generalForm")
      v-stepper-content(:step="steps.INFOS")
        v-card
          v-card-text
            dynamic-info-infos-form(v-field="form.infos", ref="infosForm")
      v-stepper-content(:step="steps.PATTERNS")
        v-card
          v-card-text
            dynamic-info-patterns-form(v-field="form.patterns", ref="patternsForm")
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
  watch: {
    stepper(value) {
      if (value === this.steps.PATTERNS) {
        this.$refs.patternsForm.callTabsUpdateTabsMethod();
      }
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
