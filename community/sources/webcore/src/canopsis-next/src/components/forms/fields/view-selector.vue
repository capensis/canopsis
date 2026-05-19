<template>
  <c-select-field
    v-field="value"
    :items="viewItems"
    :label="label || $t('role.defaultView')"
    :name="name"
    :loading="groupsPending"
    :menu-props="menuProps"
    clearable
    ellipsis
  >
    <template #item="{ item, attrs, on }">
      <v-subheader
        v-if="item.header"
        :key="item.header"
        class="view-selector__header"
      >
        {{ item.header }}
      </v-subheader>
      <v-list-item
        v-else
        v-bind="attrs"
        class="view-selector__item"
        v-on="on"
      >
        <v-list-item-content>
          <v-list-item-title>{{ item.text }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </template>
    <template #selection="">
      <span class="text-truncate">
        {{ selectedViewLabel }}
      </span>
    </template>
  </c-select-field>
</template>

<script>
import { computed, onMounted } from 'vue';

import { useViewGroup } from '@/hooks/store/modules/view';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
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
  setup(props) {
    const {
      groups,
      groupsPending,
      getViewById,
      fetchAllGroupsListWithWidgets,
    } = useViewGroup();

    const menuProps = {
      contentClass: 'view-selector-menu',
    };

    const viewItems = computed(() => (
      (groups.value ?? []).reduce((acc, group) => {
        const views = group.views ?? [];

        if (!views.length) {
          return acc;
        }

        acc.push({ header: group.title || group.name });

        views.forEach((view) => {
          acc.push({
            value: view._id,
            text: view.title || view.name,
            groupTitle: group.title || group.name,
          });
        });

        return acc;
      }, [])
    ));

    const selectedViewLabel = computed(() => {
      if (!props.value) {
        return '';
      }

      const selectedItem = viewItems.value.find(({ value }) => value === props.value);

      if (selectedItem) {
        return `${selectedItem.groupTitle} / ${selectedItem.text}`;
      }

      try {
        const view = getViewById.value(props.value);

        return view?.title ?? '';
      } catch (error) {
        console.error(error);

        return '';
      }
    });

    onMounted(fetchAllGroupsListWithWidgets);

    return {
      menuProps,
      viewItems,
      groupsPending,
      selectedViewLabel,
    };
  },
};
</script>

<style lang="scss">
.view-selector-menu {
  .view-selector__header {
    font-weight: 700;
    height: auto;
    min-height: 32px;
  }

  .view-selector__item {
    padding-left: 32px !important;
  }
}
</style>
