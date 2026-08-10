<template>
  <c-layout gap="3" column>
    <c-enabled-field v-model="visible" :label="t('pbehavior.visible')" with-background />

    <c-name-field
      v-model="name"
      :disabled="onlyColor"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="t('modals.createPbehaviorType.fields.description')">
        <c-description-field
          v-model="description"
          :disabled="onlyColor"
          name="description"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="t('modals.createPbehaviorType.fields.type')">
        <pbehavior-type-default-type-field
          v-model="type"
          :disabled="onlyColor"
          @sync:color="syncColor"
          @update:color="updateColor"
        />
      </c-form-block-row>

      <c-form-block-row :label="t('modals.createPbehaviorType.fields.priority')">
        <c-priority-field
          v-model="priority"
          :disabled="onlyColor"
          :loading="pendingPriority"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="t('modals.createPbehaviorType.fields.iconName')">
        <c-icon-field
          v-model="iconName"
          :placeholder="t('modals.createPbehaviorType.iconNamePlaceholder')"
          :disabled="onlyColor"
          :required="!onlyColor"
          name="icon_name"
        >
          <template #no-data>
            <v-list-item>
              <v-list-item-title v-html="t('modals.createPbehaviorType.errors.iconName')" />
            </v-list-item>
          </template>
        </c-icon-field>
      </c-form-block-row>

      <c-form-block-row :label="t('common.color')" align-center>
        <c-layout gap="2" align-center justify-start>
          <span v-if="onlyColor" class="text-medium-emphasis">
            {{ t('pbehavior.types.defaultType') }}
          </span>
          <c-color-picker-field
            v-model="color"
            :justify-end="onlyColor"
            name="color"
            required
            @update:model-value="setColorWasChanged"
          />
        </c-layout>
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { ref } from 'vue';

import { useFormContext, useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import PbehaviorTypeDefaultTypeField from './pbehavior-type-default-type-field.vue';

export default {
  components: { PbehaviorTypeDefaultTypeField },

  props: {
    onlyColor: {
      type: Boolean,
      default: false,
    },
    pendingPriority: {
      type: Boolean,
      default: false,
    },
    syncOriginalForm: {
      type: Function,
      default: undefined,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const form = useFormContext();
    const { name, description, type, priority, color, visible, iconName } = useFormModels({
      name: 'name',
      description: 'description',
      type: 'type',
      priority: 'priority',
      color: 'color',
      visible: 'visible',
      iconName: 'icon_name',
    });

    const colorWasChanged = ref(false);

    const setColorWasChanged = () => {
      colorWasChanged.value = true;
    };

    const syncColor = (nextColor) => {
      if (color.value) {
        return;
      }

      form.reset({
        ...form.values.value,

        color: nextColor,
      });
      props.syncOriginalForm?.();
    };

    const updateColor = (nextColor) => {
      if (color.value) {
        if (colorWasChanged.value) {
          return;
        }

        if (color.value !== nextColor) {
          colorWasChanged.value = true;

          return;
        }
      }

      color.value = nextColor;
    };

    return {
      t,
      name,
      description,
      type,
      priority,
      color,
      visible,
      iconName,
      setColorWasChanged,
      syncColor,
      updateColor,
    };
  },
};
</script>
