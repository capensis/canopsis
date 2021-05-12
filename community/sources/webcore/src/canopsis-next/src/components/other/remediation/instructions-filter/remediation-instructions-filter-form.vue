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
      v-flex(md3, xs6)
        v-checkbox(
          :input-value="form.auto",
          :label="$t('remediationInstructions.types.automatic')",
          :disabled="form.all || hasAnyAnotherOppositeFilterWithAuto",
          color="primary",
          @change="changeType('auto', $event)"
        )
      v-checkbox(
        :input-value="form.manual",
        :label="$t('remediationInstructions.types.manual')",
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
        :disabled="form.all || (form.auto && form.manual)",
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
import { find, pick } from 'lodash';

import { MAX_LIMIT } from '@/constants';

import { isRemediationInstructionIntersectsWithFilterByType } from '@/helpers/forms/remediation-instruction-filter';

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
      return (this.form.all || this.form.manual || this.form.auto) ? {} : { required: true };
    },

    hasAnyAnotherOppositeFilter() {
      return this.filters.some(filter => this.form.with !== filter.with);
    },

    hasAnyAnotherOppositeFilterWithAuto() {
      return this.filters.some(filter =>
        this.form.with !== filter.with && (filter.auto || filter.all));
    },

    hasAnyAnotherOppositeFilterWithManual() {
      return this.filters.some(filter =>
        this.form.with !== filter.with && (filter.manual || filter.all));
    },

    preparedRemediationInstructions() {
      return this.remediationInstructions.reduce((acc, instruction) => {
        if (isRemediationInstructionIntersectsWithFilterByType(this.form, instruction)) {
          return acc;
        }

        const filtersSomeComparator = filter =>
          this.form.with !== filter.with
          && (
            (filter.all || find(filter.instructions, { _id: instruction._id }))
            || (isRemediationInstructionIntersectsWithFilterByType(filter, instruction))
          );

        const instructionAlreadyInForm = find(this.form.instructions, { _id: instruction._id });
        const disabled = !instructionAlreadyInForm && this.filters.some(filtersSomeComparator);

        acc.push(disabled ? { ...instruction, disabled } : instruction);

        return acc;
      }, []);
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

    changeType(key, value) {
      const newForm = {
        ...this.form,

        [key]: value,
      };

      newForm.instructions =
        newForm.instructions
          .filter(instruction => !isRemediationInstructionIntersectsWithFilterByType(newForm, instruction));

      this.updateModel(newForm);
    },

    changeInstructions(instructions) {
      this.updateField('instructions', instructions.map(instruction => pick(instruction, ['_id', 'name', 'type'])));
    },
  },
};
</script>
