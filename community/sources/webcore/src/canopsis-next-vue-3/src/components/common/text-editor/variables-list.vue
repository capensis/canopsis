<template>
  <v-list
    :density="dense ? 'compact' : 'default'"
    class="variables-list py-0"
  >
    <slot
      :items="items"
      name="prepend"
    />

    <template
      v-for="item in items"
      :key="item[itemValue] || item[itemText]"
    >
      <v-list-subheader
        v-if="item.header"
        class="variables-list__header"
      >
        {{ item.header }}
      </v-list-subheader>

      <v-list-group
        v-else-if="item[childrenKey]?.length"
        :value="item[itemText]"
      >
        <template #activator="{ props: activatorProps }">
          <v-list-item
            v-bind="activatorProps"
            :title="item[itemText]"
            @click="clickableParent && emitValue(item)"
          >
            <template
              v-if="item.alias"
              #append
            >
              <c-icon
                class="ml-1"
                color="primary"
                size="small"
              >
                alternate_email
              </c-icon>
            </template>
          </v-list-item>
        </template>

        <variables-list
          :items="item[childrenKey]"
          :children-key="childrenKey"
          :item-value="itemValue"
          :item-text="itemText"
          :show-value="showValue"
          :hide-empty-value="hideEmptyValue"
          :dense="dense"
          :return-object="returnObject"
          :clickable-parent="clickableParent"
          :search="search"
          @update:value="$emit('update:value', $event)"
        />
      </v-list-group>

      <v-list-item
        v-else
        :active="isActive(item)"
        @click="emitValue(item)"
      >
        <template #title>
          <c-layout
            align="center"
            justify="space-between"
          >
            <c-layout
              v-if="item[itemText] || item.alias"
              align="center"
            >
              <span>{{ item[itemText] }}</span>

              <c-icon
                v-if="item.alias"
                class="ml-1"
                color="primary"
                size="small"
              >
                alternate_email
              </c-icon>
            </c-layout>

            <span
              v-if="showValue && (!hideEmptyValue || item[itemValue] !== undefined)"
              class="text-medium-emphasis"
            >
              {{ item[itemValue] }}
            </span>

            <span
              v-if="item.defined"
              class="text-medium-emphasis"
            >
              {{ t('common.defined') }}
            </span>
          </c-layout>
        </template>
      </v-list-item>
    </template>

    <slot
      v-if="!items.length"
      name="no-data"
    />
  </v-list>
</template>

<script>
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    value: {
      type: [Array, String, Number, Object],
      default: () => [],
    },
    search: {
      type: [String, Number],
      default: '',
    },
    items: {
      type: Array,
      default: () => [],
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    itemText: {
      type: String,
      default: 'text',
    },
    childrenKey: {
      type: String,
      default: 'items',
    },
    showValue: {
      type: Boolean,
      default: false,
    },
    hideEmptyValue: {
      type: Boolean,
      default: false,
    },
    dense: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    clickableParent: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:value'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const getValue = item => (props.returnObject ? item : item[props.itemValue]);
    const emitValue = item => emit('update:value', getValue(item));
    const isActive = item => props.value === getValue(item);

    return {
      t,
      emitValue,
      isActive,
    };
  },
};
</script>

<style lang="scss" scoped>
.variables-list__header {
  min-height: 32px;
  pointer-events: none;
  font-size: 1.1rem;
  font-weight: 700;
  line-height: 1.5rem;
  opacity: 1;
  color: rgba(var(--v-theme-on-surface), .78);
}
</style>
