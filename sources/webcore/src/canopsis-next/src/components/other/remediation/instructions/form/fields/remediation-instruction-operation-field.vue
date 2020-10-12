<template lang="pug">
  v-layout(row)
    v-flex.pr-1(xs11)
      v-layout(row)
        expand-button.operation-expand(v-model="expanded")
        v-layout(column)
          v-text-field(
            v-field="operation.name",
            v-validate="'required'",
            :label="$t('common.name')",
            :error-messages="nameErrors",
            :name="nameFieldName",
            box
          )
          v-expand-transition(mode="out-in")
            v-layout(v-show="expanded", column)
              remediation-instruction-time-to-complete-field(
                v-field="operation.time_to_complete",
                :name="timeToCompleteFieldName"
              )
              text-editor-field(
                v-field="operation.description",
                :label="$t('common.description')",
                focusable
              )
    v-flex.mt-3(xs1)
      v-layout(justify-center)
        v-btn.ma-0(icon, small, @click.prevent="$emit('remove')")
          v-icon(color="error") delete
</template>

<script>
import formMixin from '@/mixins/form';

import TextEditorField from '@/components/forms/fields/text-editor-field.vue';
import ExpandButton from '@/components/other/buttons/expand-button.vue';

import RemediationInstructionTimeToCompleteField from './remediation-instruction-time-to-complete-field.vue';

export default {
  inject: ['$validator'],
  components: {
    TextEditorField,
    ExpandButton,
    RemediationInstructionTimeToCompleteField,
  },
  mixins: [formMixin],
  model: {
    prop: 'operation',
    event: 'input',
  },
  props: {
    operation: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      expanded: true,
    };
  },
  computed: {
    fieldName() {
      return this.operation.key ? this.operation.key : '';
    },

    nameFieldName() {
      return `${this.fieldName}.name`;
    },

    timeToCompleteFieldName() {
      return `${this.fieldName}.timeToComplete`;
    },

    descriptionFieldName() {
      return `${this.fieldName}.description`;
    },

    nameErrors() {
      return this.getErrors(this.nameFieldName, this.$t('common.name'));
    },

    descriptionErrors() {
      return this.getErrors(this.descriptionFieldName, this.$t('common.description'));
    },
  },
  methods: {
    getErrors(name, nameReplacer) {
      return this.errors.collect(name).map(error => error.replace(name, nameReplacer));
    },
  },
};
</script>

<style lang="scss" scoped>
  .operation-expand {
    margin: 24px 2px 0 2px !important;
    width: 20px;
    height: 20px;
  }
</style>
