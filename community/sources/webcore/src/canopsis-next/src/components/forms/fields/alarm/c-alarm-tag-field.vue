<template lang="pug">
  v-select(
    v-field="value",
    :items="alarmTags",
    :label="label || $tc('common.tag')",
    :loading="alarmTagsPending",
    :disabled="disabled",
    :name="name",
    item-text="value",
    item-value="value",
    hide-details,
    multiple,
    chips,
    dense,
    small-chips,
    clearable
  )
    template(#selection="{ item, index }")
      c-alarm-action-chip(:color="item.color", closable, @close="removeItemFromArray(index)") {{ item.value }}
</template>

<script>
import { entitiesAlarmTagMixin } from '@/mixins/entities/alarm-tag';
import { formArrayMixin } from '@/mixins/form';

export default {
  mixins: [entitiesAlarmTagMixin, formArrayMixin],
  props: {
    value: {
      type: [Array],
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'tag',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  mounted() {
    if (!this.alarmTagsPending) {
      this.fetchAlarmTagsList({ params: { paginate: false } });
    }
  },
};
</script>
