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
      v-checkbox(
        v-field="form.auto",
        :label="$t(`remediationInstructions.types.automatic`)",
        :disabled="form.all",
        color="primary"
      )
      v-checkbox(
        v-field="form.manual",
        :label="$t('remediationInstructions.types.manual')",
        :disabled="form.all",
        color="primary"
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
        item-value="_id",
        name="instructions",
        multiple,
        clearable,
        return-object,
        @change="changeInstructions"
      )
</template>

<script>
import { MAX_LIMIT } from '@/constants';

import { formMixin } from '@/mixins/form';
import { entitiesRemediationInstructionsMixin } from '@/mixins/entities/remediation/instructions';

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
        newForm.manual = true;
        newForm.auto = true;
      }

      this.updateModel(newForm);
    },

    changeInstructions(instructions) {
      this.updateField('instructions', instructions.map(({ _id, name }) => ({ _id, name })));
    },
  },
};
</script>
