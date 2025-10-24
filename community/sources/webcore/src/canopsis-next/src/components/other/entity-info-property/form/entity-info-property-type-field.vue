<template>
  <v-select
    v-field="value"
    v-validate="rules"
    :items="types"
    :label="label || $t('common.type')"
    :name="name"
    :error-messages="errors.collect(name)"
    clearable
  />
</template>

<script>
import { computed } from 'vue';
import { Validator } from 'vee-validate';

import { ENTITY_INFO_PROPERTY_TYPES, ENTITY_INFO_PROPERTY_TYPE_I18N_KEYS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  inject: {
    $validator: {
      default: new Validator(),
    },
  },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Number,
      required: false,
    },
    label: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'type',
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t } = useI18n();

    const rules = computed(() => ({
      required: props.required,
    }));

    const types = computed(() => Object.values(ENTITY_INFO_PROPERTY_TYPES).map(value => ({
      text: t(ENTITY_INFO_PROPERTY_TYPE_I18N_KEYS[value]),
      value,
    })));

    return {
      rules,
      types,
    };
  },
};
</script>
