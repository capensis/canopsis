<template>
  <widget-settings-item :title="$t('settings.fastPbehavior')">
    <c-card-iterator-form
      v-field="form"
      :handle="`.${dragHandleClass}`"
      @add="add"
    >
      <template #header>
        <span class="text-body-2">{{ $t('settings.fastPbehavior') }}</span>
        <span>{{ $t('settings.fastPbehaviorDescription') }}</span>
      </template>
      <template #item="{ index }">
        <c-card-iterator-item
          :drag-handle-class="dragHandleClass"
          small
          @remove="remove(index)"
        >
          <template #header>
            <v-layout class="gap-2" column>
              <v-text-field
                v-field="form[index].name_prefix"
                v-validate="'required'"
                :label="$t('common.namePrefix')"
                :error-messages="errors.collect(getName(form[index], 'name_prefix'))"
                :name="getName(form[index], 'name_prefix')"
              />
              <c-pbehavior-type-field
                v-field="form[index].type"
                :label="$t('common.type')"
                :types="types"
                :name="getName(form[index], 'type')"
                independent
                clearable
              />
              <c-pbehavior-reason-field
                v-field="form[index].reason"
                :label="$t('common.reason')"
                :name="getName(form[index], 'reason')"
                clearable
              />
            </v-layout>
          </template>
        </c-card-iterator-item>
      </template>
    </c-card-iterator-form>
  </widget-settings-item>
</template>
<script>
import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { alarmListFastPbehaviorParametersItemToForm } from '@/helpers/entities/widget/forms/alarm';

import { useArrayModelField } from '@/hooks/form/array-model-field';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  inject: ['$validator'],
  components: { WidgetSettingsItem },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const types = [PBEHAVIOR_TYPE_TYPES.pause];
    const dragHandleClass = 'fast-pbehavior-drag-handle';

    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    /**
     * Adds a new fast pbehavior parameters item with default values to the form.
     */
    const add = () => addItemIntoArray(alarmListFastPbehaviorParametersItemToForm());

    /**
     * Builds a unique field name for vee-validate by combining item key with suffix.
     *
     * @param {Object} item - Form item with key property
     * @param {string} suffix - Suffix to append (e.g. 'name_prefix', 'type', 'reason')
     * @returns {string} Unique field name
     */
    const getName = (item, suffix) => `${item.key}${suffix}`;

    return {
      types,
      dragHandleClass,

      add,
      remove: removeItemFromArray,

      getName,
    };
  },
};
</script>
