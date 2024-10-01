<template>
  <v-layout wrap>
    <v-flex xs5>
      <v-checkbox
        v-validate
        v-field="duration.enabled"
        :error-messages="errors.collect(enabledFieldName)"
        :name="enabledFieldName"
        color="primary"
      >
        <template #label="">
          {{ label }}
          <c-help-icon
            v-if="helpText"
            :text="helpText"
            icon-class="ml-2"
            color="info"
            max-width="300"
            top
          />
        </template>
      </v-checkbox>
    </v-flex>
    <v-flex xs4>
      <c-duration-field
        v-field="duration"
        :units-label="$t('common.unit')"
        :disabled="!duration.enabled"
        :required="duration.enabled"
        :units="timeUnits"
        :name="name"
        :min="min"
        @input="validate"
      />
    </v-flex>
    <v-flex xs9>
      <v-messages
        :value="errors.collect(name)"
        color="error"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { computed, toRef } from 'vue';

import { useEnabledDurationField } from '@/components/forms/fields/duration/hooks/enabled-duration-field';

export default {
  inject: ['$validator'],
  model: {
    prop: 'duration',
    event: 'input',
  },
  props: {
    duration: {
      type: Object,
      required: true,
    },
    label: {
      type: String,
      required: true,
    },
    helpText: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      required: false,
    },
    units: {
      type: Array,
      required: false,
    },
    after: {
      type: Object,
      required: false,
    },
  },
  setup(props) {
    const enabledFieldName = computed(() => `${props.name}.enabled`);

    const {
      timeUnits,
      min,
      validate,
    } = useEnabledDurationField({
      duration: toRef(props, 'duration'),
      name: toRef(props, 'name'),
      units: toRef(props, 'units'),
      after: toRef(props, 'after'),
    });

    return {
      enabledFieldName,
      timeUnits,
      min,
      validate,
    };
  },
};
</script>
