<template lang="pug">
  v-select.c-alarm-tag-field(
    v-field="value",
    :items="alarmTags",
    :label="label || $tc('common.tag')",
    :loading="alarmTagsPending",
    :disabled="disabled",
    :name="name",
    :menu-props="{ contentClass: 'c-alarm-tag-field__list' }",
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
      c-alarm-action-chip.c-alarm-tag-field__tag(
        :color="item.color",
        :title="item.value",
        closable,
        ellipsis,
        @close="removeItemFromArray(index)"
      ) {{ item.value }}
    template(#item="{ item, tile, parent }")
      v-list-tile.c-alarm-tag-field__list-item(v-bind="tile.props", v-on="tile.on")
        v-list-tile-action
          v-checkbox(:input-value="tile.props.value", :color="parent.color")
        v-list-tile-content.c-word-break-all {{ item.value }}
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

<style lang="scss">
$selectIconsWidth: 56px;

.c-alarm-tag-field {
  .v-select__selections {
    width: calc(100% - #{$selectIconsWidth});
  }

  &__tag {
    max-width: 100%;
  }

  &__list {
    max-width: 400px;
  }

  &__list-item .v-list__tile__action {
    height: 40px !important;
    width: 40px !important;
  }

  &__list-item .v-list__tile {
    height: unset !important;
  }
}
</style>
