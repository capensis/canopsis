<template>
  <widget-settings-group :title="t('settings.commentTemplates.title')">
    <v-container>
      <c-layout gap="2" column>
        <span class="text-body-2">
          {{ t('common.template', 2) }}
        </span>
        <span>{{ t('settings.commentTemplates.description') }}</span>
        <c-card-iterator-field
          :data="modelValue"
          :handle="`.${dragItemHandleClass}`"
          item-key="key"
          @update:data="$emit('update:modelValue', $event)"
        >
          <template #item="{ index }">
            <c-card-iterator-item
              :drag-handle-class="dragItemHandleClass"
              small
              @remove="removeTemplate(index)"
            >
              <template #header>
                <c-select
                  :model-value="modelValue[index].template"
                  :items="getAvailableTemplates(modelValue[index].template)"
                  :label="t('common.template', 1)"
                  :name="`commentTemplates[${index}].template`"
                  :loading="pending"
                  item-title="name"
                  item-value="_id"
                  item-disabled="disabled"
                  hide-details
                  required
                  @update:model-value="updateTemplate(index, $event)"
                />
              </template>
            </c-card-iterator-item>
          </template>
        </c-card-iterator-field>
        <c-layout class="mt-2" justify="end">
          <v-btn color="primary" @click="addTemplate">
            {{ t('common.add') }}
          </v-btn>
        </c-layout>
      </c-layout>
    </v-container>
  </widget-settings-group>
</template>

<script>
import { onMounted, ref } from 'vue';

import { uid } from '@/helpers/uid';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useI18n } from '@/hooks/i18n';
import { useCommentTemplate } from '@/hooks/store/modules/comment-template';

import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';

export default {
  components: { WidgetSettingsGroup },

  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const dragItemHandleClass = 'template-drag-handler';
    const availableTemplates = ref([]);

    const { fetchCommentTemplatesListWithoutStore } = useCommentTemplate();
    const { addItemIntoArray, updateFieldInArrayItem, removeItemFromArray } = useArrayModelField(
      props,
      emit,
    );

    const pending = ref(false);

    const fetchTemplates = async () => {
      try {
        pending.value = true;

        const response = await fetchCommentTemplatesListWithoutStore({
          params: { limit: 1000 },
        });

        availableTemplates.value = response.data;
      } finally {
        pending.value = false;
      }
    };

    const getAvailableTemplates = currentTemplate => availableTemplates.value.map(template => ({
        ...template,
        disabled: props.modelValue.some(
          item => item.template === template._id && item.template !== currentTemplate,
        ),
      }));

    const addTemplate = () => {
      addItemIntoArray({
        key: uid(),
        template: '',
      });
    };

    const updateTemplate = (index, value) => updateFieldInArrayItem(index, 'template', value);
    const removeTemplate = index => removeItemFromArray(index);

    onMounted(fetchTemplates);

    return {
      t,
      dragItemHandleClass,
      pending,
      getAvailableTemplates,
      addTemplate,
      updateTemplate,
      removeTemplate,
    };
  },
};
</script>
