<template>
  <v-layout column>
    <c-entity-info-property-key-field
      v-field="form.infos_key"
      :label="$t('entityInfoProperties.infosKey')"
      name="infos_key"
      required
    />
    <c-description-field
      v-field="form.description"
      :label="$t('common.description')"
      :max-length="255"
      name="description"
    />
    <c-name-field
      v-field="form.alias"
      :label="$t('common.alias')"
      name="alias"
    />
    <v-select
      v-field="form.type"
      v-validate="'required'"
      :items="typeOptions"
      :label="$t('common.type')"
      :error-messages="errors.collect('type')"
      name="type"
      required
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import { ENTITY_INFO_PROPERTY_TYPES, ENTITY_INFO_PROPERTY_TYPE_I18N_KEYS } from '@/constants/entity-info-properties';

export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isNew: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();

    const typeOptions = computed(() => Object.values(ENTITY_INFO_PROPERTY_TYPES).map(value => ({
      text: t(ENTITY_INFO_PROPERTY_TYPE_I18N_KEYS[value]),
      value,
    })));

    return {
      typeOptions,
    };
  },
};
</script>
