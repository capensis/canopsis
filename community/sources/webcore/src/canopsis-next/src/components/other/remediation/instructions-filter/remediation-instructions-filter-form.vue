<template lang="pug">
  div
    v-layout(row)
      v-radio-group(v-field="form.with", name="with", hide-details)
        v-radio(
          :label="$t('remediationInstructionsFilters.with')",
          :value="true",
          color="primary"
        )
        v-radio(
          :label="$t('remediationInstructionsFilters.without')",
          :value="false",
          color="primary"
        )
    v-layout(row)
      v-switch(
        :input-value="form.all",
        :label="$t('remediationInstructionsFilters.selectAll')",
        :disabled="hasAnyAnotherOppositeFilter",
        color="primary",
        @change="changeSelectedAll"
      )
    v-layout(row)
      v-flex(md3, xs6)
        v-checkbox(
          :input-value="form.auto",
          :label="$t(`remediationInstructions.types.${$constants.REMEDIATION_INSTRUCTION_TYPES.auto}`)",
          :disabled="form.all || hasAnyAnotherOppositeFilterWithAuto",
          color="primary",
          @change="changeType('auto', $event)"
        )
      v-checkbox(
        :input-value="form.manual",
        :label="$t(`remediationInstructions.types.${$constants.REMEDIATION_INSTRUCTION_TYPES.manual}`)",
        :disabled="form.all || hasAnyAnotherOppositeFilterWithManual",
        color="primary",
        @change="changeType('manual', $event)"
      )
    v-layout(row)
      v-select(
        v-validate="selectValidationRules",
        :value="form.instructions",
        :items="preparedRemediationInstructions",
        :loading="remediationInstructionsPending",
        :disabled="isAll",
        :label="$t('remediationInstructionsFilters.selectedInstructions')",
        :error-messages="errors.collect('instructions')",
        item-text="name",
        item-value="_id",
        name="instructions",
        multiple,
        clearable,
        return-object,
        @change="changeInstructions"
      )
        c-help-icon(
          slot="append-outer",
          :text="$t('remediationInstructionsFilters.selectedInstructionsHelp')",
          color="grey darken-1",
          icon="help",
          left
        )
</template>

<script>
import { find, pick } from 'lodash';

import { MAX_LIMIT } from '@/constants';

import { isRemediationInstructionIntersectsWithFilterByType } from '@/helpers/forms/remediation-instruction-filter';

import { formMixin } from '@/mixins/form';
import { entitiesRemediationInstructionMixin } from '@/mixins/entities/remediation/instruction';

export default {
  inject: ['$validator'],
  mixins: [formMixin, entitiesRemediationInstructionMixin],
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
    isAll() {
      return this.form.all || (this.form.auto && this.form.manual);
    },

    selectValidationRules() {
      return (this.form.all || this.form.manual || this.form.auto) ? {} : { required: true };
    },

    hasAnyAnotherOppositeFilter() {
      return this.filters.some(filter => this.form.with !== filter.with);
    },

    hasAnyAnotherOppositeFilterWithAuto() {
      return this.filters.some(filter => this.form.with !== filter.with && (filter.auto || filter.all));
    },

    hasAnyAnotherOppositeFilterWithManual() {
      return this.filters.some(filter => this.form.with !== filter.with && (filter.manual || filter.all));
    },

    preparedRemediationInstructions() {
      return this.remediationInstructions.reduce((acc, instruction) => {
        if (isRemediationInstructionIntersectsWithFilterByType(this.form, instruction)) {
          return acc;
        }

        const filtersSomeComparator = filter => this.form.with !== filter.with
          && (
            (filter.all || find(filter.instructions, { _id: instruction._id }))
            || isRemediationInstructionIntersectsWithFilterByType(filter, instruction)
          );

        const instructionAlreadyInForm = find(this.form.instructions, { _id: instruction._id });
        const disabled = !instructionAlreadyInForm && this.filters.some(filtersSomeComparator);

        acc.push(disabled ? { ...instruction, disabled } : instruction);

        return acc;
      }, []);
    },
  },
  mounted() {
    this.fetchRemediationInstructionsList({ params: { limit: MAX_LIMIT } });
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

    changeType(key, value) {
      const newForm = {
        ...this.form,

        [key]: value,
      };

      newForm.instructions = newForm.instructions
        .filter(instruction => !isRemediationInstructionIntersectsWithFilterByType(newForm, instruction));

      this.updateModel(newForm);
    },

    changeInstructions(instructions) {
      this.updateField('instructions', instructions.map(instruction => pick(instruction, ['_id', 'name', 'type'])));
    },
  },
};
</script>
