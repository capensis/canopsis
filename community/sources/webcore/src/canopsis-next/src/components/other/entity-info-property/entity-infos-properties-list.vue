<template>
  <c-advanced-data-table
    :headers="headers"
    :items="entityInfosProperties"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    advanced-pagination
    search
    @update:options="$emit('update:options', $event)"
  >
    <template #toolbar>
      <v-flex xs2>
        <entity-info-property-type-field
          :value="options.type"
          :label="$t('common.filterByType')"
          @input="updateType"
        />
      </v-flex>
    </template>
    <template #mass-actions="{ selected }">
      <c-action-btn
        :tooltip="$t('common.bulkDelete')"
        icon="delete"
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
    </template>
    <template #type="{ item }">
      <span>{{ $t(ENTITY_INFO_PROPERTY_TYPE_I18N_KEYS[item.type]) }}</span>
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="updatable"
          type="edit"
          @click="$emit('edit', item)"
        />
        <c-action-btn
          v-if="duplicable"
          type="duplicate"
          @click="$emit('duplicate', item)"
        />
        <c-action-btn
          v-if="removable"
          type="delete"
          @click="$emit('remove', item)"
        />
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { ENTITY_INFO_PROPERTY_TYPE_I18N_KEYS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import EntityInfoPropertyTypeField from './form/entity-info-property-type-field.vue';

export default {
  components: {
    EntityInfoPropertyTypeField,
  },
  props: {
    entityInfosProperties: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      default: 0,
    },
    options: {
      type: Object,
      default: () => ({}),
    },
    duplicable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    removable: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('entityInfoProperties.infosKey'),
        value: 'name',
      },
      {
        text: t('common.description'),
        value: 'description',
        sortable: false,
      },
      {
        text: t('common.type'),
        value: 'type',
      },
      {
        text: t('common.alias'),
        value: 'alias',
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    /**
     * Updates the type option in the options object and emits the change to the parent component
     *
     * @param {string} type - The new type value to set in the options
     */
    const updateType = type => emit('update:options', { ...props.options, type });

    return {
      ENTITY_INFO_PROPERTY_TYPE_I18N_KEYS,

      headers,

      updateType,
    };
  },
};
</script>
