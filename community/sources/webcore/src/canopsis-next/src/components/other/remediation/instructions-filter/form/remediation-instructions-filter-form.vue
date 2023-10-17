<template>
  <div>
    <v-radio-group
      v-field="form.with"
      :label="$t('remediation.instructionsFilter.filterByInstructions')"
      name="with"
      hide-details
    >
      <v-radio
        :label="$t('remediation.instructionsFilter.with')"
        :value="true"
        color="primary"
      />
      <v-radio
        :label="$t('remediation.instructionsFilter.without')"
        :value="false"
        color="primary"
      />
    </v-radio-group>
    <v-switch
      :input-value="form.all"
      :label="$t('remediation.instructionsFilter.selectAll')"
      :disabled="hasAnyAnotherOppositeFilter"
      color="primary"
      @change="changeSelectedAll"
    />
    <v-layout>
      <v-flex
        md3
        xs6
      >
        <v-checkbox
          :input-value="form.auto"
          :label="$t(`remediation.instruction.types.${$constants.REMEDIATION_INSTRUCTION_TYPES.auto}`)"
          :disabled="form.all || hasAnyAnotherOppositeFilterWithAuto"
          color="primary"
          @change="changeType('auto', $event)"
        />
      </v-flex>
      <v-checkbox
        :input-value="form.manual"
        :label="$t(`remediation.instruction.types.${$constants.REMEDIATION_INSTRUCTION_TYPES.manual}`)"
        :disabled="form.all || hasAnyAnotherOppositeFilterWithManual"
        color="primary"
        @change="changeType('manual', $event)"
      />
    </v-layout>
    <v-select
      v-validate="selectValidationRules"
      :value="form.instructions"
      :items="preparedRemediationInstructions"
      :loading="remediationInstructionsPending"
      :disabled="isAll"
      :label="$t('remediation.instructionsFilter.selectedInstructions')"
      :error-messages="errors.collect('instructions')"
      item-text="name"
      item-value="_id"
      name="instructions"
      multiple
      clearable
      return-object
      @change="changeInstructions"
    >
      <template #append-outer="">
        <c-help-icon
          :text="$t('remediation.instructionsFilter.selectedInstructionsHelp')"
          icon="help"
          left
        />
      </template>
    </v-select>
    <v-radio-group
      v-field="form.running"
      :label="$t('remediation.instructionsFilter.alarmsListDisplay')"
      name="with"
      hide-details
    >
      <v-radio
        :label="$t('remediation.instructionsFilter.allAlarms')"
        :value="null"
        color="primary"
      />
      <v-radio
        :label="showInProgressLabel"
        :value="true"
        color="primary"
      />
      <v-radio
        :label="hideInProgressLabel"
        :value="false"
        color="primary"
      />
    </v-radio-group>
  </div>
</template>

<script>
import { find, pick } from 'lodash';

import { MAX_LIMIT } from '@/constants';

import {
  isRemediationInstructionIntersectsWithFilterByType,
} from '@/helpers/entities/remediation/instruction-filter/form';

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

    showInProgressLabel() {
      const key = this.form.with
        ? 'remediation.instructionsFilter.showWithInProgress'
        : 'remediation.instructionsFilter.showWithoutInProgress';

      return this.$t(key);
    },

    hideInProgressLabel() {
      const key = this.form.with
        ? 'remediation.instructionsFilter.hideWithInProgress'
        : 'remediation.instructionsFilter.hideWithoutInProgress';

      return this.$t(key);
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
