<template lang="pug">
  v-select(
    v-validate="'required'",
    v-field="value",
    :items="availableTriggers",
    :disabled="disabled",
    :label="label || $t('common.triggers')",
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
          c-help-icon(:text="item.helpText", color="info", size="20", top)
</template>

<script>
import { SCENARIO_TRIGGERS, PRO_SCENARIO_TRIGGERS } from '@/constants';

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
      required: true,
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
  },
  computed: {
    availableTriggers() {
      return Object.values(SCENARIO_TRIGGERS)
        .reduce((acc, type) => {
          if (!PRO_SCENARIO_TRIGGERS.includes(type) || this.isProVersion) {
            const { text, helpText } = this.$t(`common.scenarioTriggers.${type}`);

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
