<template>
  <v-stepper
    v-model="stepper"
    non-linear
  >
    <v-stepper-header>
      <v-stepper-step
        :complete="stepper > steps.GENERAL"
        :step="steps.GENERAL"
        :rules="[() => !hasGeneralFormAnyError]"
        class="py-0"
        editable
      >
        {{ $t('common.general') }}
        <small v-if="hasGeneralFormAnyError">{{ $t('modals.createDynamicInfo.errors.invalid') }}</small>
      </v-stepper-step>
      <v-divider />
      <v-stepper-step
        :complete="stepper > steps.INFOS"
        :step="steps.INFOS"
        :rules="[() => !hasInfosFormAnyError]"
        class="py-0"
        editable
      >
        {{ $t('modals.createDynamicInfo.steps.infos.title') }}
        <small v-if="hasInfosFormAnyError">{{ $t('modals.createDynamicInfo.errors.invalid') }}</small>
      </v-stepper-step>
      <v-divider />
      <v-stepper-step
        :complete="stepper > steps.PATTERNS"
        :step="steps.PATTERNS"
        :rules="[() => !hasPatternsFormAnyError]"
        class="py-0"
        editable
      >
        {{ $t('modals.createDynamicInfo.steps.patterns.title') }}
        <small v-if="hasPatternsFormAnyError">{{ $t('modals.createDynamicInfo.errors.invalid') }}</small>
      </v-stepper-step>
    </v-stepper-header>
    <v-stepper-items>
      <v-stepper-content
        :step="steps.GENERAL"
        class="pa-0"
      >
        <dynamic-info-general-form
          v-field="form"
          ref="generalForm"
          :is-disabled-id-field="isDisabledIdField"
          class="pa-4"
        />
      </v-stepper-content>
      <v-stepper-content
        :step="steps.INFOS"
        class="pa-0"
      >
        <dynamic-info-infos-form
          v-field="form.infos"
          ref="infosForm"
          class="pa-4"
        />
      </v-stepper-content>
      <v-stepper-content
        :step="steps.PATTERNS"
        class="pa-0"
      >
        <dynamic-info-patterns-form
          v-field="form.patterns"
          ref="patternsForm"
          class="pa-4"
        />
      </v-stepper-content>
    </v-stepper-items>
  </v-stepper>
</template>

<script>
import DynamicInfoGeneralForm from './fields/dynamic-info-general-form.vue';
import DynamicInfoInfosForm from './fields/dynamic-info-infos-form.vue';
import DynamicInfoPatternsForm from './fields/dynamic-info-patterns-form.vue';

export default {
  inject: ['$validator'],
  components: {
    DynamicInfoGeneralForm,
    DynamicInfoInfosForm,
    DynamicInfoPatternsForm,
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
