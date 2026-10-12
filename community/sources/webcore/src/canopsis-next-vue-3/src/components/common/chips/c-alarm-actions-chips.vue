<template>
  <c-layout
    :key="wrapperKey"
    :class="{ 'c-alarm-actions-chips--small': small }"
    class="c-alarm-actions-chips"
    align="center"
    wrap
  >
    <c-chip
      v-for="item in inlineItems"
      :key="item[itemValue]"
      :class="itemClass"
      :color="item.color"
      :small="small"
      :closable="isClosableItem(item)"
      :outlined="outlined"
      :text-color="getItemTextColor(item)"
      @click="selectItem(item)"
      @close="closeItem(item)"
    >
      <slot :item="item" name="item">
        <span>{{ item[itemText] }}</span>
      </slot>
    </c-chip>

    <v-menu
      v-if="dropDownItems.length"
      key="more"
      max-height="400px"
      location="bottom end"
    >
      <template #activator="{ props: activatorProps }">
        <v-btn
          v-bind="activatorProps"
          :class="{ 'c-alarm-actions-chips__more-btn--small': small }"
          :size="small ? 'x-small' : 'small'"
          class="c-alarm-actions-chips__more-btn ma-0"
          color="grey"
          icon
          @click="$emit('activate')"
        >
          <v-icon
            :size="small ? 16 : 20"
            color="white"
          >
            more_horiz
          </v-icon>
        </v-btn>
      </template>

      <v-card>
        <v-card-text>
          <c-layout
            :class="{ 'c-alarm-actions-chips--small': small }"
            class="c-alarm-actions-chips__more"
            wrap
          >
            <c-chip
              v-for="item in dropDownItems"
              :key="item[itemValue]"
              :class="itemClass"
              :color="item.color"
              :closable="isClosableItem(item)"
              :text-color="getItemTextColor(item)"
              :outlined="outlined"
              class="mx-0"
              @click="selectItem(item)"
              @close="closeItem(item)"
            >
              <slot :item="item" name="item">
                <span>{{ item[itemText] }}</span>
              </slot>
            </c-chip>
          </c-layout>
        </v-card-text>
      </v-card>
    </v-menu>
  </c-layout>
</template>

<script>
import { isNil } from 'lodash';
import { computed, ref, watch } from 'vue';

import { getMostReadableTextColor } from '@/helpers/color';
import { uid } from '@/helpers/uid';

export default {
  emits: ['select', 'close', 'activate'],

  props: {
    items: {
      type: Array,
      default: () => [],
    },
    activeItem: {
      type: String,
      required: false,
    },
    activeItems: {
      type: Array,
      default: () => [],
    },
    inlineCount: {
      type: [Number, String],
      default: 2,
    },
    closable: {
      type: Boolean,
      default: false,
    },
    closableActive: {
      type: Boolean,
      default: false,
    },
    small: {
      type: Boolean,
      default: false,
    },
    outlined: {
      type: Boolean,
      default: false,
    },
    itemValue: {
      type: String,
      default: 'text',
    },
    itemText: {
      type: String,
      default: 'text',
    },
    itemClass: {
      type: [String, Object],
      required: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    textColor: {
      type: String,
      required: false,
    },
  },

  setup(props, { emit }) {
    const wrapperKey = ref(uid());

    const sortedItems = computed(() => [...props.items].sort((first, second) => {
        if (first[props.itemValue] === props.activeItem || props.activeItems.includes(first[props.itemValue])) {
          return -1;
        }

        if (second[props.itemValue] === props.activeItem || props.activeItems.includes(second[props.itemValue])) {
          return 0;
        }

        if (first[props.itemText] < second[props.itemText]) {
          return -1;
        }

        if (first[props.itemText] > second[props.itemText]) {
          return 1;
        }

        return 0;
      }),
    );
    const inlineItems = computed(() => sortedItems.value.slice(0, props.inlineCount));
    const dropDownItems = computed(() => sortedItems.value.slice(props.inlineCount));

    watch(inlineItems, () => {
      wrapperKey.value = uid();
    });

    /**
     * Emits `select` with the full item or its value (per `returnObject`), unless the item is already in `activeItems`.
     *
     * @param {Object} item - Chip item; value key is `itemValue`.
     */
    const selectItem = (item) => {
      if (props.activeItems.includes(item[props.itemValue])) {
        return;
      }

      emit('select', props.returnObject ? item : item[props.itemValue]);
    };

    /**
     * Emits `close` with the full item or its value (per `returnObject`).
     *
     * @param {Object} item - Chip item; value key is `itemValue`.
     */
    const closeItem = item => emit('close', props.returnObject ? item : item[props.itemValue]);

    const getItemTextColor = (item) => {
      if (!isNil(props.textColor)) {
        return props.textColor;
      }

      return getMostReadableTextColor(item.color, { level: 'AA', size: 'large' });
    };

    /**
     * Whether the chip should show a close control: when `closable`, or when `closableActive` and the item is active.
     *
     * @param {Object} item - Chip item; compared using `itemValue` and `activeItem` / `activeItems`.
     * @returns {boolean}
     */
    const isClosableItem = (item) => {
      const isActive =
        props.activeItem === item[props.itemValue]
        || props.activeItems.includes(item[props.itemValue]);

      return props.closable || (props.closableActive && isActive);
    };

    return {
      closeItem,
      dropDownItems,
      getItemTextColor,
      inlineItems,
      isClosableItem,
      selectItem,
      wrapperKey,
    };
  },
};
</script>

<style lang="scss">
.c-alarm-actions-chips {
  &,
  &__more {
    column-gap: 8px;
    row-gap: 4px;
  }

  &--small {
    column-gap: 4px;
  }

  & &__more-btn {
    width: 24px;
    height: 24px;
    min-width: 24px;
  }

  & &__more-btn--small {
    min-width: 20px;
  }
}
</style>
