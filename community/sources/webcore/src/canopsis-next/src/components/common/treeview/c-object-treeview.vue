<template>
  <c-treeview :items="resolvedItems" item-key="path" item-title="name">
    <template #prepend="{ item }">
      <div v-if="item.isArray" class="text-caption font-italic text-medium-emphasis">
        {{ `(${t('common.variableTypes.array')})` }}
      </div>
    </template>

    <template #label="{ item, leaf }">
      <div class="c-object-treeview__item">
        <strong v-if="strongTitle">
          {{ item.name }}
        </strong>
        <template v-else>
          {{ item.name }}
        </template>
        <span v-if="leaf" class="pl-1">:
          <c-ellipsis
            :text="formatValue(item.value)"
            :style="valueStyle"
            class="pl-1 d-inline-block text-medium-emphasis text-body-1 pre-wrap c-object-treeview__value"
          />
        </span>
        <span v-else-if="!item.children?.length" class="pl-1">:
          <div class="pl-1 d-inline-block text-medium-emphasis text-body-1 font-italic">
            {{ t('common.emptyObject') }}
          </div>
        </span>
      </div>
    </template>

    <template v-if="$slots.append" #append="slotProps">
      <slot v-bind="slotProps" name="append" />
    </template>
  </c-treeview>
</template>

<script>
import { computed } from 'vue';

import { convertObjectToTreeview } from '@/helpers/treeview';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    items: {
      type: Array,
      default: undefined,
    },
    value: {
      type: Object,
      default: undefined,
    },
    rootName: {
      type: String,
      default: 'root',
    },
    strongTitle: {
      type: Boolean,
      default: false,
    },
    valueMaxWidth: {
      type: [Number, String],
      default: undefined,
    },
    valueFormatter: {
      type: Function,
      default: value => String(value),
    },
  },

  setup(props) {
    const { t } = useI18n();

    const resolvedItems = computed(() =>
      props.items ?? [convertObjectToTreeview(props.value ?? {}, props.rootName)],
    );

    const valueStyle = computed(() => {
      if (!props.valueMaxWidth) {
        return {};
      }

      return {
        maxWidth: Number.isInteger(props.valueMaxWidth)
          ? `${props.valueMaxWidth}px`
          : props.valueMaxWidth,
      };
    });
    const formatValue = value => props.valueFormatter(value);

    return {
      t,
      resolvedItems,
      valueStyle,
      formatValue,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-object-treeview {
  &__item {
    min-width: 0;
  }

  &__value {
    vertical-align: bottom;
  }
}
</style>
