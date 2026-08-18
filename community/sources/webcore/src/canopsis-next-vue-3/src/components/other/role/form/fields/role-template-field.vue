<template>
  <c-select
    :model-value="selectedTemplate"
    :label="label || t('role.selectTemplate')"
    :loading="pending"
    :items="preparedItems"
    :disabled="disabled"
    name="template"
    item-title="name"
    item-value="_id"
    return-object
    clearable
    @update:model-value="updatePermissions"
  >
    <template #item="{ props: itemProps, item }">
      <v-list-item v-bind="getTemplateItemProps(itemProps)">
        <c-layout justify="space-between">
          {{ item?.name }}

          <span class="role-template-field__item-description">
            {{ item?.description }}
          </span>
        </c-layout>
      </v-list-item>
    </template>
  </c-select>
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { rolePermissionsToForm } from '@/helpers/entities/role/form';

import { useI18n } from '@/hooks/i18n';
import { useRole } from '@/hooks/store/modules/role';

export default {
  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { fetchRoleTemplatesListWithoutStore } = useRole();
    const items = ref([]);
    const pending = ref(false);
    const selectedTemplate = ref(null);

    const preparedItems = computed(() =>
      items.value.map(template => ({
        ...template,
        permissions: rolePermissionsToForm(template.permissions),
      })),
    );

    const updatePermissions = (template) => {
      selectedTemplate.value = template;
      emit('update:modelValue', template ? template.permissions : {});
    };

    const getTemplateItemProps = ({ title, subtitle, ...itemProps } = {}) => itemProps;

    const fetchList = async () => {
      try {
        pending.value = true;

        const { data = [] } = await fetchRoleTemplatesListWithoutStore({
          params: { limit: MAX_LIMIT },
        });

        items.value = data;
      } finally {
        pending.value = false;
      }
    };

    onMounted(fetchList);

    return {
      t,
      pending,
      selectedTemplate,
      preparedItems,
      updatePermissions,
      getTemplateItemProps,
    };
  },
};
</script>

<style lang="scss">
.role-template-field__item-description {
  opacity: 0.7;
}
</style>
