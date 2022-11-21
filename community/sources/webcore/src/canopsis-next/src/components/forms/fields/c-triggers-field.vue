<template lang="pug">
  v-select(
    v-validate="'required'",
    v-field="value",
    :items="availableTriggers",
    :disabled="disabled",
    :label="label || $tc('common.trigger', 2)",
    :error-messages="errors.collect(name)",
    :name="name",
    multiple,
    chips
  )
    template(#item="{ item, tile, parent }")
      v-list-tile(v-bind="tile.props", v-on="tile.on")
        v-list-tile-action
          v-checkbox(:input-value="tile.props.value", :color="parent.color")
        v-list-tile-content {{ item.text }}
        v-list-tile-action(v-if="item.helpText")
          c-help-icon(:text="item.helpText", size="20", top)
</template>

<script>
import { TRIGGERS, PRO_TRIGGERS } from '@/constants';

import { entitiesInfoMixin } from '@/mixins/entities/info';

export default {
  inject: ['$validator'],
  mixins: [entitiesInfoMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'triggers',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    triggers: {
      type: Array,
      default: () => Object.values(TRIGGERS),
    },
  },
  computed: {
    availableTriggers() {
      return Object.values(this.triggers)
        .reduce((acc, type) => {
          if (!PRO_TRIGGERS.includes(type) || this.isProVersion) {
            const { text, helpText } = this.$t(`common.triggers.${type}`);

            acc.push({
              text,
              helpText,
              value: type,
            });
          }

          return acc;
        }, []);
    },
  },
};
</script>
