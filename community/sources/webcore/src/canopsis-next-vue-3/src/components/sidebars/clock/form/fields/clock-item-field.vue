<template>
  <c-card-iterator-item
    :drag-handle-class="dragHandleClass"
    :removable="removable"
    :expandable="false"
    class="clock-item-field"
    small
    @remove="$emit('remove')"
  >
    <template #header>
      <v-row>
        <v-col cols="12">
          <c-timezone-field
            :model-value="modelValue.timezone"
            :name="`${name}.timezone`"
            :label="t('clock.fields.timezone')"
            required
            @update:model-value="updateField('timezone', $event)"
          />
        </v-col>
        <v-col cols="12">
          <c-text-field
            :model-value="modelValue.label"
            :name="`${name}.label`"
            :label="t('clock.fields.customLabel')"
            :max-length="255"
            @update:model-value="updateField('label', $event)"
          />
        </v-col>
      </v-row>
    </template>
  </c-card-iterator-item>
</template>

<script>
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    dragHandleClass: {
      type: String,
      required: true,
    },
    removable: {
      type: Boolean,
      default: true,
    },
  },

  emits: ['update:modelValue', 'remove'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    return {
      t,
      updateField,
    };
  },
};
</script>

<style lang="scss" scoped>
.clock-item-field {
  :deep(.c-card-iterator-item__header) {
    padding-right: 28px;
  }
}
</style>
