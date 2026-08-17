<template>
  <alarms-list-table-with-pagination
    :parent-alarm="alarm"
    :widget="widget"
    :alarms="alarms"
    :meta="meta"
    :query="query"
    :columns="widget.parameters.widgetGroupColumns"
    :loading="pending"
    :refresh-alarms-list="refreshAlarmsList"
    :columns-settings="columnsSettings"
    expandable
    eager
    @update:query="$emit('update:query', $event)"
    @select:tag="$emit('select:tag', $event)"
    @remove:tag="$emit('remove:tag', $event)"
  />
</template>

<script>
import { computed, defineAsyncComponent } from 'vue';

export default {
  components: {
    /**
     * We need to do it by defineAsyncComponent to avoid circular dependency.
     */
    AlarmsListTableWithPagination: defineAsyncComponent(
      () => import('../partials/alarms-list-table-with-pagination.vue'),
    ),
  },

  props: {
    children: {
      type: Object,
      required: true,
    },
    alarm: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    query: {
      type: Object,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    refreshAlarmsList: {
      type: Function,
      default: () => {},
    },
    columnsSettings: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:query', 'select:tag', 'remove:tag'],

  setup(props) {
    const alarms = computed(() =>
      (Array.isArray(props.children) ? props.children : (props.children?.data ?? [])),
    );
    const meta = computed(() => (Array.isArray(props.children) ? {} : (props.children?.meta ?? {})));

    return {
      alarms,
      meta,
    };
  },
};
</script>
