<template>
  <v-card>
    <v-card-text>
      <v-layout column>
        <v-layout align-center>
          <c-payload-text-field
            v-field="form.label"
            :label="$t('common.label')"
            :variables="templateVars.label"
            :name="labelFieldName"
            class="mr-2"
            required
          />
          <c-action-btn
            type="delete"
            @click="remove"
          />
        </v-layout>
        <c-form-block>
          <c-form-block-row :label="$t('common.category')">
            <v-text-field
              v-field="form.category"
              :label="$t('common.category')"
            />
          </c-form-block-row>

          <c-form-block-row :label="$tc('common.icon', 1)">
            <c-icon-field
              v-field="form.icon_name"
              :label="$tc('common.icon', 1)"
              :name="iconFieldName"
              required
            />
          </c-form-block-row>

          <c-form-block-row
            v-if="isAlarmType"
            :label="$t('linkRule.alarmListParameters')"
          >
            <c-enabled-field
              v-field="form.single"
              :label="$t('linkRule.single')"
              hide-details
            />
            <c-enabled-field
              v-field="form.hide_in_menu"
              :label="$t('linkRule.hideInMenu')"
            />
          </c-form-block-row>

          <c-form-block-row :label="$t('common.url')">
            <c-payload-text-field
              v-field="form.url"
              :label="$t('common.url')"
              :variables="templateVars.url"
              :name="form.key"
              required
            />
          </c-form-block-row>

          <c-form-block-row :label="$t('common.actionType')">
            <v-radio-group
              v-field="form.action"
              :label="$t('common.actionType')"
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
          </c-form-block-row>
        </c-form-block>
      </v-layout>
    </v-card-text>
  </v-card>
</template>

<script>
import { computed } from 'vue';

import { LINK_RULE_TYPES } from '@/constants';

export default {
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
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const isAlarmType = computed(() => props.type === LINK_RULE_TYPES.alarm);
    const labelFieldName = computed(() => `${props.name}.label`);
    const iconFieldName = computed(() => `${props.name}.icon`);

    const remove = () => emit('remove');

    return {
      isAlarmType,
      labelFieldName,
      iconFieldName,

      remove,
    };
  },
};
</script>
