<template>
  <widget-settings-item :title="t('common.pause')">
    <c-radio-group
      :model-value="modelValue.mode"
      class="pause-form__mode mt-0"
      color="primary"
      hide-details
      mandatory
      @update:model-value="updateMode"
    >
      <v-radio
        :value="ALARM_LIST_PAUSE_METHODS.pause"
        color="primary"
      >
        <template #label>
          <span>
            {{ t('settings.pause.pauseAction') }}
            <span class="font-italic text-medium-emphasis">
              ({{ t('settings.pause.recommended') }})
            </span>
          </span>
        </template>
      </v-radio>

      <v-row
        v-if="isPauseAction"
        class="pause-form__parameters pb-4 ma-0"
      >
        <v-col class="pa-0 pb-3" cols="12">
          <div>
            <div class="text-body-2">
              {{ t('settings.pause.types') }}
            </div>
            <div class="text-caption text-medium-emphasis">
              {{ t('settings.pause.typesDescription') }}
            </div>
          </div>
        </v-col>

        <v-col class="pa-0" cols="12">
          <c-card-iterator-form
            :value="modelValue.types"
            :handle="`.${pauseDragHandleClass}`"
            :name="`${name}.types`"
            :add-button-label="t('common.add')"
            add-button-variant="flat"
            item-key="key"
            @update:value="updatePauseTypes"
            @add="addPauseType"
          >
            <template #item="{ index }">
              <c-card-iterator-item
                :drag-handle-class="pauseDragHandleClass"
                small
                @remove="removePauseType(index)"
              >
                <template #header>
                  <v-row class="ma-0">
                    <v-col class="pa-0 pb-2" cols="12">
                      <c-name-field
                        :model-value="modelValue.types[index].name"
                        :max-length="FORM_MAX_STRING_LENGTH"
                        :name="`${name}.types.${modelValue.types[index].key}.name`"
                        required
                        @update:model-value="updatePauseTypeField(index, 'name', $event)"
                      />
                    </v-col>
                    <v-col class="pa-0 pb-2" cols="12">
                      <c-description-field
                        :model-value="modelValue.types[index].comment"
                        :label="t('common.comment', 1)"
                        :max-length="FORM_MAX_STRING_LENGTH"
                        :name="`${name}.types.${modelValue.types[index].key}.comment`"
                        @update:model-value="updatePauseTypeField(index, 'comment', $event)"
                      />
                    </v-col>
                    <v-col class="pa-0" cols="12">
                      <c-text-field
                        :model-value="modelValue.types[index].type"
                        :label="t('common.type')"
                        :max-length="FORM_MAX_STRING_LENGTH"
                        :name="`${name}.types.${modelValue.types[index].key}.type`"
                        @update:model-value="updatePauseTypeField(index, 'type', $event)"
                      />
                    </v-col>
                  </v-row>
                </template>
              </c-card-iterator-item>
            </template>
          </c-card-iterator-form>
        </v-col>
      </v-row>

      <v-radio
        :value="ALARM_LIST_PAUSE_METHODS.pbehavior"
        :label="t('settings.pause.pbehavior')"
        color="primary"
      />

      <v-row
        v-if="isPbehaviorAction"
        class="pause-form__parameters ma-0"
      >
        <v-col class="pa-0 pb-3" cols="12">
          <div>
            <div class="text-body-2">
              {{ t('settings.pause.pbehaviorTypes') }}
            </div>
            <div class="text-caption text-medium-emphasis">
              {{ t('settings.pause.typesDescription') }}
            </div>
          </div>
        </v-col>

        <v-col class="pa-0" cols="12">
          <c-card-iterator-form
            :value="fastPbehaviors"
            :handle="`.${pbehaviorDragHandleClass}`"
            :add-button-label="t('common.add')"
            add-button-variant="flat"
            name="fast_pbehaviors"
            item-key="key"
            @update:value="$emit('update:fastPbehaviors', $event)"
            @add="addPbehavior"
          >
            <template #item="{ index }">
              <c-card-iterator-item
                :drag-handle-class="pbehaviorDragHandleClass"
                small
                @remove="removePbehavior(index)"
              >
                <template #header>
                  <v-row class="ma-0">
                    <v-col class="pa-0 pb-2" cols="12">
                      <c-name-field
                        :model-value="fastPbehaviors[index].name_prefix"
                        :label="t('common.name')"
                        :max-length="FORM_MAX_STRING_LENGTH"
                        :name="`fast_pbehaviors.${fastPbehaviors[index].key}.name_prefix`"
                        required
                        @update:model-value="updatePbehaviorField(index, 'name_prefix', $event)"
                      />
                    </v-col>
                    <v-col class="pa-0 pb-2" cols="12">
                      <c-pbehavior-type-field
                        :model-value="fastPbehaviors[index].type"
                        :label="t('common.type')"
                        :types="pbehaviorTypes"
                        :name="`fast_pbehaviors.${fastPbehaviors[index].key}.type`"
                        independent
                        clearable
                        required
                        @update:model-value="updatePbehaviorField(index, 'type', $event)"
                      />
                    </v-col>
                    <v-col class="pa-0" cols="12">
                      <c-pbehavior-reason-field
                        :model-value="fastPbehaviors[index].reason"
                        :label="t('common.reason')"
                        :name="`fast_pbehaviors.${fastPbehaviors[index].key}.reason`"
                        clearable
                        required
                        @update:model-value="updatePbehaviorField(index, 'reason', $event)"
                      />
                    </v-col>
                  </v-row>
                </template>
              </c-card-iterator-item>
            </template>
          </c-card-iterator-form>
        </v-col>
      </v-row>
    </c-radio-group>
  </widget-settings-item>
