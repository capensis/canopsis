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
        v-validate="selectValidationRules",
        :value="form.instructions",
        :items="preparedRemediationInstructions",
        :loading="remediationInstructionsPending",
        :disabled="form.all",
        :label="$t('remediationInstructionsFilters.fields.selectedInstructions')",
        :error-messages="errors.collect('instructions')",
        item-text="name",
        item-value="name",
        name="instructions",
        multiple,
        clearable,
        @change="updateField('instructions', $event)"
      )
</template>

<script>
import { MAX_LIMIT } from '@/constants';

import formMixin from '@/mixins/form';
import entitiesRemediationInstructionsMixin from '@/mixins/entities/remediation/instructions';

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
        const filtersSomeComparator = filter => this.form.with !== filter.with
          && (filter.all || filter.instructions.includes(instruction.name));

        const instructionAlreadyInForm = this.form.instructions.includes(instruction.name);
        const disabled = !instructionAlreadyInForm && this.filters.some(filtersSomeComparator);

        if (disabled) {
          return { ...instruction, disabled };
        }

        return instruction;
      });
    },
  },
  mounted() {
    this.fetchRemediationInstructionsList({ limit: MAX_LIMIT });
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
