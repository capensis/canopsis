<template lang="pug">
  v-card
    v-card-text
      v-layout(column)
        v-layout(row, align-center)
          c-name-field.mr-2(
            v-field="form.label",
            :label="$t('common.label')",
            :name="labelFieldName",
            required
          )
          c-action-btn(type="delete", @click="remove")
        v-layout(row, align-center)
          v-flex.pr-2(xs8)
            v-text-field(
              v-field="form.category",
              :label="$t('common.category')"
            )
          v-flex(xs4)
            c-icon-field(
              v-field="form.icon_name",
              :label="$t('common.icon')",
              :name="iconFieldName",
              required
            )
        c-enabled-field(
          v-if="isAlarmType",
          v-field="form.single",
          :label="$t('linkRule.single')"
        )
        c-payload-text-field(
          v-field="form.url",
          :label="$t('common.url')",
          :variables="urlVariables",
          :name="form.key",
          required
        )
</template>

<script>
import { ENTITY_PAYLOADS_VARIABLES, LINK_RULE_TYPES } from '@/constants';

import { payloadVariablesMixin } from '@/mixins/payload/variables';

export default {
  inject: ['$validator'],
  mixins: [payloadVariablesMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    type: {
      type: String,
      default: LINK_RULE_TYPES.alarm,
    },
    name: {
      type: String,
      default: 'link',
    },
  },
  computed: {
    isAlarmType() {
      return this.type === LINK_RULE_TYPES.alarm;
    },

    labelFieldName() {
      return `${this.name}.label`;
    },

    iconFieldName() {
      return `${this.name}.icon`;
    },

    alarmUrlVariables() {
      return [
        ...this.alarmPayloadRangeVariables,
        ...this.externalDataAlarmPayloadVariables,
      ];
    },

    entityUrlVariables() {
      return [
        {
          value: ENTITY_PAYLOADS_VARIABLES.entity,
          enumerable: true,
          variables: [
            {
              value: ENTITY_PAYLOADS_VARIABLES.infosValue,
              text: this.$t('common.infos'),
            },
          ],
        },

        ...this.externalDataEntityPayloadVariables,
      ];
    },

    urlVariables() {
      return this.isAlarmType
        ? this.alarmUrlVariables
        : this.entityUrlVariables;
    },
  },
  methods: {
    remove() {
      this.$emit('remove');
    },
  },
};
</script>
