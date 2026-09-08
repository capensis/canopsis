<template>
  <v-card>
    <v-card-text>
      <c-layout column>
        <c-layout align="center">
          <c-payload-text-field
            :value="form.label"
            :label="$t('common.label')"
            :variables="templateVars.label"
            :name="labelFieldName"
            class="mr-2"
            required
            @update:value="updateField('label', $event)"
          />
          <c-action-btn type="delete" @click="remove" />
        </c-layout>
        <c-form-block>
          <c-form-block-row :label="$t('common.category')">
            <c-text-field
              :model-value="form.category"
              :label="$t('common.category')"
              @update:model-value="updateField('category', $event)"
            />
          </c-form-block-row>

          <c-form-block-row :label="$t('common.icon', 1)">
            <c-icon-field
              :model-value="form.icon_name"
              :label="$t('common.icon', 1)"
              :name="iconFieldName"
              required
              @update:model-value="updateField('icon_name', $event)"
            />
          </c-form-block-row>

          <c-form-block-row v-if="isAlarmType" :label="$t('linkRule.alarmListParameters')">
            <c-enabled-field
              :model-value="form.single"
              :label="$t('linkRule.single')"
              hide-details
              @update:model-value="updateField('single', $event)"
            />
            <c-enabled-field
              :model-value="form.hide_in_menu"
              :label="$t('linkRule.hideInMenu')"
              @update:model-value="updateField('hide_in_menu', $event)"
            />
          </c-form-block-row>

          <c-form-block-row :label="$t('common.url')">
            <c-payload-text-field
              :value="form.url"
              :label="$t('common.url')"
              :variables="templateVars.url"
              :name="form.key"
              required
              @update:value="updateField('url', $event)"
            />
          </c-form-block-row>

          <c-form-block-row :label="$t('common.actionType')">
            <c-radio-group
              :model-value="form.action"
              :label="$t('common.actionType')"
              @update:model-value="updateField('action', $event)"
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
            </c-radio-group>
          </c-form-block-row>
        </c-form-block>
      </c-layout>
    </v-card-text>
  </v-card>
</template>

<script>
import { computed } from 'vue';

import { LINK_RULE_TYPES } from '@/constants';

export default {
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

  emits: ['update:form', 'remove', 'change'],

  setup(props, { emit }) {
    const isAlarmType = computed(() => props.type === LINK_RULE_TYPES.alarm);
    const labelFieldName = computed(() => `${props.name}.label`);
    const iconFieldName = computed(() => `${props.name}.icon`);

    const updateField = (field, value) => {
      emit('update:form', {
        ...props.form,
        [field]: value,
      });
      emit('change');
    };
    const remove = () => emit('remove');

    return {
      isAlarmType,
      labelFieldName,
      iconFieldName,

      updateField,
      remove,
    };
  },
};
</script>
