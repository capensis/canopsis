<template>
  <c-select
    :model-value="modelValue"
    :items="viewItems"
    :label="label || t('role.defaultView')"
    :name="name"
    :loading="viewGroupsPending"
    :menu-props="menuProps"
    item-title="title"
    item-value="value"
    clearable
    @update:model-value="updateModel"
  >
    <template #subheader="{ props: subheaderProps }">
      <v-list-subheader
        v-bind="subheaderProps"
        class="view-selector-menu__subheader"
      />
    </template>

    <template #item="{ props: itemProps }">
      <v-list-item
        v-bind="itemProps"
        class="view-selector-menu__item"
      />
    </template>
  </c-select>
</template>

<script>
import { computed, onMounted } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useView } from '@/hooks/store/modules/view';

export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'defaultview',
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const {
      viewGroups,
      viewGroupsPending,
      fetchViewGroupsList,
    } = useView();
    const menuProps = {
      contentClass: 'view-selector-menu',
    };

    const viewItems = computed(() => (
      (viewGroups.value ?? []).reduce((acc, group) => {
        const views = group.views ?? [];
        const groupTitle = group.title || group.name;

        acc.push({
          type: 'subheader',
          title: groupTitle,
        });

        views.forEach((view) => {
          acc.push({
            value: view._id,
            title: view.title || view.name,
          });
        });

        return acc;
      }, [])
    ));

    const updateModel = value => {
      emit('update:modelValue', value);
    };

    onMounted(() => fetchViewGroupsList({
      params: {
        limit: MAX_LIMIT,
        page: 1,
        with_views: true,
        with_tabs: true,
        with_widgets: true,
        with_flags: true,
        with_private: true,
      },
    }));

    return {
      t,
      menuProps,
      viewItems,
      viewGroupsPending,
      updateModel,
    };
  },
};
</script>

<style lang="scss">
.view-selector-menu__subheader {
  .v-list-item-title {
    font-weight: 600;
  }
}
</style>
