<template lang="pug">
  div
    v-layout(row)
      v-radio-group(
        v-field="form.with",
        hide-details
      )
        v-radio(
          :label="$t('remediationInstructionsFilters.fields.with')",
          :value="true",
          color="primary"
        )
        v-radio(
          :label="$t('remediationInstructionsFilters.fields.without')",
          :value="false",
          color="primary"
        )
    v-layout(row)
      v-switch(
        :input-value="form.all",
        :label="$t('remediationInstructionsFilters.fields.selectAll')",
        :disabled="hasAnyAnotherOppositeFilter",
        color="primary",
        @change="changeSelectedAll"
      )
    v-layout(row)
      v-select(
        v-field="form.instructions",
        v-validate="selectValidationRules",
        :items="preparedRemediationInstructions",
        :loading="remediationInstructionsPending",
        :disabled="form.all",
        :label="$t('remediationInstructionsFilters.fields.selectedInstructions')",
        :error-messages="errors.collect('instructions')",
        itemText="name",
        itemValue="name",
        name="instructions",
        multiple,
        clearable
      )
</template>

<script>
import entitiesRemediationInstructionsMixin from '@/mixins/entities/remediation/instructions';
import formMixin from '@/mixins/form/object';

export default {
  inject: ['$validator'],
  mixins: [formMixin, entitiesRemediationInstructionsMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    filters: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    selectValidationRules() {
      return this.form.all ? {} : { required: true };
    },

    hasAnyAnotherOppositeFilter() {
      return this.filters.some(filter => this.form.with !== filter.with);
    },

    preparedRemediationInstructions() {
      return this.remediationInstructions.map((instruction) => {
        const disabled = !this.form.instructions.includes(instruction.name)
          && this.filters.some(filter =>
            this.form.with !== filter.with && (filter.all || filter.instructions.includes(instruction.name)));

        if (disabled) {
          return { ...instruction, disabled };
        }

        return instruction;
      });
    },
  },
  mounted() {
    this.fetchRemediationInstructionsList({ limit: 10000 });
  },
  methods: {
    changeSelectedAll(all) {
      const newForm = { ...this.form, all };

      if (all) {
        newForm.instructions = [];
      }

      this.updateModel(newForm);
    },
  },
};
</script>