</template>

<script>
import { computed } from 'vue';

import {
  ALARM_LIST_PAUSE_METHODS,
  FORM_MAX_STRING_LENGTH,
  PBEHAVIOR_TYPE_TYPES,
} from '@/constants';

import {
  alarmListFastPbehaviorParametersItemToForm,
  alarmListPauseTypeToForm,
} from '@/helpers/entities/widget/forms/alarm';

import { useI18n } from '@/hooks/i18n';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Object,
      default: () => ({
        mode: ALARM_LIST_PAUSE_METHODS.pause,
        types: [],
      }),
    },
    fastPbehaviors: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'pause',
    },
  },

  emits: ['update:modelValue', 'update:fastPbehaviors'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const pauseDragHandleClass = 'pause-type-drag-handle';
    const pbehaviorDragHandleClass = 'pause-pbehavior-drag-handle';
    const pbehaviorTypes = [PBEHAVIOR_TYPE_TYPES.pause];

    const isPauseAction = computed(
      () => props.modelValue.mode === ALARM_LIST_PAUSE_METHODS.pause,
    );
    const isPbehaviorAction = computed(
      () => props.modelValue.mode === ALARM_LIST_PAUSE_METHODS.pbehavior,
    );

    /**
     * Set the selected pause mode.
     *
     * @param {string} mode
     */
    const updateMode = mode => emit('update:modelValue', { ...props.modelValue, mode });

    /**
     * Replace the pause types list.
     *
     * @param {Object[]} types
     */
    const updatePauseTypes = types => emit('update:modelValue', { ...props.modelValue, types });

    /**
     * Append an empty pause type.
     */
    const addPauseType = () => updatePauseTypes([
      ...(props.modelValue.types ?? []),
      alarmListPauseTypeToForm(),
    ]);

    /**
     * Update one field of a pause type.
     *
     * @param {number} index
     * @param {string} field
     * @param {string} value
     */
    const updatePauseTypeField = (index, field, value) => {
      const types = [...(props.modelValue.types ?? [])];
      types[index] = { ...types[index], [field]: value };
      updatePauseTypes(types);
    };

    /**
     * Remove a pause type by index.
     *
     * @param {number} index
     */
    const removePauseType = index => updatePauseTypes(
      props.modelValue.types.filter((item, itemIndex) => itemIndex !== index),
    );

    /**
     * Replace the fast pbehavior settings.
     *
     * @param {Object[]} fastPbehaviors
     */
    const updateFastPbehaviors = fastPbehaviors =>
      emit('update:fastPbehaviors', fastPbehaviors);

    /**
     * Append an empty fast pbehavior.
     */
    const addPbehavior = () => updateFastPbehaviors([
      ...props.fastPbehaviors,
      alarmListFastPbehaviorParametersItemToForm(),
    ]);

    /**
     * Update one field of a fast pbehavior.
     *
     * @param {number} index
     * @param {string} field
     * @param {string} value
     */
    const updatePbehaviorField = (index, field, value) => {
      const fastPbehaviors = [...props.fastPbehaviors];
      fastPbehaviors[index] = { ...fastPbehaviors[index], [field]: value };
      updateFastPbehaviors(fastPbehaviors);
    };

    /**
     * Remove a fast pbehavior by index.
     *
     * @param {number} index
     */
    const removePbehavior = index => updateFastPbehaviors(
      props.fastPbehaviors.filter((item, itemIndex) => itemIndex !== index),
    );

    return {
      t,
      FORM_MAX_STRING_LENGTH,
      ALARM_LIST_PAUSE_METHODS,
      pbehaviorTypes,
      pauseDragHandleClass,
      pbehaviorDragHandleClass,
      isPauseAction,
      isPbehaviorAction,
      updateMode,
      updatePauseTypes,
      addPauseType,
      updatePauseTypeField,
      removePauseType,
      addPbehavior,
      updatePbehaviorField,
      removePbehavior,
    };
  },
};
</script>

<style lang="scss" scoped>
.pause-form {
  &__mode {
    width: 100%;
  }

  &__parameters {
    width: 100%;
  }
}
</style>
