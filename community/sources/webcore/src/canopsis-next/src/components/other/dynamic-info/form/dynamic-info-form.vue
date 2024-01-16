<template>
  <v-stepper
    class="dynamic-info-form"
    v-model="stepper"
    non-linear
  >
    <v-stepper-header>
      <v-stepper-step
        class="py-0"
        :complete="stepper > steps.GENERAL"
        :step="steps.GENERAL"
        :rules="[() => !hasGeneralFormAnyError]"
        editable
      >
        {{ $t('common.general') }}
        <small v-if="hasGeneralFormAnyError">{{ $t('modals.createDynamicInfo.errors.invalid') }}</small>
      </v-stepper-step>
      <v-divider />
      <v-stepper-step
        class="py-0"
        :complete="stepper > steps.INFOS"
        :step="steps.INFOS"
        :rules="[() => !hasInfosFormAnyError]"
        editable
      >
        {{ $t('modals.createDynamicInfo.steps.infos.title') }}
        <small v-if="hasInfosFormAnyError">{{ $t('modals.createDynamicInfo.errors.invalid') }}</small>
      </v-stepper-step>
      <v-divider />
      <v-stepper-step
        class="py-0"
        :complete="stepper > steps.PATTERNS"
        :step="steps.PATTERNS"
        :rules="[() => !hasPatternsFormAnyError]"
        editable
      >
        {{ $t('modals.createDynamicInfo.steps.patterns.title') }}
        <small v-if="hasPatternsFormAnyError">{{ $t('modals.createDynamicInfo.errors.invalid') }}</small>
      </v-stepper-step>
    </v-stepper-header>
    <v-stepper-items>
      <v-stepper-content
        class="pa-0"
        :step="steps.GENERAL"
      >
        <dynamic-info-general-form
          class="pa-4"
          v-field="form"
          :is-disabled-id-field="isDisabledIdField"
          ref="generalForm"
        />
      </v-stepper-content>
      <v-stepper-content
        class="pa-0"
        :step="steps.INFOS"
      >
        <dynamic-info-infos-form
          class="pa-4"
          v-field="form.infos"
          ref="infosForm"
        />
      </v-stepper-content>
      <v-stepper-content
        class="pa-0"
        :step="steps.PATTERNS"
      >
        <dynamic-info-patterns-form
          class="pa-4"
          v-field="form.patterns"
          ref="patternsForm"
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

<style lang="scss">
.dynamic-info-form {
  background-color: transparent !important;
}
</style>
