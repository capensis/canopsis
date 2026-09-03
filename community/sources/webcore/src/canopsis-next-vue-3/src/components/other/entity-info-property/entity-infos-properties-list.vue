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
      <div class="entity-infos-properties-list__type-filter">
        <entity-info-property-type-field
          :value="options.type"
          :label="$t('common.filterByType')"
          @update:value="updateType"
        />
      </div>
    </template>

    <template #mass-actions="{ selected, clearSelected }">
      <c-action-btn
        :tooltip="$t('common.bulkDelete')"
        icon="delete"
        type="delete"
        small
        @click="$emit('remove-selected', selected, clearSelected)"
      />
    </template>

    <template #type="{ item }">
      <span>{{ $t(ENTITY_INFO_PROPERTY_TYPE_I18N_KEYS[item.type]) }}</span>
    </template>

    <template #actions="{ item }">
      <c-layout>
        <c-action-btn v-if="updatable" type="edit" @click="$emit('edit', item)" />
        <c-action-btn v-if="duplicable" type="duplicate" @click="$emit('duplicate', item)" />
        <c-action-btn v-if="removable" type="delete" @click="$emit('remove', item)" />
      </c-layout>
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

  emits: ['edit', 'duplicate', 'remove', 'remove-selected', 'update:options'],

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
     * @param {number} type - The new type value to set in the options
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

<style lang="scss">
.entity-infos-properties-list__type-filter {
  display: flex;
  align-items: center;
  width: 180px;
}
</style>
