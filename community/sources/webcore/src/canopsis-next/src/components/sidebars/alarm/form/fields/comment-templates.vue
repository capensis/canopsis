<template>
  <widget-settings-group :title="$t('settings.commentTemplates.title')">
    <v-container>
      <v-layout class="gap-2" column>
        <span class="text-body-2">{{ $tc('common.template', 2) }}</span>
        <span>{{ $t('settings.commentTemplates.description') }}</span>
        <c-card-iterator-field
          v-field="templates"
          :handle="`.${dragItemHandleClass}`"
          item-key="key"
        >
          <template #item="{ index }">
            <c-card-iterator-item
              :drag-handle-class="dragItemHandleClass"
              small
              @remove="removeTemplate(index)"
            >
              <template #header>
                <c-select-field
                  v-field="templates[index].template"
                  :items="getAvailableTemplates(templates[index].template)"
                  :label="$tc('common.template', 1)"
                  :name="`commentTemplates[${index}].template`"
                  :loading="pending"
                  item-text="name"
                  item-value="_id"
                  item-disabled="disabled"
                  hide-details
                  required
                />
              </template>
            </c-card-iterator-item>
          </template>
        </c-card-iterator-field>
        <v-layout
          class="mt-2"
          justify-end
        >
          <v-btn
            color="primary"
            @click="addTemplate"
          >
            {{ $t('common.add') }}
          </v-btn>
        </v-layout>
      </v-layout>
    </v-container>
  </widget-settings-group>
</template>

<script>
import { ref, onMounted } from 'vue';

import { uid } from '@/helpers/uid';

import { useCommentTemplates } from '@/hooks/store/modules/comment-template';
import { useArrayModelField } from '@/hooks/form/array-model-field';
import { usePendingHandler } from '@/hooks/query/pending';

import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';

export default {
  inject: ['$validator'],
  components: { WidgetSettingsGroup },
  model: {
    prop: 'templates',
    event: 'input',
  },
  props: {
    templates: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const dragItemHandleClass = 'template-drag-handler';

    const availableTemplates = ref([]);

    const { fetchCommentTemplatesListWithoutStore } = useCommentTemplates();
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    /**
     * Fetch available comment templates
     */
    const {
      pending,
      handler: fetchTemplates,
    } = usePendingHandler(async () => {
      const response = await fetchCommentTemplatesListWithoutStore({
        params: { limit: 1000 },
      });

      availableTemplates.value = response.data;
    });

    /**
     * Get available templates with disabled state for already selected ones
     *
     * @param {string} currentTemplate - Current template ID to exclude from disabling
     * @returns {Array} Array of templates with disabled property
     */
    const getAvailableTemplates = currentTemplate => (
      availableTemplates.value.map(template => ({
        ...template,
        disabled: props.templates.some(item => item.template === template._id && item.template !== currentTemplate),
      }))
    );

    /**
     * Add new template to the templates array
     */
    const addTemplate = () => {
      addItemIntoArray({
        key: uid(),
        template: '',
      });
    };

    /**
     * Remove template from the templates array
     *
     * @param {number} index - Index of the template to remove
     */
    const removeTemplate = index => removeItemFromArray(index);

    onMounted(() => {
      fetchTemplates();
    });

    return {
      dragItemHandleClass,
      pending,
      getAvailableTemplates,
      addTemplate,
      removeTemplate,
    };
  },
};
</script>
