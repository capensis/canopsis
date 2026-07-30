<template>
  <v-list :density="dense ? 'compact' : 'default'" class="variables-list py-0">
    <slot :items="items" name="prepend" />

    <template v-for="item in items" :key="item[itemValue] || item[itemText]">
      <v-list-subheader v-if="item.header" class="variables-list__header">
        {{ item.header }}
      </v-list-subheader>

      <v-list-item
        v-else
        :active="isActive(item)"
        :disabled="isDisabled(item)"
        @click="selectItem(item)"
        @mouseenter="showSubItems(item, $event)"
      >
        <template #title>
          <c-layout gap="4" align="center" justify="space-between">
            <c-layout v-if="item[itemText] || item.alias" align="center">
              <span>{{ item[itemText] }}</span>

              <c-icon
v-if="item.alias"
class="ml-1"
color="primary"
size="small">
                alternate_email
              </c-icon>
            </c-layout>

            <span
              v-if="showValue && (!hideEmptyValue || item[itemValue] !== undefined)"
              class="text-medium-emphasis"
            >
              {{ item[itemValue] }}
            </span>

            <span v-if="item.defined" class="text-medium-emphasis">
              {{ t('common.defined') }}
            </span>
          </c-layout>
        </template>

        <template v-if="hasChildren(item)" #append>
          <c-icon>arrow_right</c-icon>
        </template>
      </v-list-item>
    </template>

    <slot v-if="!items.length" name="no-data" />

    <v-menu
      v-if="subItemsShown"
      v-model="subItemsShown"
      :target="[subItemsPosition.x, subItemsPosition.y]"
      :close-on-content-click="clickableParent"
      location="bottom start"
      submenu
    >
      <variables-list
        :value="value"
        :items="parentItemChildren"
        :children-key="childrenKey"
        :item-value="itemValue"
        :item-text="itemText"
        :show-value="showValue"
        :hide-empty-value="hideEmptyValue"
        :dense="dense"
        :return-object="returnObject"
        :clickable-parent="clickableParent"
        :search="search"
        @update:value="selectSubItem"
      />
    </v-menu>
  </v-list>
</template>

<script>
import { isArray, isObject } from 'lodash';
import { computed, ref } from 'vue';

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
    const subItemsShown = ref(false);
    const parentItem = ref();
    const subItemsPosition = ref({ x: 0, y: 0 });

    const parentItemChildren = computed(() => getChildren(parentItem.value) ?? []);

    const getChildren = item =>
      item?.[props.childrenKey] ??
      (isArray(item?.[props.itemValue]) ? item[props.itemValue] : undefined);

    const hasChildren = item => !!getChildren(item);
    const hasAvailableChildren = item => !!getChildren(item)?.length;
    const isDisabled = item => item.disabled || item.loading || getChildren(item)?.length === 0;

    const prepareItem = (item) => {
      if (props.returnObject) {
        return item;
      }

      return isObject(item) ? item?.[props.itemValue] : item;
    };

    const getValue = value => (props.returnObject ? value?.[props.itemValue] : value);
    const isActive = item => props.value === item[props.itemValue];

    const selectItem = (item) => {
      if (hasAvailableChildren(item)) {
        subItemsShown.value = true;
        return;
      }

      emit('update:value', prepareItem(item));
    };

    const selectSubItem = (item) => {
      const hasParentValue = !!parentItem.value?.[props.itemValue];
      const parentValue = String(parentItem.value?.[props.itemValue]).replace(/\s*}}$/, '');
      const value = getValue(item);
      const newValue =
        String(value).startsWith(parentValue) || !hasParentValue
          ? value
          : `${parentValue}.${value}`;

      if (props.returnObject) {
        selectItem({
          ...item,
          [props.itemValue]: newValue,
        });
      } else {
        selectItem(newValue);
      }

      subItemsShown.value = false;
    };

    const showSubItems = (item, event) => {
      if (isDisabled(item)) {
        return;
      }

      if (hasChildren(item)) {
        const { left, top, width } = event.currentTarget.getBoundingClientRect();

        subItemsPosition.value = {
          x: left + width,
          y: top,
        };
        parentItem.value = item;
        subItemsShown.value = true;
      } else {
        parentItem.value = undefined;
        subItemsShown.value = false;
      }
    };

    return {
      t,
      subItemsShown,
      subItemsPosition,
      parentItemChildren,
      hasChildren,
      isDisabled,
      selectItem,
      selectSubItem,
      showSubItems,
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
  color: rgba(var(--v-theme-on-surface), 0.78);
}
</style>
