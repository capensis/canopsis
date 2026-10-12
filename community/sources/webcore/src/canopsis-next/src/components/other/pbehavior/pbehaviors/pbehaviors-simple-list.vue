<template>
  <c-layout class="pbehaviors-simple-list" column>
    <c-layout align="center" justify="space-between">
      <c-timezone-field
        v-if="shownUserTimezone"
        v-model="timezone"
        class="pbehaviors-simple-list__timezone"
        server
      />

      <c-layout gap="2" justify="end" class="ml-auto">
        <c-action-fab-btn
          v-if="addable"
          :tooltip="$t('modals.createPbehavior.create.title')"
          icon="add"
          color="primary"
          left
          @click="showCreatePbehaviorModal"
        />
        <c-action-fab-btn
          :tooltip="$t('modals.pbehaviorsCalendar.title')"
          icon="calendar_today"
          color="secondary"
          left
          @click="showPbehaviorsCalendarModal"
        />
      </c-layout>
    </c-layout>

    <c-advanced-data-table
      :items="pbehaviors"
      :headers="headers"
      :loading="pending"
      :density="dense ? 'compact' : 'default'"
      class="pbehaviors-simple-list__table"
      hide-actions
      no-pagination
    >
      <template #enabled="{ item }">
        <c-enabled :value="item.enabled" />
      </template>
      <template #tstart="{ item }">
        {{ formatIntervalDate(item, 'tstart') }}
      </template>
      <template #tstop="{ item }">
        {{ formatIntervalDate(item, 'tstop') }}
      </template>
      <template #rrule_end="{ item }">
        {{ formatRruleEndDate(item) }}
      </template>
      <template #rrule="{ item }">
        <v-icon :icon="item.rrule ? 'check' : 'clear'" />
      </template>
      <template #icon="{ item }">
        <v-icon :icon="item.type?.icon_name" color="primary" />
      </template>
      <template #status="{ item }">
        <v-icon
          :color="item.is_active_status ? 'primary' : 'error'"
          icon="$vuetify.icons.settings_sync"
        />
      </template>
      <template #actions="{ item }">
        <pbehavior-actions
          :pbehavior="item"
          :removable="removable"
          :updatable="updatable"
          :with-inherited="withInherited"
          @refresh="fetchList"
        />
      </template>
    </c-advanced-data-table>
  </c-layout>
</template>

<script>
import { computed, inject, onBeforeUnmount, onMounted, ref } from 'vue';

import { MODALS } from '@/constants';

import Observer from '@/services/observer';

import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { usePbehavior } from '@/hooks/store/modules/pbehavior';

import { usePbehaviorDateFormat } from '@/components/other/pbehavior/pbehaviors/hooks/pbehavior-date-format';

import PbehaviorActions from './partials/pbehavior-actions.vue';

export default {
  components: { PbehaviorActions },

  props: {
    entity: {
      type: Object,
      required: true,
    },
    removable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    dense: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    withActiveStatus: {
      type: Boolean,
      default: false,
    },
    withInherited: {
      type: Boolean,
      default: false,
    },
    defaultInherited: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const modals = useModals();
    const { t } = useI18n();
    const { fetchPbehaviorsByEntityIdWithoutStore } = usePbehavior();
    const { timezone, shownUserTimezone, formatIntervalDate, formatRruleEndDate } =
      usePbehaviorDateFormat();

    const periodicRefresh = inject('$periodicRefresh', new Observer());

    const pbehaviors = ref([]);

    const headers = computed(() => {
      const result = [
        { title: t('common.name'), value: 'name' },
        { title: t('common.author'), value: 'author.display_name' },
        { title: t('pbehavior.isEnabled'), value: 'enabled' },
        { title: t('pbehavior.begins'), value: 'tstart' },
        { title: t('pbehavior.ends'), value: 'tstop' },
        { title: t('pbehavior.rruleEnd'), value: 'rrule_end' },
        { title: t('common.recurrence'), value: 'rrule' },
        { title: t('common.type'), value: 'type.name' },
        { title: t('common.reason'), value: 'reason.name' },
        { title: t('common.icon', 1), value: 'icon' },
      ];

      if (props.withActiveStatus) {
        result.push({ title: t('common.status'), value: 'status', sortable: false });
      }

      if (props.updatable || props.removable) {
        result.push({ title: t('common.actionsLabel'), value: 'actions', sortable: false });
      }

      return result;
    });

    const { pending, fetchHandlerWithQuery: fetchList } = usePendingWithLocalQuery({
      fetchHandler: async () => {
        try {
          pbehaviors.value = await fetchPbehaviorsByEntityIdWithoutStore({
            id: props.entity._id,
            params: {
              with_flags: true,
            },
          });
        } catch (err) {
          console.warn(err);
        }
      },
    });

    const showPbehaviorsCalendarModal = () => {
      modals.show({
        name: MODALS.pbehaviorsCalendar,
        config: {
          title: t('modals.pbehaviorsCalendar.entity.title', { name: props.entity.name }),
          entityId: props.entity._id,
          defaultInherited: props.defaultInherited,
          withInherited: props.withInherited,
        },
      });
    };

    const showCreatePbehaviorModal = () => {
      modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: createEntityIdPatternByValue(props.entity._id),
          entities: [props.entity],
          defaultInherited: props.defaultInherited,
          withInherited: props.withInherited,
          afterSubmit: fetchList,
        },
      });
    };

    onMounted(() => {
      fetchList();

      periodicRefresh.register(fetchList);
    });

    onBeforeUnmount(() => periodicRefresh.unregister(fetchList));

    return {
      timezone,
      shownUserTimezone,
      pending,
      pbehaviors,
      headers,

      formatIntervalDate,
      formatRruleEndDate,
      fetchList,
      showPbehaviorsCalendarModal,
      showCreatePbehaviorModal,
    };
  },
};
</script>

<style lang="scss" scoped>
.pbehaviors-simple-list {
  max-width: 100%;
  min-width: 0;
}

.pbehaviors-simple-list__table {
  max-width: 100%;
  min-width: 0;
}

.pbehaviors-simple-list__timezone {
  max-width: 420px;
}
</style>
