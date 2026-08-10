<template>
  <c-layout gap="2" column>
    <span class="text-subtitle-2">
      {{ t('common.fields') }}
    </span>

    <c-card-iterator-form
      :value="fields"
      :handle="`.${dragItemHandleClass}`"
      :empty-message="t('commentTemplate.errors.fieldsRequired')"
      item-key="key"
      name="fields"
      required
      @update:value="fields = $event"
      @add="addField"
    >
      <template #item="{ index }">
        <c-card-iterator-item
          :item-number="index + 1"
          :drag-handle-class="dragItemHandleClass"
          small
          @remove="removeField(index)"
        >
          <template #header>
            <c-layout class="comment-template-fields-form__field" column>
              <c-name-field
                :model-value="fields[index].name"
                :label="t('common.fieldName')"
                :name="`fields[${index}].name`"
                :max-length="255"
                required
                @update:model-value="updateField(index, 'name', $event)"
              />

              <c-enabled-field
                :model-value="fields[index].required"
                :label="t('common.required')"
                class="mt-0"
                hide-details
                @update:model-value="updateField(index, 'required', $event)"
              />
            </c-layout>
          </template>
        </c-card-iterator-item>
      </template>
    </c-card-iterator-form>
  </c-layout>
</template>

<script>
import { onMounted } from 'vue';

import { uid } from '@/helpers/uid';
import { useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  setup() {
    const { t } = useI18n();
    const form = useFormContext();
    const fields = form.model('fields');

    const dragItemHandleClass = 'field-drag-handler';

    /**
     * Add new field to the fields array
     */
    const addField = () => {
      fields.value = [
        ...fields.value,
        {
          key: uid(),
          name: '',
          required: false,
        },
      ];
    };

    /**
     * Remove field from the fields array
     *
     * @param {number} index - Index of the field to remove
     */
    const removeField = (index) => {
      fields.value = fields.value.filter((field, fieldIndex) => fieldIndex !== index);
    };

    const updateField = (index, field, value) => {
      fields.value = fields.value.map((item, itemIndex) =>
        itemIndex === index ? { ...item, [field]: value } : item,
      );
    };

    onMounted(() => {
      if (!fields.value.length) {
        addField();
      }
    });

    return {
      t,
      fields,
      dragItemHandleClass,
      addField,
      removeField,
      updateField,
    };
  },
};
</script>

<style lang="scss" scoped>
.comment-template-fields-form__field {
  flex: 1 1 auto;
  min-width: 0;
}
</style>
