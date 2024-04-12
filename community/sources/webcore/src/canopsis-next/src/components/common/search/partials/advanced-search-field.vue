<template>
  <advanced-search-field-wrapper :is-menu-active="isMenuActive">
    <template #activator>
      <div ref="activatorElement" class="v-input__slot">
        <div class="v-select__slot">
          <div class="v-select__selections">
            <advanced-search-field-chip
              v-for="(item, index) in value"
              :key="item.key"
              :item="item"
              :active="index === activeIndex"
              :last="index === value.length - 1"
              @input="updateItemInArray(index, $event)"
              @click:item="clickChip(index)"
              @remove:item="removeChip(index)"
              @keydown:input="keydownInputChip"
            />
            <advanced-search-field-input
              v-model="internalSearch"
              ref="internalSearchElement"
              :is-value-type="isItemTypeValue"
              :placeholder="isMenuActive || value.length ? '' : $t('advancedSearch.title')"
              @click="clickInput"
              @apply="applyInput"
              @keydown:backspace="keydownBackspaceInput"
              @keydown:navigate="keydownNavigateInput"
            />
          </div>
        </div>
      </div>
    </template>
    <template #items>
      <advanced-search-field-list
        ref="listElement"
        :value="isMenuActive"
        :activator="$refs.activatorElement"
        :not-switcher-value="notValues[activeIndex]"
        :active-type="activeType"
        :active-value="activeItem?.value"
        :items="filteredItems"
        :list-message="listMessage"
        @input="changeMenu"
        @select:item="selectItem"
        @change:not="setNotValue"
      />
    </template>
  </advanced-search-field-wrapper>
</template>

<script>
import {
  computed,
  ref,
  toRef,
  nextTick,
  watch,
  set,
} from 'vue';

import { ADVANCED_SEARCH_CONDITIONS, ADVANCED_SEARCH_ITEM_TYPES, KEY_CODES } from '@/constants';

import { uid } from '@/helpers/uid';

import { useI18n } from '@/hooks/i18n';
import { useArrayModelField } from '@/hooks/form/array-model-field';

import {
  useAdvancedSearchItemType,
  useAdvancedSearchInternalSearch,
  useAdvancedSearchActiveIndex,
  useAdvancedSearchMenu,
  useAdvancedSearchItems,
  useAdvancedSearchNotValues,
} from '../hooks/advanced-search';

import AdvancedSearchFieldWrapper from './advanced-search-field-wrapper.vue';
import AdvancedSearchFieldList from './advanced-search-field-list.vue';
import AdvancedSearchFieldChip from './advanced-search-field-chip.vue';
import AdvancedSearchFieldInput from './advanced-search-field-input.vue';

