<template lang="pug">
  div
    v-stepper.state-setting-form(v-model="stepper")
      v-stepper-header
        v-stepper-step(
          :complete="stepper > steps.BASICS",
          :step="steps.BASICS",
          :rules="[() => !hasBasicsFormAnyError]",
          editable
        ) {{ $t('common.basics') }}
        v-divider
        v-stepper-step(
          :complete="stepper > steps.ENTITIES",
          :step="steps.ENTITIES",
          :rules="[() => !hasEntitiesFormAnyError]",
          editable
        ) {{ $t('common.entities') }}
        v-divider
        v-stepper-step(
          :complete="stepper > steps.CONDITIONS",
          :step="steps.CONDITIONS",
          :rules="[() => !hasConditionsFormAnyError]",
          editable
        ) {{ $t('common.conditions') }}
      v-stepper-items
        v-stepper-content(:step="steps.BASICS")
          v-layout(column)
            v-layout(row)
              v-flex(xs7)
                c-name-field.mr-2(
                  v-field="form.title",
                  :label="$t('common.title')",
                  name="title",
                  required
                )
              v-flex(xs3)
                c-priority-field.mx-2(
                  v-field="form.priority",
                  required
                )
              v-flex(xs2)
                c-enabled-field.ml-2
            state-setting-compute-method-field(v-field="form.method")
            state-setting-target-entity-state-field(v-field="form", label="critical")
        v-stepper-content(:step="steps.ENTITIES")
          span ENTITIES
        v-stepper-content(:step="steps.CONDITIONS")
          span CONDITIONS
</template>

<script>
import { STATE_SETTING_COMPUTE_METHODS, STATE_SETTING_METHODS } from '@/constants';

import { formMixin } from '@/mixins/form';

import CNameField from '@/components/forms/fields/c-name-field.vue';
import CEnabledField from '@/components/forms/fields/c-enabled-field.vue';
import CPriorityField from '@/components/forms/fields/c-priority-field.vue';

import StateSettingComputeMethodField from './fields/state-setting-compute-method-field.vue';
import StateSettingTargetEntityStateField from './fields/state-setting-target-entity-state-field.vue';

export default {
  inject: ['$validator'],
  components: {
    CPriorityField,
    CEnabledField,
    CNameField,
    StateSettingComputeMethodField,
    StateSettingTargetEntityStateField,
  },
  mixins: [formMixin],
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
  data() {
    return {
      stepper: 1,
      hasBasicsFormAnyError: false,
      hasEntitiesFormAnyError: false,
      hasConditionsFormAnyError: false,
    };
  },
  computed: {
    STATE_SETTING_COMPUTE_METHODS() {
      return STATE_SETTING_COMPUTE_METHODS;
    },
    steps() {
      return {
        BASICS: 1,
        ENTITIES: 2,
        CONDITIONS: 3,
      };
    },

    isWorstOfShareMethod() {
      return this.form.method === STATE_SETTING_METHODS.worstOfShare;
    },
  },

  mounted() {
    /*    this.$watch(() => this.$refs.generalForm.hasAnyError, (value) => {
      this.hasBasicsFormAnyError = value;
    });

    this.$watch(() => this.$refs.infosForm.hasAnyError, (value) => {
      this.hasEntitiesFormAnyError = value;
    });

    this.$watch(() => this.$refs.patternsForm.hasAnyError, (value) => {
      this.hasConditionsFormAnyError = value;
    }); */
  },
};
</script>

<style lang="scss">
.state-setting-form {
  background-color: transparent !important;
}
</style>
