<template>
  <v-card>
    <v-card-text>
      <v-layout column>
        <v-layout align-center>
          <c-payload-text-field
            v-field="form.label"
            :label="$t('common.label')"
            :variables="variables"
            :name="labelFieldName"
            class="mr-2"
            required
          />
          <c-action-btn
            type="delete"
            @click="remove"
          />
        </v-layout>
        <v-layout align-center>
          <v-flex
            class="pr-2"
            xs8
          >
            <v-text-field
              v-field="form.category"
              :label="$t('common.category')"
            />
          </v-flex>
          <v-flex xs4>
            <c-icon-field
              v-field="form.icon_name"
              :label="$tc('common.icon', 1)"
              :name="iconFieldName"
              required
            />
          </v-flex>
        </v-layout>
        <template v-if="isAlarmType">
          <c-enabled-field
            v-field="form.single"
            :label="$t('linkRule.single')"
            hide-details
          />
          <c-enabled-field
            v-field="form.hide_in_menu"
            :label="$t('linkRule.hideInMenu')"
          />
        </template>
        <c-payload-text-field
          v-field="form.url"
          :label="$t('common.url')"
          :variables="variables"
          :name="form.key"
          required
        />
        <v-radio-group
          v-field="form.action"
          :label="$t('linkRule.actionType')"
        >
          <v-radio
            :value="$constants.LINK_RULE_ACTIONS.open"
            :label="$t('linkRule.actionTypes.open')"
            color="primary"
          />
          <v-radio
            :value="$constants.LINK_RULE_ACTIONS.copy"
            :label="$t('linkRule.actionTypes.copy')"
            color="primary"
          />
        </v-radio-group>
      </v-layout>
    </v-card-text>
  </v-card>
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

    alarmVariables() {
      return [
        ...this.alarmPayloadRangeVariables,
        ...this.externalDataAlarmPayloadVariables,
        ...this.userPayloadVariables,
      ];
    },

    entityVariables() {
      return [
        {
          value: ENTITY_PAYLOADS_VARIABLES.entities,
          enumerable: true,
          variables: [
            {
              value: ENTITY_PAYLOADS_VARIABLES.infosValue,
              text: this.$t('common.infos'),
            },
          ],
        },

        ...this.externalDataEntityPayloadVariables,
        ...this.userPayloadVariables,
      ];
    },

    variables() {
      return this.isAlarmType
        ? this.alarmVariables
        : this.entityVariables;
    },
  },
  methods: {
    remove() {
      this.$emit('remove');
    },
  },
};
</script>