export default {
  components: {
    AdvancedSearchFieldWrapper,
    AdvancedSearchFieldList,
    AdvancedSearchFieldChip,
    AdvancedSearchFieldInput,
  },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
    fields: {
      type: Array,
      default: () => [],
    },
    conditions: {
      type: Array,
      default: () => Object.values(ADVANCED_SEARCH_CONDITIONS),
    },
    initialInternalSearch: {
      type: String,
      default: '',
    },
  },
  setup(props, { emit }) {
    const {
      addItemIntoArray,
      updateItemInArray,
      removeItemFromArray,
    } = useArrayModelField(props, emit);

    const {
      internalSearch,
      internalSearchElement,
      internalSearchElementFocus,
      clearInternalSearch,
    } = useAdvancedSearchInternalSearch({ initialInternalSearch: toRef(props, 'initialInternalSearch') });

    const {
      activeIndex,
      activeItem,
      activeType,
      isActiveIndexLast,
      isActiveIndexNew,
      goToActiveIndex,
      goToNextActiveIndex,
      goToPrevActiveIndex,
      clearActiveIndex,
    } = useAdvancedSearchActiveIndex({
      value: toRef(props, 'value'),
      onChange: clearInternalSearch,
    });

    const onMenuChange = (value) => {
      if (!value) {
        clearActiveIndex();
        clearInternalSearch();
      }
    };

    const {
      isMenuActive,
      changeMenu,
      focusMenu,
      blurMenu,
    } = useAdvancedSearchMenu({ onChange: onMenuChange });

    const {
      notValues,
      activeIndexNotValue,
      setNotValue,
    } = useAdvancedSearchNotValues({
      value: toRef(props, 'value'),
      activeIndex,
      updateItemInArray,
    });

    const selectItem = (item) => {
      if (isActiveIndexNew.value) {
        addItemIntoArray({ ...item, not: activeIndexNotValue.value, key: uid() });
        goToNextActiveIndex();
        internalSearchElementFocus();
      } else {
        updateItemInArray(activeIndex.value, { ...item, not: activeIndexNotValue.value });
        blurMenu();
      }
    };

    const clickChip = (index = props.value.length) => {
      goToActiveIndex(index);
      focusMenu();
    };

    const removeChip = (index = props.value.length - 1) => {
      removeItemFromArray(index);
      goToActiveIndex(index);
    };

    const keydownInputChip = (event) => {
      if (event.keyCode === KEY_CODES.backspace && !event.target.value && isActiveIndexLast.value) {
        removeItemFromArray(activeIndex.value);
        internalSearchElementFocus();
      } else if (event.keyCode === KEY_CODES.enter) {
        if (isActiveIndexLast.value) {
          internalSearchElementFocus();
        }

        goToNextActiveIndex();
      }
    };

    const listElement = ref();

    const { isItemTypeValue } = useAdvancedSearchItemType({ type: activeType });

    const clickInput = () => clickChip();
    const applyInput = () => {
      addItemIntoArray({
        key: uid(),
        value: internalSearch.value,
        type: ADVANCED_SEARCH_ITEM_TYPES.value,
        text: internalSearch.value,
      });
      goToNextActiveIndex();
    };
    const keydownBackspaceInput = () => {
      removeItemFromArray(props.value.length - 1);
      goToPrevActiveIndex();
    };
    const keydownNavigateInput = (event) => {
      const { menuElement } = listElement.value?.$refs ?? {};

      if (!menuElement) {
        return;
      }

      nextTick(() => {
        menuElement.changeListIndex(event);
      });
    };

    const { filteredItems } = useAdvancedSearchItems({
      fields: toRef(props, 'fields'),
      conditions: toRef(props, 'conditions'),
      activeType,
      internalSearch,
    });

    const { t } = useI18n();
    const listMessage = computed(() => {
      const valueToCheck = isItemTypeValue.value && !isActiveIndexNew.value
        ? activeItem.value?.value
        : internalSearch.value;

      return t(`advancedSearch.${!valueToCheck ? 'valueTypeListEmptyMessage' : 'valueTypeListMessage'}`);
    });

    watch(() => props.value, (newValue) => {
      const keys = Object.keys(notValues.value);
      const maxNotValueIndex = Math.max(...keys);

      if (maxNotValueIndex < newValue.length) {
        return;
      }

      keys.forEach((key) => {
        if (key < newValue.length) {
          return;
        }

        set(notValues.value, key, false);
      });

      clearInternalSearch();
    });

    return {
      updateItemInArray,

      internalSearch,
      internalSearchElement,

      activeIndex,
      activeItem,
      activeType,

      isMenuActive,
      changeMenu,

      notValues,
      setNotValue,

      selectItem,

      clickChip,
      removeChip,
      keydownInputChip,

      listElement,
      isItemTypeValue,
      applyInput,
      clickInput,
      keydownBackspaceInput,
      keydownNavigateInput,

      filteredItems,

      listMessage,
    };
  },
};
</script>

<style lang="scss" scoped>
.v-select__selections {
  height: 32px;
}
</style>
