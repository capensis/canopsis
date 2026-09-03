<template>
  <v-menu
    v-model="menu"
    :min-width="listMinWidth"
    :transition="false"
    content-class="c-advanced-search__history-menu"
    location="bottom"
    offset="1"
    @update:model-value="handleMenuToggle"
  >
    <template #activator="{ props }">
      <c-action-btn
        :tooltip="$t('common.search')"
        :disabled="!searches.length"
        icon="history"
        v-bind="props"
      />
    </template>
    <div ref="scrollContainerElement" class="c-advanced-search__history-list">
      <v-list class="pa-0" density="compact">
        <div
          v-for="search in items"
          :key="search._id"
          :ref="setItemRef(search._id)"
          class="c-advanced-search__history-item"
        >
          <v-list-item v-if="isItemVisible(search._id)" @click="select(search)">
            <div class="c-advanced-search__history-item-content">
              <advanced-search-rules :rules="search.rules" :attributes="attributes" disabled />
              <advanced-search-history-item-btns
                :id="search._id"
                :pinned="search.pinned"
                @remove="remove"
                @toggle-pin="togglePin"
              />
            </div>
          </v-list-item>
          <div
            v-else
            :style="getItemPlaceholderStyle(search._id)"
            class="c-advanced-search__history-item__placeholder"
          />
        </div>
        <div ref="appendItemElement" class="c-advanced-search__history-list__append-item" />
      </v-list>
    </div>
  </v-menu>
</template>

<script>
import { nextTick, ref, toRef } from 'vue';

import { useAdvancedSearchHistory } from '../hooks/advanced-search-history';

import AdvancedSearchHistoryItemBtns from './advanced-search-history-item-btns.vue';
import AdvancedSearchRules from './advanced-search-rules.vue';

export default {
  components: { AdvancedSearchHistoryItemBtns, AdvancedSearchRules },

  props: {
    searches: {
      type: Array,
      default: () => [],
    },
    attributes: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['select', 'remove', 'toggle-pin'],

  setup(props, { emit }) {
    const menu = ref(false);
    const {
      items,
      scrollContainerElement,
      appendItemElement,
      setItemRef,
      isItemVisible,
      getItemPlaceholderStyle,
      listMinWidth,
      updateListWidth,
    } = useAdvancedSearchHistory(toRef(props, 'searches'));

    /**
     * Recalculates list width after the menu content is mounted.
     *
     * @param {boolean} isOpen
     */
    const handleMenuToggle = (isOpen) => {
      if (!isOpen) {
        return;
      }

      nextTick(() => {
        nextTick(() => updateListWidth());
      });
    };

    const select = search => emit('select', search);
    const remove = id => emit('remove', id);
    const togglePin = id => emit('toggle-pin', id);

    return {
      menu,
      items,
      scrollContainerElement,
      appendItemElement,
      setItemRef,
      isItemVisible,
      getItemPlaceholderStyle,
      listMinWidth,
      handleMenuToggle,
      select,
      remove,
      togglePin,
    };
  },
};
</script>

<style lang="scss">
.c-advanced-search__history-menu {
  max-height: none !important;
  overflow-y: hidden !important;
  width: max-content;
  min-width: max-content;
  background: rgb(var(--v-theme-surface));
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  border-radius: 4px;
  box-shadow:
    0 5px 5px -3px rgba(0, 0, 0, 0.2),
    0 8px 10px 1px rgba(0, 0, 0, 0.14),
    0 3px 14px 2px rgba(0, 0, 0, 0.12);

  @media (min-width: 1264px) {
    max-width: 95vw;
  }
}

.c-advanced-search__history-list {
  max-height: 95vh;
  overflow-y: auto;
  width: max-content;
  min-width: 100%;
}

.c-advanced-search__history-item {
  width: max-content;
  min-width: 100%;
}

.c-advanced-search__history-item-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  width: max-content;
  min-width: 100%;
}

.c-advanced-search__history-item__placeholder {
  width: 100%;
}

.c-advanced-search__history-list__append-item {
  height: 1px;
}
</style>
