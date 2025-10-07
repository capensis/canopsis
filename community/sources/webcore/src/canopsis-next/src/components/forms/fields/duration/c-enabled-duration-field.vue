<template>
  <v-layout wrap>
    <v-flex xs5>
      <component
        :is="component"
        v-field="duration.enabled"
        v-validate
        :error-messages="errors.collect(enabledFieldName)"
        :name="enabledFieldName"
        color="primary"
        class="pb-2"
      >
        <template #label="">
          <span class="full-width">
            {{ label }}
            <c-help-icon
              v-if="helpText"
              :text="helpText"
              icon-class="ml-2"
              color="info"
              max-width="300"
              top
            />
            <v-fade-transition>
              <span v-if="suffix && !hideDuration" class="pr-6 float-right text-lowercase">{{ suffix }}</span>
            </v-fade-transition>
          </span>
        </template>
      </component>
    </v-flex>
    <v-fade-transition>
      <v-flex v-if="!hideDuration" xs4>
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
    </v-fade-transition>
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
    suffix: {
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
    switcher: {
      type: Boolean,
      default: false,
    },
    hideValueOnFalse: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const enabledFieldName = computed(() => `${props.name}.enabled`);

    const hideDuration = computed(() => props.hideValueOnFalse && !props.duration.enabled);

    const component = computed(() => (props.switcher ? 'c-enabled-field' : 'v-checkbox'));

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
      hideDuration,
      component,

      timeUnits,
      min,
      validate,
    };
  },
};
</script>
